(function () {

    var photos = [];
    var pos = 0;

    var mod = function (n, m) {
        return ((n % m) + m) % m;
    }

    var getID = function (id) {
        return document.getElementById(id);
    }

    var updateFooter = function () {
        getID('toolbar').innerHTML = "&nbsp;&nbsp;" + photos[pos] + " [" + (pos + 1) + "/" + photos.length + "]";
    }


    var displayPhoto = function () {
        getID("photo").src = "t/" + photos[pos];
        updateFooter();
    }

    var displayHiResPhoto = function () {
        getID("photo").src = "p/" + photos[pos];
        updateFooter();
    }


    var request = new XMLHttpRequest();
    request.open('GET', '/photos.json', true);
    request.onload = function () {
        if (request.status >= 200 && request.status < 400) {
            photos = JSON.parse(request.responseText);
            displayHiResPhoto();
        }
    };
    request.send();

    document.onkeydown = function (e) {
        loadHiRes = false;
        if (e.keyCode == 37) { // left
            pos = mod(pos - 1, photos.length);
            displayPhoto();
        } else if (e.keyCode == 39 || e.keyCode == 32) { // right
            pos = (pos + 1) % photos.length;
            displayPhoto();

        }
    };

    document.onkeyup = function (e) {
        displayHiResPhoto()
    }

}).call(this);