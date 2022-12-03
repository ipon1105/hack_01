package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

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

var wg sync.WaitGroup

// Функция для многопоточности или подругому говоря для горутин
func goroutina(branches []Coord, ends []Coord, start int, stop int, files []string) (int, float64) {
	defer wg.Done()
	t0 := time.Now()

	var (
		nowAccuracy  float64 = 0
		bestAccuracy float64 = 0
		bestFile     int     = 0
	)

	for i := start; i < stop && i < len(files); i++ {
		f, _ := os.Open(files[i])
		pixelArr, _ := getPixels(f)
		nowAccuracy = specialPointCompare(branches, ends, binarization(pixelArr))
		if bestAccuracy < nowAccuracy {
			bestAccuracy = nowAccuracy
			bestFile = i
		}
	}

	t1 := time.Now()
	fmt.Printf("accuracy: %f;\tworkTime: %v;\tfile: %s.\n", bestAccuracy, t1.Sub(t0), files[bestFile])
	return bestFile, bestAccuracy
}

func main() {

	THREADS := MaxParallelism()
	// THREADS := 1

	targetFile, _ := os.Open("Датасет/Real/1__M_Left_index_finger.BMP")

	pixelArr, _ := getPixels(targetFile)
	bImg := binarization(pixelArr)
	skeletonization(bImg)
	var bBranches, bEnds = findPoints(bImg)
	bBranches, bEnds = delNoisePoint(bBranches, bEnds)

	root := "Датасет/Altered/Altered-Hard"
	files, err := FilePathWalkDir(root)

	if false {
		grayImage := grayscale(pixelArr)
		sobelay(grayImage)
		//fmt.Println(grayImage)
		//fmt.Println()
		//gabor(1, 1, 1, 1, 1)
		return
	}

	if err != nil {
		panic(err)
	}

	length := int(math.Ceil(float64(len(files) / THREADS)))
	for i := 0; i < THREADS; i++ {
		wg.Add(1)

		if i == 0 {
			go goroutina(bBranches, bEnds, 0, length, files)
			continue
		}
		if i == length-1 {
			go goroutina(bBranches, bEnds, i*length+1, len(files)-1, files)
			continue
		}
		go goroutina(bBranches, bEnds, i*length+1, (i+1)*length, files)
	}
	wg.Wait()
	//time.Sleep(time.Minute)
}

// Поварачивать изображения

//D * DPI = PIXELS_RESOLUTION

//Сранение по узору
