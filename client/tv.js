document.querySelectorAll('.TV').forEach((tv) => {
    tv.addEventListener('click', () => {
        console.log(tv.id);
        if (tv.classList.contains('clicked')) {
            tv.classList.remove('clicked');
            tv.classList.add('floating');
        } else {
            tv.classList.add('clicked');
            tv.classList.remove('floating');
            document.querySelectorAll('.TV').forEach((t) => {
                if (t.id !== tv.id) {
                    t.classList.remove('clicked');
                    t.classList.add('floating');
                    
                }
            })
        }
    });
});
