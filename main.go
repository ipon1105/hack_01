package main

import (
	"fmt"
	"io"
	"os"

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

func main() {
	f1, _ := os.Open("Датасет/Real/1__M_Left_index_finger.BMP")
	f2, _ := os.Open("Датасет/Real/600__M_Right_little_finger.BMP")

	pixelArr1, _ := getPixels(f1)
	pixelArr2, _ := getPixels(f2)

	fmt.Println(specialPointCompare(binaryzation(pixelArr1), binaryzation(pixelArr2)))
}

//D * DPI = PIXELS_RESOLUTION

//Сранение по узору

//Корреляционное сравнение
//Высокая трудоёмкость

//Сравнение по особым точкам
func specialPointCompare(bImg [][]int, cImg [][]int) (int, int) {
	skeletization(bImg)
	skeletization(cImg)

	var bBranches, bEnds = findPoints(bImg)
	bBranches, bEnds = delNoisePoint(bBranches, bEnds)

	var cBranches, cEnds = findPoints(cImg)
	cBranches, cEnds = delNoisePoint(cBranches, cEnds)

	fmt.Println("logs")
	fmt.Println("\tbBranches = ", bBranches)
	fmt.Println("\tbEnds = ", bEnds)
	fmt.Println("\tcBranches = ", cBranches)
	fmt.Println("\tcEnds = ", cEnds)

	return matchingPoints([][][]int{bBranches, bEnds}, [][][]int{cBranches, cEnds})
}

//Функция сравнения
func matchingPoints(origin [][][]int, target [][][]int) (int, int) {
	var all int = 0
	var match int = 0

	//Сравниваем ветвления
	var originBranches [][]int = origin[0]
	for i := 0; i < len(originBranches); i++ {
		var widthRange = [...]int{originBranches[i][0] - 15, originBranches[i][0] + 15}
		var heightRange = [...]int{originBranches[i][1] - 15, originBranches[i][1] + 15}
		var targetBranches [][]int = target[0]

		all++
		for j := 0; j < len(targetBranches); j++ {

			if targetBranches[j][0] >= widthRange[0] && targetBranches[j][0] <= widthRange[1] && targetBranches[j][1] >= heightRange[0] && targetBranches[j][1] <= heightRange[1] {
				match++
				break
			}
		}
	}

	//Сравниваем ветвления
	var originEnds [][]int = origin[1]
	for i := 0; i < len(originEnds); i++ {
		var widthRange = [...]int{originEnds[i][0] - 15, originEnds[i][0] + 15}
		var heightRange = [...]int{originEnds[i][1] - 15, originEnds[i][1] + 15}
		var targetEnds [][]int = target[1]

		all++
		for j := 0; j < len(targetEnds); j++ {

			if targetEnds[j][0] >= widthRange[0] && targetEnds[j][0] <= widthRange[1] && targetEnds[j][1] >= heightRange[0] && targetEnds[j][1] <= heightRange[1] {
				match++
				break
			}
		}
	}

	return match, all
}

//Удаляем повторения
func delNoisePoint(branchPoints [][]int, endPoints [][]int) ([][]int, [][]int) {
	var tmp, tmp2 [][]int

	for i := 0; i < len(endPoints); i++ {
		var widthRange = [...]int{endPoints[i][0] - 5, endPoints[i][0] + 5}
		var heightRange = [...]int{endPoints[i][1] - 5, endPoints[i][1] + 5}
		for j := 0; j < len(branchPoints); j++ {

			if branchPoints[j][0] >= widthRange[0] && branchPoints[j][0] <= widthRange[1] && branchPoints[j][1] >= heightRange[0] && branchPoints[j][1] <= heightRange[1] {
				tmp = append(tmp, endPoints[i])
				tmp2 = append(tmp2, branchPoints[j])
			}
		}
	}

	return removeDouble(branchPoints, tmp2), removeDouble(endPoints, tmp)
}

//Возвращает список элементов, у которых нет одинакового в другом  списке
func removeDouble(x [][]int, y [][]int) [][]int {
	var z [][]int
	for hx := 0; hx < len(x); hx++ {
		c := true
		for hy := 0; hy < len(y); hy++ {
			if intArrayEquals(x[hx], y[hy]) {
				c = false
			}
		}
		if c {
			z = append(z, x[hx])
		}
	}
	for hy := 0; hy < len(y); hy++ {
		c := true
		for hx := 0; hx < len(x); hx++ {
			if intArrayEquals(y[hy], x[hx]) {
				c = false
			}
		}
		if c {
			z = append(z, y[hy])
		}
	}
	return z
}

// Функция подсчёта количество чёрных точек в округе
// TODO: Здесь возможна ошибка при подсчёте точек на краях массива.
func getBlackArround(img [][]int, x int, y int) int {
	if y-1 < 0 || x-1 < 0 || y+1 >= len(img) || x+1 >= len(img[0]) {
		return 0
	}
	var c int = 0

	for j := y - 1; j < y+1; j++ {
		for i := x - 1; i < x+1; i++ {
			if img[j][i] == 1 {
				c++
			}
		}
	}
	return c
}

// Функция составления списка особых точек
func findPoints(img [][]int) ([][]int, [][]int) {
	var branchPoints [][]int
	var endPoints [][]int

	for h, vh := range img {
		for w := range vh {
			if img[h][w] == 0 {
				var tmp int = getBlackArround(img, w, h)
				if tmp == 1 {
					var arr = []int{w, h}
					endPoints = append(endPoints, arr)
				}
				if tmp == 3 {
					var arr = []int{w, h}
					branchPoints = append(branchPoints, arr)
				}
			}
		}
	}

	return branchPoints, endPoints
}

//Функция бинаризации
//TODO: Доделать функцию бинаризации
func binaryzation(img [][]Pixel) [][]int {
	var bImg [][]int

	for i := 0; i < len(img); i++ {
		var p []int
		for j := 0; j < len(img[i]); j++ {
			tmp := int(float64(img[i][j].R)*0.3 + float64(img[i][j].G)*0.59 + float64(img[i][j].B)*0.11)

			if tmp > 128 { // Чёрный
				tmp = 1
			} else { // Белый
				tmp = 0
			}

			p = append(p, tmp)
		}

		bImg = append(bImg, p)
	}
	return bImg
}

//Скелетизация
func skeletization(img [][]int) {
	var count int = 1
	for count != 0 {
		count = deleteMain(img)
		if count > 0 {
			deleteNoise(img)
		}
	}
}

//Удаление пикселя по набору шумов
func deleteNoise(img [][]int) {
	for h := 1; h < len(img)-1; h++ {
		for w := 1; w < len(img[h])-1; w++ {
			if img[h][w] == 0 && fringe(getTripleVector(img, w, h)) {
				img[h][w] = 1
			}
		}
	}
}

//Удаление пикселя по основному набору
func deleteMain(img [][]int) int {
	var count int = 0
	for h := 1; h < len(img)-1; h++ {
		for w := 1; w < len(img[h])-1; w++ {
			if img[h][w] == 0 && check(getTripleVector(img, w, h)) {
				img[h][w] = 1
				count++
			}
		}
	}
	return count
}

// Проверка на удаление
func getTripleVector(img [][]int, x int, y int) []int {
	if y-1 < 0 || x-1 < 0 || y+1 >= len(img) || x+1 >= len(img[0]) {
		return nil
	}

	var a = make([]int, 9)
	for j := y - 1; j < y+1; j++ {
		for i := x - 1; i < x+1; i++ {
			a = append(a, img[j][i])
		}
	}
	return a
}

// Функция сравения шаблонов с вектором
func check(a []int) bool {
	if a == nil {
		return false
	}

	//4 шаблона
	if templateCompare(a, []int{1, 1, 0, 0, 1, 0}, []int{1, 2, 3, 4, 5, 7}) {
		return true
	}
	if templateCompare(a, []int{1, 1, 1, 0, 0, 0}, []int{0, 1, 3, 4, 5, 7}) {
		return true
	}
	if templateCompare(a, []int{0, 1, 0, 0, 1, 1}, []int{1, 3, 4, 5, 6, 7}) {
		return true
	}
	if templateCompare(a, []int{0, 0, 0, 1, 1, 1}, []int{1, 3, 4, 5, 7, 8}) {
		return true
	}

	//4 шаблона
	if templateCompare(a, []int{1, 1, 1, 0, 0, 0, 0}, []int{0, 1, 2, 3, 4, 5, 7}) {
		return true
	}
	if templateCompare(a, []int{1, 0, 1, 0, 0, 1, 0}, []int{0, 1, 3, 4, 5, 6, 7}) {
		return true
	}
	if templateCompare(a, []int{0, 0, 0, 0, 1, 1, 1}, []int{1, 3, 4, 5, 6, 7, 8}) {
		return true
	}
	if templateCompare(a, []int{0, 1, 0, 0, 1, 0, 1}, []int{1, 2, 3, 4, 5, 7, 8}) {
		return true
	}

	return false
}

// Функция сравения шаблонов с вектором
func fringe(a []int) bool {
	if a == nil {
		return false
	}

	var matrix = make([][]int, 13)

	matrix[0] = []int{1, 1, 1, 1, 0, 1, 1, 1, 1}
	matrix[1] = []int{1, 1, 1, 1, 0, 1, 1, 0, 0}
	matrix[2] = []int{1, 1, 1, 0, 0, 1, 0, 1, 1}
	matrix[3] = []int{0, 0, 1, 1, 0, 1, 1, 1, 1}
	matrix[4] = []int{1, 1, 0, 1, 0, 0, 1, 1, 1}
	matrix[5] = []int{1, 1, 1, 1, 0, 1, 0, 0, 1}
	matrix[6] = []int{0, 1, 1, 0, 0, 1, 1, 1, 1}
	matrix[7] = []int{1, 0, 0, 1, 0, 1, 1, 1, 1}
	matrix[8] = []int{1, 1, 1, 1, 0, 0, 1, 1, 0}
	matrix[9] = []int{1, 1, 1, 1, 0, 1, 0, 0, 0}
	matrix[10] = []int{0, 1, 1, 0, 0, 1, 0, 1, 1}
	matrix[11] = []int{0, 0, 0, 1, 0, 1, 1, 1, 1}
	matrix[12] = []int{1, 1, 0, 1, 0, 0, 1, 1, 0}

	for i := 0; i < len(matrix); i++ {
		if intArrayEquals(a, matrix[i]) {
			return true
		}
	}
	return false
}

// Функция сранения двух массивов
func intArrayEquals(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

//Функция сравнения шаблона c вектором
func templateCompare(vector []int, template []int, allowed []int) bool {
	var count = 0
	for i := 0; i < len(vector); i++ {
		if i != allowed[count] {
			continue
		}

		if vector[i] != template[count] {
			return false
		}
		count++
	}

	return true
}
