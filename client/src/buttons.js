document.querySelector('.select').addEventListener('click', () => {
    if(document.querySelector('.TV-container').classList.contains('selection')) {
        document.querySelector('.TV-container').classList.remove('selection');
        
    }
    else {
        document.querySelector('.TV-container').classList.add('selection')
        console.log('selection mode');
    }
});