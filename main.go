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

// Функция бинаризации
func binarization2(img [][]Pixel) [][]int {
	var bImg [][]int

	for _, row := range img {
		var p []int
		for _, col := range row {
			p = append(p, int(float64(col.R)*BINARY_RATION_R+float64(col.G)*BINARY_RATION_G+float64(col.B)*BINARY_RATION_B))
		}

		bImg = append(bImg, p)
	}

	return bImg
}

// Функция бинаризации
func binarization3(img [][]int) {
	for a, row := range img {
		for b, col := range row {
			if col > 128 {
				img[a][b] = 1 // Чёрный
			} else {
				img[a][b] = 1 // Белый
			}

		}

	}

}

func main() {

	f1, _ := os.Open("Датасет/Real/1__M_Left_index_finger.BMP")
	//f2, _ := os.Open("Датасет/Real/1__M_Left_index_finger.BMP")

	pixelArr1, _ := getPixels(f1)
	//pixelArr2, _ := getPixels(f2)

	bImg := binarization2(pixelArr1)
	fmt.Println("bImg2")
	fmt.Println(bImg)
	gabor(1, 2, 3, 4, 5)
	//fmt.Println("gabor")
	//gaboraImagination(bImg)
	//fmt.Println("bImg3")
	//binarization3(bImg)
	//fmt.Println(bImg)

	//bImg := binarization(pixelArr1)
	//skeletonization(bImg)
	//var branches, ends = findPoints(bImg)
	//fmt.Println(orientation(branches, ends, len(bImg[0]), len(bImg)))

	//fmt.Println(specialPointCompare(binarization(pixelArr1), binarization(pixelArr2)))
}

// Поварачивать изображения

//D * DPI = PIXELS_RESOLUTION

//Сранение по узору
