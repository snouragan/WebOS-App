document.querySelector('.play').addEventListener('click', ()=> {
    console.log('clicked play');
    fetch('http://192.168.1.15:8069/ctl/play', {
        method: 'GET',
    });
});

document.querySelector('.pause').addEventListener('click', ()=> {
    console.log('clicked pause');
    fetch('http://192.168.1.15:8069/ctl/pause', {
        method: 'GET',
    });
});