const INFO = 0;
const WARN = 1;
const ERROR = 2;
const CRITICAL = 3;

const styling = [
    'color: black',
    'color: gold',
    'color: red',
    'background-color: darkred; color: white; font-weight: bold'
]

function log(context, content, severity = 0) {
    if (LOGGING) {
        console.log('%c' + new Date().toDateString() + ' [' + context.toUpperCase() + '] - ' + content, styling[severity]);
    }
}
