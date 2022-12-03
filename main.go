package main

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/sergeymakinen/go-bmp"
)

func getPixels(file io.Reader) ([][]Pixel, error) {
	img, _ := bmp.Decode(file)

	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var pixels [][]Pixel
	for y := 0; y < height; y++ {
		var row []Pixel
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels, nil
}

// img.At(x, y).RGBA() returns four uint32 values; we want a Pixel
func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

// Pixel struct example
type Pixel struct {
	R int
	G int
	B int
	A int
}

func MaxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}

func main() {
	f1, _ := os.Open("Датасет/Real/1__M_Left_index_finger.BMP")
	f2, _ := os.Open("Датасет/Real/2__F_Left_index_finger.BMP")

	pixelArr1, _ := getPixels(f1)
	pixelArr2, _ := getPixels(f2)

	fmt.Println(specialPointCompare(binarization(pixelArr1), binarization(pixelArr2)))
}

//D * DPI = PIXELS_RESOLUTION

//Сранение по узору
