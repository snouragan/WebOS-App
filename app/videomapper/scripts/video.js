// Store video player and image elements:
const videoElement = document.getElementById("video");
const imageElement = document.querySelector("img.fullscreen");

// Object to store video file and image:
let videoFile = undefined;
let iamgeFile = undefined;

// Flag for deciding between video and image:
let displayVideo = true;

// Flag for looping videos:
let loop = true;

// Logging context:
const videoContext = 'VIDEO';

/**
 * Resume video playback.
 */
function resumePlayback() {
    if (!displayVideo) {
        return;
    }
    syncPlay();
}

/**
 * Pause video playback.
 */
function pausePlayback() {
    if (!displayVideo) {
        return;
    }
    videoElement.pause();
    syncCurrentPosition();
}

/**
 * Set player to a specific position.
 * 
 * @param {*} position - position to seek to.
 */
function seekVideo(position) {
    if (!displayVideo) {
        return;
    }
    videoElement.pause();
    videoElement.currentTime = position;
    syncPlay();
}

/**
 * Fetch video file and save as blob.
 */
function loadVideo(src) {
    stopPolling();
    fetch('http://' + SERVER_IP + src, {
        method: 'GET'
    }).then(res => res.blob()).then(blob => {
        displayVideo = true;
        imageElement.classList.remove('visible');
        videoFile = blob;
        display(videoFile, videoElement);
        log(videoContext, 'Loaded video: ' + src, INFO);
        syncPlay();
    }).catch(e => {
        log(videoContext, 'Failed to fetch video. Reason: ' + e, ERROR);
        startPolling();
    });
}

function loadImage(src) {
    stopPolling();
    fetch('http://' + SERVER_IP + src, {
        method: 'GET'
    }).then(res => res.blob()).then(blob => {
        displayVideo = false;
        imageElement.classList.add('visible');
        imageFile = blob;
        let urlCreator = window.URL || window.webkitURL;
        imageElement.src = urlCreator.createObjectURL(imageFile);
        log(videoContext, 'Loaded image: ' + src, INFO);
        startPolling();
    }).catch(e => {
        log(videoContext, 'Failed to fetch image. Reason: ' + e, ERROR);
        startPolling();
    });
}

/**
 * Synchronize player with other devices.
 */
function syncPlay() {
    if (!displayVideo) {
        return;
    }
    stopPolling();
    fetch('http://' + SERVER_IP + '/sync', {
        method: 'GET'
    }).then(res => {
        videoElement.play().catch(() => {
            log(videoContext, 'Error playing video', ERROR);
        });
        startPolling();
    }).catch(e => {
        log(videoContext, 'Failed to sync. Reason: ' + e, ERROR);
        startPolling();
    });
}

/**
 * Synchronize pause position with other devices.
 */
function syncCurrentPosition() {
    if (!displayVideo) {
        return;
    }
    stopPolling();
    fetch('http://' + SERVER_IP + '/sync', {
        method: 'POST',
        body: JSON.stringify({ 'time': videoElement.currentTime })
    }).then(res => res.json()).then(data => {
        if (data && data.time) {
            videoContext.currentTime = data.time;
            startPolling();
        }
    })
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
        syncPlay();
    }
});


window.addEventListener('load', () => {
    log(videoContext, 'Loaded video player', INFO);
    // Load default video:
    // loadVideo('/split/sample.6.stretch.webm');

    startPolling();
});
