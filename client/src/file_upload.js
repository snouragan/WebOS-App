var box = document.getElementById("box");

box.addEventListener('click', () => {
    console.log('clicked box');
    if(box.classList.contains('box')) {
        box.classList.add('largebox');
        box.classList.remove('box');

        // const filelist = document.createElement('div');
        // filelist.classList.add('filelist');
        // document.querySelector('.largebox').appendChild(filelist);
    }

    else if(box.classList.contains('largebox')) {
        box.classList.add('box');
        box.classList.remove('largebox');
        // const filelist = document.querySelector('.filelist'); 
        // box.removeChild(filelist);
    }

});

