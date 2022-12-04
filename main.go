package main

import (
	"fmt"
	"math"
	"os"
	"sync"
	"time"
)

var minAccuracy float64 = 0.90000

// Топ лучших в потоке
const queue = 3

// Подробно
var moreInfo bool = false

var wg sync.WaitGroup
var mt sync.Mutex

type Answer struct {
	Accuracy    float64
	WorkingTime time.Duration
	Filename    string
}

type Info struct {
	Queue       []Answer
	WorkingTime time.Duration
}

// Функция для многопоточности или подругому говоря для горутин
func goroutina(branches []Coord, ends []Coord, start int, stop int, files []string) Info {
	var (
		nowAccuracy float64  = 0
		writeIndex  int      = 0
		answerArr   []Answer = make([]Answer, queue)
	)

	for i, _ := range answerArr {
		answerArr[i] = Answer{0.00, time.Now().Sub(time.Now()), ""}
	}

	min := func() int {
		minIndex := 0
		for i, el := range answerArr {
			if answerArr[minIndex].Accuracy > el.Accuracy {
				minIndex = i
			}
		}
		return minIndex
	}

	t0 := time.Now()
	for i := start; i < stop && i < len(files); i++ {
		writeIndex = -1

		tt0 := time.Now()
		f, _ := os.Open(files[i])
		pixelArr, _ := getPixels(f)
		nowAccuracy = specialPointCompare(branches, ends, binarization(pixelArr))
		tt1 := time.Now()

		m := min()
		if nowAccuracy > answerArr[m].Accuracy {
			writeIndex = m
		} else {
			continue
		}

		answerArr[writeIndex] = Answer{nowAccuracy, tt1.Sub(tt0), files[i]}
	}

	t1 := time.Now()

	return Info{answerArr[:], t1.Sub(t0)}
}

func fileToFileCompare(target string, compare string) float64 {
	targetFile, _ := os.Open(target)
	compareFile, _ := os.Open(compare)

	pixelTarget, _ := getPixels(targetFile)
	pixelCompare, _ := getPixels(compareFile)

	bImgTarget := binarization(pixelTarget)
	skeletonization(bImgTarget)
	var targetBranches, targetEnds = findPoints(bImgTarget)
	targetBranches, targetEnds = delNoisePoint(targetBranches, targetEnds)

	return specialPointCompare(targetBranches, targetEnds, binarization(pixelCompare))
}

func fileToDirCompare(filename string, serchdir string, threads int) []Info {

	targetFile, _ := os.Open(filename)

	pixelArr, _ := getPixels(targetFile)
	bImg := binarization(pixelArr)
	skeletonization(bImg)
	var bBranches, bEnds = findPoints(bImg)
	bBranches, bEnds = delNoisePoint(bBranches, bEnds)

	root := serchdir
	files, err := FilePathWalkDir(root)

	if err != nil {
		panic(err)
	}

	var results []Info
	length := int(math.Ceil(float64(len(files) / threads)))

	f := func(start int, stop int) {
		r := goroutina(bBranches, bEnds, start, stop, files)

		mt.Lock()
		results = append(results, r)
		mt.Unlock()

		defer wg.Done()
	}

	for i := 0; ; i += length {
		wg.Add(1)
		if i+length >= len(files) {
			go f(i, len(files)-1)
			break
		}

		go f(i, i+length)
	}
	wg.Wait()

	return results
}

func dirToDirCompare(dirTarget string, dirCompare string, threads int) [][]Info {
	targetFiles, err := FilePathWalkDir(dirTarget)
	if err != nil {
		panic(err)
	}

	var infoList [][]Info
	for _, target := range targetFiles {
		infoList = append(infoList, fileToDirCompare(target, dirCompare, threads))
	}
	return infoList
}

func main() {
	//Начальные параметры

	// Файл для сравнения
	filename := "Датасет/Real/1__M_Left_middle_finger.BMP"

	// Папка для сравнения
	dir := "Датасет/Altered/Altered-Easy"

	// Количество потоков
	threads := MaxParallelism()
	// threads := 3

	// Коэффициенты сравнения
	DEL_RANGE = 4
	FIND_OFFSET = 5

	results := fileToDirCompare(filename, dir, threads)

	for _, n := range results {
		for _, m := range n.Queue {
			if m.Accuracy > minAccuracy {
				fmt.Printf("file %s like %s as %f\n", filename, m.Filename, m.Accuracy)
				if moreInfo {
					fmt.Printf("Accuracy := %f; \t WorkingTime := %v.\n", m.Accuracy, m.WorkingTime)
				}
			}

		}
	}
}
