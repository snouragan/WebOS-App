let selectionF = 0

document.querySelectorAll('.TV').forEach((tv) => {
    tv.addEventListener('click', () => {
        if(tvurl.length == 0)
            tvurl = tvurl + tv.id.split('tv')[1];
        else 
            tvurl = tvurl + ',' + tv.id.split('tv')[1];
        console.log(tvurl);
        console.log(tv.id, "clicked");
        if (selectionF) {
            if (tv.classList.contains('clicked')) {
                tv.classList.remove('clicked');
            } else {
                tv.classList.add('clicked');

                document.querySelectorAll('.TV').forEach((t) => {
                    if (t.id !== tv.id) {
                        t.classList.remove('clicked');
                    }
                })
            }
        }
    });
});

let src = '';
let tvurl = '';
let tvthumbnail = [0, 0, 0, 0, 0, 0];

document.querySelector('.load').addEventListener('click', () => {
    fetch('http://' + ip + '192.168.1.77:8069/ctl/load?src=' + src + '&tv=' + tvurl, {
        method: 'GET'
    }); 

    let splitter = tvurl.split(',');
    splitter.forEach(i => {
        tvthumbnail[i] = src;
    });

    console.log(tvthumbnail);
    showThumbnail();
    tvthumbnail = ['none', 'none', 'none', 'none', 'none', 'none'];

    src = '';
    tvurl = '';
});

function showThumbnail() {
    document.querySelectorAll('.TV').forEach( t => {
        t.style.backgroundImage = tvthumbnail[t.id.split('tv')[1]]
    });
};

document.querySelectorAll('.TV').forEach((tv) => {
    tv.addEventListener('mouseover', () => {
        console.log(tv.id, "hover");
        if (tv.id === "tv0") {
            document.getElementById("mtv0").style.background = "red";
        }
        if (tv.id === "tv1") {
            document.getElementById("mtv1").style.background = "red";
        }
        if (tv.id === "tv2") {
            document.getElementById("mtv2").style.background = "red";
        }
        if (tv.id === "tv3") {
            document.getElementById("mtv3").style.background = "red";
        }
        if (tv.id === "tv4") {
            document.getElementById("mtv4").style.background = "red";
        }
        if (tv.id === "tv5") {
            document.getElementById("mtv5").style.background = "red";
        }
    })
})

document.querySelectorAll('.TV').forEach((tv) => {
    tv.addEventListener('mouseout', () => {
        console.log(tv.id, "default");
        if (tv.id === "tv0") {
            document.getElementById("mtv0").style.background = "black";
        }
        if (tv.id === "tv1") {
            document.getElementById("mtv1").style.background = "black";
        }
        if (tv.id === "tv2") {
            document.getElementById("mtv2").style.background = "black";
        }
        if (tv.id === "tv3") {
            document.getElementById("mtv3").style.background = "black";
        }
        if (tv.id === "tv4") {
            document.getElementById("mtv4").style.background = "black";
        }
        if (tv.id === "tv5") {
            document.getElementById("mtv5").style.background = "black";
        }
    })
})