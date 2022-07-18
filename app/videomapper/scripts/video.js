// Store video player element:
const videoElement = document.getElementById("video");

// Object to store video file:
let videoFile = undefined;

// Flag for looping videos:
let loop = true;

// Logging context:
const videoContext = 'VIDEO';

/**
 * Resume video playback.
 */
function resumePlayback() {
    sync();
}

/**
 * Pause video playback.
 */
function pausePlayback() {
    videoElement.pause();
}

/**
 * Set player to a specific position.
 * 
 * @param {*} position - position to seek to.
 */
function seekVideo(position) {
    videoElement.currentTime = position;
}

/**
 * Fetch video file and save as blob.
 */
function loadVideo(src) {
    fetch('http://' + SERVER_IP + src, {
        method: 'GET'
    }).then(res => res.blob()).then(blob => {
        videoFile = blob;
        display(videoFile, videoElement);
        log(videoContext, 'Loaded video: ' + src, INFO);
        sync();
    }).catch(e => {
        log(videoContext, 'Failed to fetch video. Reason: ' + e, ERROR);
    });
}

/**
 * Synchronize player with server.
 */
function sync() {
    fetch('http://' + SERVER_IP + '/sync', {
        method: 'GET'
    }).then(res => {
        videoElement.play().catch(() => {
            log(videoContext, 'Error playing video', ERROR);
        });
    }).catch(e => {
        log(videoContext, 'Failed to sync. Reason: ' + e, ERROR);
    });
}

/**
 * Load a video file into the player.
 * source: https://stackoverflow.com/questions/14317179/display-a-video-from-a-blob-javascript
 * 
 * @param {*} videoFile - video file to be loaded.
 * @param {*} videoEl - video player element.
 */
function display(videoFile, videoEl) {

    if (!(videoFile instanceof Blob)) throw new Error('`videoFile` must be a Blob or File object.'); // The `File` prototype extends the `Blob` prototype, so `instanceof Blob` works for both.
    if (!(videoEl instanceof HTMLVideoElement)) throw new Error('`videoEl` must be a <video> element.');

    const newObjectUrl = URL.createObjectURL(videoFile);

    const oldObjectUrl = videoEl.currentSrc;
    if (oldObjectUrl && oldObjectUrl.startsWith('blob:')) {
        videoEl.src = '';
        URL.revokeObjectURL(oldObjectUrl);
    }
    videoEl.src = newObjectUrl;
    videoEl.load();
}

// Ask for server sync when playback ends:
videoElement.addEventListener('ended', () => {
    videoElement.currentTime = 0;
    if (loop) {
        videoElement.currentTime = 0;
        sync();
    }
});

// Load default video:
loadVideo('/split/sample.6.stretch.webm');
