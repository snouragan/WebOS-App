document.querySelector('.select').addEventListener('click', () => {
    if(document.querySelector('.TV-container').classList.contains('selection')) {
        document.querySelector('.TV-container').classList.remove('selection');
        
    }
    else {
        document.querySelector('.TV-container').classList.add('selection')
        console.log('selection mode');
    }
});

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
