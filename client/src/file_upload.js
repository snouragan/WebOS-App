let box = document.getElementById("box");
let list = {};

box.addEventListener('click', () => {
    console.log('clicked box');
    if (box.classList.contains('box')) {
        box.classList.add('largebox');
        box.classList.remove('box');
    }

    else if (box.classList.contains('largebox')) {
        box.classList.add('box');
        box.classList.remove('largebox');
    }

});

document.querySelector('input').addEventListener('click', (event) => {
    event.stopPropagation();
});

let selectTvF = 0;
let selectTypeF = 0;
let titleF = 0;
let uploadF = 0;
let formData = new FormData();

document.getElementById('file-input').addEventListener('change', () => {

    let options = [1, 2, 3, 4, 5, 6];
    let selectTv = document.createElement('select');
    selectTv.id = 'tvnum';

    let type = ['stretch', 'fit'];
    let selectType = document.createElement('select');
    selectType.id = 'pictype';

    let title = document.createElement('input');
    title.id = 'title';
    title.setAttribute('type', 'text');
    title.setAttribute('value', 'title');

    let upload = document.createElement('div');
    upload.classList.add('upload');
    upload.classList.add('card');
    upload.id = 'upload';

    file = document.getElementById('file-input').files[0];

    if (formData != null)
        formData.delete('file');

    formData.append('file', file);

    console.log(file);

    if (selectTvF == 0) {
        document.querySelector('.addfile').appendChild(selectTv);

        for (var i = 0; i < options.length; i++) {
            var option = document.createElement("option");
            option.value = options[i];
            option.text = options[i];
            option.id = 'option';
            selectTv.appendChild(option);
        }

        selectTvF = 1;

    }

    if (selectTypeF == 0) {
        document.querySelector('.addfile').appendChild(selectType);

        for (var i = 0; i < type.length; i++) {
            var types = document.createElement("option");
            types.value = type[i];
            types.text = type[i];
            types.id = 'types';
            selectType.appendChild(types);
        }

        selectTypeF = 1;

    }

    if (titleF == 0) {
        document.querySelector('.addfile').appendChild(title);

        titleF = 1;
    }

    if (uploadF == 0) {
        document.querySelector('.addfile').appendChild(upload);

        uploadF = 1;
    }

}, false);

let num = 0;
let type = 'stretch';
let title = 'title';

document.querySelector('.addfile').addEventListener('click', (event) => {
    if (event.target.id == 'tvnum') {
        event.stopPropagation();
    }
    if (event.target.id == 'option') {
        event.stopPropagation();
    }
    if (event.target.id == 'pictype') {
        event.stopPropagation();
    }
    if (event.target.id == 'title') {
        event.stopPropagation();
    }
    if (event.target.id == 'types') {
        event.stopPropagation();
    }
    if (event.target.id == 'upload') {

        formData.append('n', num);
        formData.append('sf', type);
        formData.append('title', title);

        console.log(num);
        console.log(type);
        console.log(title);

        fetch('http://192.168.1.109:8069/ctl/upload',
            {
                method: 'POST',
                body: formData
            });
        event.stopPropagation();
    }
});

document.querySelector('.addfile').addEventListener('change', (event) => {
    if (event.target.id == 'tvnum') {
        num = event.target.value;
    }
    if (event.target.id == 'pictype') {
        type = event.target.value;
    }
    if (event.target.id == 'title') {
        title = event.target.value;
    }
});

let thumbnailKeys = [];

(function getList() {
    const filelist = document.querySelector('.filelist');

    fetch('http://192.168.1.109:8069/ctl/list')
        .then(res => res.json())
        .then(data => {
            for (let key in data.resources) {
                let obj = data.resources[key];

                if (!thumbnailKeys.includes(key)) {
                    thumbnailKeys.push(key);
                    addThumbnail(obj.thumbnail);
                }
            }
        });

    setTimeout(getList, 5000);
})();

function addThumbnail(link) {
    let thumbnailURL = "http://192.168.1.109:8069".concat(link);

    let thumbnail = document.createElement('div');
    thumbnail.classList.add('thumbnail');
    thumbnail.classList.add('thumbnail');
    console.log("url('" + thumbnailURL + "')");
    thumbnail.style.backgroundImage = "url('" + thumbnailURL + "')";
    // // thumbnail.style.backgroundImage = 'url(/home/snou/Pictures/rat-car-9553952.png)';
    // thumbnail.style.backgroundImage = 'url(/assets/example.png)';
    console.log(thumbnailURL);

    document.querySelector('.filelist').appendChild(thumbnail);
};

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