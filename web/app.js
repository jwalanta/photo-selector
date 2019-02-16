var photos = [];
var pos = 0;
var buckets = [[],[],[],[],[],[],[],[],[],[]];
var bucketLabels = ["0", "1", "2", "3", "4", "5", "6", "7", "8", "9"]

var mod = function (n, m) {
    return ((n % m) + m) % m;
}

var getID = function (id) {
    return document.getElementById(id);
}

var updateFooter = function () {
    getID('info').innerHTML = photos[pos] + " [" + (pos + 1) + "/" + photos.length + "]";
}

var renderBucket = function(){
    var bucketStr = "";
    for (var i = 1; i <= 10; i++) {
        var n = i % 10
        bucketStr += "<span class='bucket' id='bucket" + n + "' onclick='displayBucket(" + n + ")'>";
        bucketStr += bucketLabels[n] == n ? n : n + ": " + bucketLabels[n];
        bucketStr += "</span>";
    }
    getID("buckets").innerHTML = bucketStr;
    updateBucketDisplay();
}

var displayPhoto = function () {
    getID("photo").src = "t/" + photos[pos];
    updateFooter();
    updateBucketDisplay()
}

var displayHiResPhoto = function () {

    var p = getID("photo")

    // dont update if the existing photo is big enough
    if (p.naturalHeight < 1000 && p.naturalWidth < 1000) {
        p.src = "p/" + photos[pos];
        updateFooter();
    }

}

var toggleBucket = function (i) {
    var index = buckets[i].indexOf(photos[pos]);

    if (index > -1) {
        buckets[i].splice(index, 1);
    } else {
        buckets[i].push(photos[pos]);
    }

    updateBucketDisplay()
}

var updateBucketDisplay = function () {
    for (var i = 0; i < 10; i++) {
        var b = getID("bucket" + i)
        if (buckets[i].indexOf(photos[pos]) > -1) {
            if (!b.classList.contains("selected"))
                b.classList.add("selected")

        } else {
            b.classList.remove("selected")

        }
    }
}

var displayBucket = function(n){

    getID("dialog").innerHTML = buckets[n].join(" ");

    // show
    getID("dialogbg").style.display = "block";
}


getID("dialogbg").addEventListener("click", function(e){
    if (e.target.id == "dialogbg"){
        getID("dialogbg").style.display = "none";
    }
});

(function(){
    var request = new XMLHttpRequest();
    request.open('GET', '/photos.json', true);
    request.onload = function () {
        if (request.status >= 200 && request.status < 400) {
            photos = JSON.parse(request.responseText);
            displayHiResPhoto();
        }
    };
    request.send();
})();

document.onkeydown = function (e) {
    if (e.keyCode == 37) { // left
        pos = mod(pos - 1, photos.length);
        displayPhoto();
    } else if (e.keyCode == 39 || e.keyCode == 32) { // right
        pos = (pos + 1) % photos.length;
        displayPhoto();
    } else if (e.keyCode >= 48 && e.keyCode <= 57) {
        var i = e.keyCode - 48;
        toggleBucket(i)
    }
};

document.onkeyup = function (e) {
    displayHiResPhoto()
};

(function(){
    var request = new XMLHttpRequest();
    request.open('GET', '/labels.json', true);
    request.onload = function () {
        if (request.status >= 200 && request.status < 400) {
            bucketLabels = JSON.parse(request.responseText);
        }
        renderBucket()
    };
    request.send();
})();

(function(){
    var request = new XMLHttpRequest();
    request.open('GET', '/selections.json', true);
    request.onload = function () {
        if (request.status >= 200 && request.status < 400) {
            buckets = JSON.parse(request.responseText);
        }
        renderBucket()
    };
    request.send();
})();


var saveBuckets = function(){
    var request = new XMLHttpRequest();
    request.open('POST', '/selections.json', true);
    request.setRequestHeader("Content-Type", "application/json");
    // request.onload = function () {
    //     var s = getID("save-btn");
    //     if (request.status == 200) {
    //         s.classList.add("state-success");
    //     }
    //     else{
    //         s.classList.add("state-error");
    //     }
    // };
    
    var data = JSON.stringify(buckets);
    request.send(data);
}