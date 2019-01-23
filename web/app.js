(function () {

    var photos = [];
    var pos = 0;
    var buckets = [[],[],[],[],[],[],[],[],[],[]]

    var mod = function (n, m) {
        return ((n % m) + m) % m;
    }

    var getID = function (id) {
        return document.getElementById(id);
    }

    var updateFooter = function () {
        getID('info').innerHTML = photos[pos] + " [" + (pos + 1) + "/" + photos.length + "]";
    }


    var displayPhoto = function () {
        getID("photo").src = "t/" + photos[pos];
        updateFooter();
        updateBucketDisplay()
    }

    var displayHiResPhoto = function () {
        getID("photo").src = "p/" + photos[pos];
        updateFooter();
        
    }

    var updateBucketDisplay = function (){
        for (var i=0;i<10;i++){
            var b = getID("bucket"+i)
            if (buckets[i].indexOf(photos[pos]) > -1){
                if (!b.classList.contains("selected"))
                    b.classList.add("selected")

            }
            else{
                b.classList.remove("selected")

            }
        }
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
        if (e.keyCode == 37) { // left
            pos = mod(pos - 1, photos.length);
            displayPhoto();
        } else if (e.keyCode == 39 || e.keyCode == 32) { // right
            pos = (pos + 1) % photos.length;
            displayPhoto();
        } else if (e.keyCode >=48 && e.keyCode <=57){
            var i = e.keyCode - 48;

            var index = buckets[i].indexOf(photos[pos]);

            if (index > -1){
                buckets[i].splice(index,1);
            }
            else{
                buckets[i].push(photos[pos]);
            }

            updateBucketDisplay()

        }
    };

    document.onkeyup = function (e) {
        displayHiResPhoto()
    }

    var bucketStr = "";
    
    for (var i=1; i<=10;i++){
        bucketStr += "<span id='bucket"+(i%10)+"'>"+(i%10)+"</span>";
    }
    getID("buckets").innerHTML = bucketStr;



}).call(this);