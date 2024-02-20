package main

import (
	"data2video/common"
	"fmt"
	vidio "github.com/AlexEidt/Vidio"
	"github.com/disintegration/imaging"
	"image"
	"image/color"
	"os"
	"path/filepath"
)

func combineImages(rows int, cols int, savePath string, images []image.Image) {
	squareImage := imaging.New(cols*images[0].Bounds().Dx(), rows*images[0].Bounds().Dy(), color.NRGBA{0, 0, 0, 0})

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if (i*cols)+j >= len(images) {
				break
			}
			img := images[(i*cols)+j]
			squareImage = imaging.Paste(squareImage, img, image.Pt(j*img.Bounds().Dx(), i*img.Bounds().Dy()))
		}
	}
	err := imaging.Save(squareImage, savePath)
	common.IsErr(err, true)
}

func createFrames(imagesDir string, destDir string) {
	rows := 2
	cols := 3

	walker := 0
	indexer := 0
	var imagesArr []image.Image
	err := filepath.Walk(imagesDir, func(path string, info os.FileInfo, err error) error {
		common.IsErr(err, true)

		if !info.IsDir() {

			if walker >= rows*cols {
				resultPath := fmt.Sprintf(destDir+"frame-%09d.jpg", indexer)
				indexer = indexer + 1

				combineImages(rows, cols, resultPath, imagesArr)

				walker = 0
				imagesArr = nil
			}

			fmt.Println(path, info.Size())
			imageQr, err := imaging.Open(path)
			common.IsErr(err, true)

			imagesArr = append(imagesArr, imageQr)

			walker = walker + 1
			return nil

		}
		return nil
	})
	common.IsErr(err, true)

	if len(imagesArr) > 0 {
		resultPath := fmt.Sprintf(destDir+"frame-%09d.jpg", indexer)
		combineImages(rows, cols, resultPath, imagesArr)
	}
}

func main() {

	data, _ := os.ReadFile("./t.pdf")

	encoded := base45.Encode(data)

	chunkSize := 4296
	for i := 0; i <= len(encoded)/chunkSize; i++ {
		startBit := i * chunkSize
		endBit := (i + 1) * chunkSize

		if endBit > len(encoded) {
			endBit = len(encoded)
		}

		numString := fmt.Sprintf("%09d", i)
		name := "./assets/qr-base45" + numString + ".png"

		err := qrcode.WriteFile(string(encoded[startBit:endBit]), qrcode.Low, 540, name)
		if err != nil {
			fmt.Printf("could not generate QRCode: %v", err)
			return
		}
	}

	// -------------------------------------------------------------------

	createFrames("./assets", "./combines/")

	options := vidio.Options{FPS: 60}
	gif, _ := vidio.NewVideoWriter("output.mp4", 1620, 1080, &options)
	defer gif.Close()

	err := filepath.Walk("./combines", func(path string, info os.FileInfo, err error) error {
		common.IsErr(err, true)

		if !info.IsDir() {
			println(path)
			_, _, img, _ := vidio.Read(path)
			gif.Write(img)
			return nil
		}
		return nil
	})
	common.IsErr(err, true)



	// -------------------------------------------------------------------

}
