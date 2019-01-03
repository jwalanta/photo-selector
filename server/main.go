package server

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rwcarlsen/goexif/exif"
)

var photosFolder = "." // default to current folder
var verbose = true

func photosJSONHandler(w http.ResponseWriter, r *http.Request) {

	files, _ := filepath.Glob(photosFolder + "/*")
	jsonStr := ""

	for _, file := range files {
		ext := filepath.Ext(file)
		if ext == ".jpg" || ext == ".JPG" {
			jsonStr = jsonStr + "\"" + filepath.Base(file) + "\","
		}
	}

	if jsonStr != "" {
		// remove the trailing comma
		jsonStr = jsonStr[:len(jsonStr)-1]
	}
	jsonStr = "[" + jsonStr + "]"

	w.Write([]byte(jsonStr))
}

func thumbHandler(w http.ResponseWriter, r *http.Request) {

	path := strings.Replace(r.URL.Path, "/t/", photosFolder+"/", 1)

	if verbose {
		fmt.Print("Request: " + r.URL.Path)
	}

	f, err := os.Open(path)

	if err == nil {
		exifData, err := exif.Decode(f)

		if err == nil {
			thumbnailByte, err := exifData.JpegThumbnail()

			if err == nil {
				fmt.Println(". Serving EXIF thumbnail")
				w.Write(thumbnailByte)
				return
			}

		}

	}

	fmt.Println(". EXIF thumbnail not found, serving: " + path)
	http.ServeFile(w, r, path)

}

func photoHandler(w http.ResponseWriter, r *http.Request) {

	path := strings.Replace(r.URL.Path, "/p/", photosFolder+"/", 1)

	if verbose {
		fmt.Println("Request: " + r.URL.Path + ". Serving: " + path)
	}

	http.ServeFile(w, r, path)

}

// Run -- Start the web server
func Run() {

	portPtr := flag.Int("p", 8080, "Port")
	verbosePtr := flag.Bool("v", true, "Verbose")

	flag.Parse()

	if flag.NArg() > 0 {
		photosFolder = flag.Arg(0)
	}

	if *verbosePtr {
		verbose = true
	}

	// serve web
	// generated using:
	// go-bindata-assetfs web/...
	http.Handle("/", http.FileServer(assetFS()))
	// staticFs := http.FileServer(http.Dir("web"))
	// http.Handle("/", staticFs)

	// serve thumbnail
	http.HandleFunc("/t/", thumbHandler)

	// serve photos
	http.HandleFunc("/p/", photoHandler)

	// generate json
	http.HandleFunc("/photos.json", photosJSONHandler)

	fmt.Println("Running server at http://localhost:" + strconv.Itoa(*portPtr))

	if err := http.ListenAndServe(":"+strconv.Itoa(*portPtr), nil); err != nil {
		panic(err)
	}

}
