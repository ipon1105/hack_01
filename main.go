package main

import (
	"fmt"
	"io"

	"os"
	"path/filepath"
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

func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func main() {

	f1, _ := os.Open("Датасет/Real/1__M_Left_index_finger.BMP")

	pixelArr1, _ := getPixels(f1)
	var bBranches, bEnds = findPoints(binarization(pixelArr1))
	bBranches, bEnds = delNoisePoint(bBranches, bEnds)
	f2, _ := os.Open("Датасет/Real/2__F_Left_index_finger.BMP")
	pixelArr2, _ := getPixels(f2)

	var (
		root  string
		files []string
		err   error
		i     int
		s     string
	)
	i = 0
	root = "Датасет/Real"
	files, err = FilePathWalkDir(root)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		//fmt.Println(file)
		f2, _ = os.Open(file)
		pixelArr2, _ = getPixels(f2)
		if specialPointCompare(bBranches, bEnds, binarization(pixelArr2)) == 100 {
			s = file
		}
	}
	//fmt.Println(MaxParallelism())
	fmt.Print(s)
	fmt.Print("\n")
}

// Поварачивать изображения

//D * DPI = PIXELS_RESOLUTION

//Сранение по узору
