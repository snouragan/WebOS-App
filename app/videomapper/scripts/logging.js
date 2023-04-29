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

    if (VISUAL_LOGGING) {
        if (content === 'No command') {
            return;
        }
        document.querySelector('p.visual-debug').innerHTML = content;
        document.querySelector('p.visual-debug').classList.add('shown');
        document.querySelector('p.visual-debug').addEventListener('transitionend', () => {
            setTimeout(() => {
                document.querySelector('p.visual-debug').classList.remove('shown');
            }, 2000);
        })
    }
}
