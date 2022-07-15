const fileInput = document.getElementById('file-input');
const dropZone = document.getElementById('drop-zone');
const prompt = document.getElementById('drop-zone-prompt');
const container = document.querySelector('.container');


container.addEventListener('dragover', (e) => {
    e.preventDefault();
    e.stopPropagation();    

    prompt.innerHTML = 'Drop here to upload';
    dropZone.classList.add('dragging');
});

window.addEventListener('dragleave', (e) => {
    e.preventDefault();

    prompt.innerHTML = 'Drop file here or click to upload';
    dropZone.classList.remove('dragging');
});

window.addEventListener('drop', (e) => {
    e.preventDefault();

    prompt.innerHTML = 'Drop file here or click to upload';
    dropZone.classList.remove('dragging');
});

dropZone.addEventListener('drop', (e) => {
    e.preventDefault();
    e.stopPropagation();

    // const file = e.dataTransfer.items[0].getAsFile();
    // console.log('dropped ' + file);
    // const data = new FormData(document.getElementById('form'));
    // data.append('file', file)
    // fetch('http://192.168.84.7:8069/ctl/upload', {
    //     method: 'POST',
    //     body: data
    // });
});

fileInput.addEventListener('change', () => {
    console.log(fileInput.files[0]);
    const data = new FormData(document.getElementById('form'));
    data.append('file', fileInput.files[0])
    fetch('http://192.168.84.7:8069/ctl/upload', {
        method: 'POST',
        body: data
    });
});

dropZone.addEventListener('click', () => {
    fileInput.click();
});