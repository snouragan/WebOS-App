    document.querySelectorAll('.TV').forEach((tv) => {
        tv.addEventListener('click', () => {
            console.log(tv.id);
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
