    document.querySelectorAll('.TV').forEach((tv) => {
        tv.addEventListener('click', () => {
            console.log(tv.id, "clicked");
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
        });
    });

    document.querySelectorAll('.TV').forEach((tv) => {
        tv.addEventListener('mouseover', () => {
            console.log(tv.id, "hover");
            if(tv.id === "tv0") {
                document.getElementById("mtv0").style.background = "red";
            }
            if(tv.id === "tv1") {
                document.getElementById("mtv1").style.background = "red";
            }
            if(tv.id === "tv2") {
                document.getElementById("mtv2").style.background = "red";
            }
            if(tv.id === "tv3") {
                document.getElementById("mtv3").style.background = "red";
            }
            if(tv.id === "tv4") {
                document.getElementById("mtv4").style.background = "red";
            }
            if(tv.id === "tv5") {
                document.getElementById("mtv5").style.background = "red";
            }
        })
    })

    document.querySelectorAll('.TV').forEach((tv) => {
        tv.addEventListener('mouseout', () => {
            console.log(tv.id, "default");
            if(tv.id === "tv0") {
                document.getElementById("mtv0").style.background = "black";
            }
            if(tv.id === "tv1") {
                document.getElementById("mtv1").style.background = "black";
            }
            if(tv.id === "tv2") {
                document.getElementById("mtv2").style.background = "black";
            }
            if(tv.id === "tv3") {
                document.getElementById("mtv3").style.background = "black";
            }
            if(tv.id === "tv4") {
                document.getElementById("mtv4").style.background = "black";
            }
            if(tv.id === "tv5") {
                document.getElementById("mtv5").style.background = "black";
            }
        })
    })