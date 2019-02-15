package server

import (
	"fmt"
	"path/filepath"

	"github.com/disintegration/imaging"
)

func resizeImageNative(src string, dest string) error {
	srcImage, err := imaging.Open(src, imaging.AutoOrientation(true))

	if err != nil {
		return err
	}

	srcImage = imaging.Fit(srcImage, 1500, 1500, imaging.NearestNeighbor)

	err = imaging.Save(srcImage, dest)

	if err != nil {
		return err
	}

	return nil
}

func resizeImages(fileList []string, n int, total int, srcFolder string, targetFolder string) error {

	//fmt.Printf("%s %s\n", srcFolder, targetFolder)

	for i := n; i < len(fileList); i += total {

		srcPath := srcFolder + "/" + filepath.Base(fileList[i])
		destPath := targetFolder + "/" + filepath.Base(fileList[i])
		//cmd := exec.Command("vipsthumbnail", srcPath, "-t", "-s", "1500x1500", "-o", destPath)
		//cmd := exec.Command("epeg", "-m 1500", srcPath, destPath)
		//cmd := exec.Command("convert", "-define", "jpeg:size=1500x1500", srcPath, "-auto-orient", "-thumbnail", "1500x1500>", destPath)
		//err := cmd.Run()

		// resize using go library
		err := resizeImageNative(srcPath, destPath)

		fmt.Printf("resizer %d: [%d]=%s -- %v\n", n, i, fileList[i], err)
		//fmt.Printf("convert %s -resize 1024x %s\n", srcPath, destPath)

		//time.Sleep(100 * time.Millisecond)
	}
	return nil
}
