package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/rwcarlsen/goexif/exif"
)

var photosFolder = "." // default to current folder
var verbose = true
var resizeFolder = ""

var labelPath, _ = homedir.Expand("~/.phs.labels")

var validExtensions = []string{"jpg", "cr2", "dng", "raw", "raf"}

func isValidExtension(ext string) bool {

	for _, e := range validExtensions {
		if e == ext {
			return true
		}
	}
	return false

}

func labelsJSONHandler(w http.ResponseWriter, r *http.Request) {

	if fileExists(labelPath) {
		http.ServeFile(w, r, labelPath)
	} else {
		labels := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
		labelsJSON, _ := json.Marshal(labels)
		w.Write(labelsJSON)
	}

}

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

	resizedPath := strings.Replace(r.URL.Path, "/t/", resizeFolder+"/", 1)

	if _, err := os.Stat(resizedPath); os.IsNotExist(err) {
		// resized file doesnt exist
		// generate from exif data

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

	} else {
		http.ServeFile(w, r, resizedPath)
	}

}

func photoHandler(w http.ResponseWriter, r *http.Request) {

	path := strings.Replace(r.URL.Path, "/p/", resizeFolder+"/", 1)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		// resized file doesnt exist
		path = strings.Replace(r.URL.Path, "/p/", photosFolder+"/", 1)
	}

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

	absPhotosFolder, _ := filepath.Abs(photosFolder)
	resizeFolder = "/tmp/photo-selector-" + getMD5(absPhotosFolder)

	os.MkdirAll(resizeFolder, 0777)

	files, _ := filepath.Glob(photosFolder + "/*")
	totalCPUs := runtime.NumCPU() / 2
	for i := 0; i < totalCPUs; i++ {
		go resizeImages(files, i, totalCPUs, absPhotosFolder, resizeFolder)
	}

	// serve web
	// generated using:
	// go-bindata-assetfs web/...
	// http.Handle("/", http.FileServer(assetFS()))
	staticFs := http.FileServer(http.Dir("web"))
	http.Handle("/", staticFs)

	// serve thumbnail
	http.HandleFunc("/t/", thumbHandler)

	// serve photos
	http.HandleFunc("/p/", photoHandler)

	// generate json
	http.HandleFunc("/photos.json", photosJSONHandler)

	// labels json
	http.HandleFunc("/labels.json", labelsJSONHandler)

	fmt.Println("Running server at http://localhost:" + strconv.Itoa(*portPtr))

	if err := http.ListenAndServe(":"+strconv.Itoa(*portPtr), nil); err != nil {
		panic(err)
	}

}
