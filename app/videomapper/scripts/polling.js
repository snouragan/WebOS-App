// Set server url:
const pollURL = 'http://' + SERVER_IP + '/poll';

// Store polling interval:
let pollInterval;

// Logging context:
const commandsContext = 'COMMANDS';

/**
 * Start listening for commands from the server.
 */
function startPolling() {
    pollInterval = setInterval(() => {
        fetch(pollURL, {
            method: 'GET'
        }).then(res => res.json()).then(data => {
            if (data && data.command) {
                handleCommand(data);
            }
        }).catch(e => {
            log(commandsContext, 'Failed to fetch command. Reason: ' + e, ERROR);
        });
    }, POLL_RATE);
}

/**
 * Stop listening for commands from the server.
 */
function stopPolling() {
    clearInterval(pollInterval);
}

/**
 * Handle a command from the server.
 * 
 * @param {*} command - command to be handled.
 * @returns an action depending on the command.
 */
function handleCommand(command) {
    if (!command || !command.command) {
        return log(commandsContext, 'Invalid command', WARN);
    }

    if (command.command === 'none') {
        return log(commandsContext, 'No command', INFO);
    }

    if (command.command === 'play') {
        log(commandsContext, 'command: play', INFO);
        return resumePlayback();
    }

    if (command.command === 'pause') {
        log(commandsContext, 'command: pause', INFO);
        return pausePlayback();
    }

    if (command.command === 'seek' && command.pos) {
        log(commandsContext, 'command: play', INFO);
        return seekVideo(command.pos);
    }

    if (command.command === 'load' && command.src) {
        return loadVideo(command.src)
    }

    return log(commandsContext, 'Unknown command: ' + command.command, WARN);
}

// Begin listening for commands:
startPolling();
