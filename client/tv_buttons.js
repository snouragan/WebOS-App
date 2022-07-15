function toggleEnable(id) {
    var element = document.getElementById(id);
    if (element.classList.contains("lime-btn")) {
        element.classList.remove("lime-btn");
        element.classList.add("red-btn");
    }
    else {
        element.classList.add("lime-btn");
        element.classList.remove("red-btn");
    }
}