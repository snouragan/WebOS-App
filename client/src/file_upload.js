var box = document.getElementById("box");
var list = '';

box.addEventListener('click', () => {
    console.log('clicked box');
    if(box.classList.contains('box')) {
        box.classList.add('largebox');
        box.classList.remove('box');
    }

    else if(box.classList.contains('largebox')) {
        box.classList.add('box');
        box.classList.remove('largebox');
    }

});

async function getList() {
    const filelist = document.querySelector('.filelist');

    list = await fetch('http://192.168.1.15:8069/ctl/list').json();
    list.foreach(file => {
        Object.entries(file).forEach(([key, value]) => {
            console.log('${key} ${value}');
            
        });
        console.log('got json');
    });

    setTimeout(getList, 5000);
}

// {
//     "resources": {
//       "38ea1517085ce4233b8251d343b70e5d1ba1d568746bcb86dbc8551dcff4ac5d.webm": {
//         "title": "short",
//         "inprogress": false,
//         "prepared": true,
//         "nmonitors": 2,
//         "sf": "stretch",
//         "thumbnail": "/data/38ea1517085ce4233b8251d343b70e5d1ba1d568746bcb86dbc8551dcff4ac5d.webm.thumb.jpg"
//       },
//       "67a6e9659e36df8568ac991eab70e2c7ca833cdbce38e7dd72d4cb2e7a29ab76.webm": {
//         "title": "video",
//         "inprogress": true,
//         "prepared": false,
//         "thumbnail": "/data/67a6e9659e36df8568ac991eab70e2c7ca833cdbce38e7dd72d4cb2e7a29ab76.webm.thumb.jpg"
//       }
//     }
//   }