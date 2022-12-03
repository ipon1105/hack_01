package main

// Константы для бинаризации изображения
const BINARY_RATION_R = float64(0.3)
const BINARY_RATION_G = float64(0.59)
const BINARY_RATION_B = float64(0.11)

// Константы для удаления повторений в радиусе
const DEL_RANGE = int(5)

// Смещения для поиска точек
const FIND_OFFSET = int(4)

// Структура содержащая в себе два числа
type Coord struct {
	X int
	Y int
}

// Функция сравнивания двух Coord
func isCoordsEquals(a Coord, b Coord) bool {
	if a.X != b.X {
		return false
	}
	if a.Y != b.Y {
		return false
	}
	return true
}

////////////////////
//Сравнение по особым точкам

func specialPointCompare(bBranches []Coord, bEnds []Coord, cImg [][]int) float64 {

	skeletonization(cImg)
	var cBranches, cEnds = findPoints(cImg)
	cBranches, cEnds = delNoisePoint(cBranches, cEnds)

	return matchingPoints([][]Coord{bBranches, bEnds}, [][]Coord{cBranches, cEnds})
}

// Функция сравнения
func matchingPoints(origin [][]Coord, target [][]Coord) float64 {
	var all int = 0
	var match int = 0

	// Сравниваем ветвления
	var originBranches []Coord = origin[0]
	for _, oCoord := range originBranches {
		var widthRange = Coord{oCoord.X - FIND_OFFSET, oCoord.X + FIND_OFFSET}
		var heightRange = Coord{oCoord.Y - FIND_OFFSET, oCoord.Y + FIND_OFFSET}
		var targetBranches []Coord = target[0]

		all++
		for _, tCoord := range targetBranches {
			if tCoord.X >= widthRange.X && tCoord.X <= widthRange.Y && tCoord.Y >= heightRange.X && tCoord.Y <= heightRange.Y {
				match++
				break
			}
		}
	}

	//Сравниваем ветвления
	var originEnds []Coord = origin[1]
	for _, oCoord := range originEnds {
		var widthRange = Coord{oCoord.X - FIND_OFFSET, oCoord.X + FIND_OFFSET}
		var heightRange = Coord{oCoord.Y - FIND_OFFSET, oCoord.Y + FIND_OFFSET}
		var targetEnds []Coord = target[1]

		all++
		for _, tCoord := range targetEnds {
			if tCoord.X >= widthRange.X && tCoord.X <= widthRange.Y && tCoord.Y >= heightRange.X && tCoord.Y <= heightRange.Y {
				match++
				break
			}
		}
	}

	return (float64(match) / float64(all))
}

// Удаляем повторения
func delNoisePoint(branchPoints []Coord, endPoints []Coord) ([]Coord, []Coord) {
	var branchList, endList []Coord

	for i := 0; i < len(endPoints); i++ {
		var widthRange = Coord{endPoints[i].X - DEL_RANGE, endPoints[i].X + DEL_RANGE}
		var heightRange = Coord{endPoints[i].Y - DEL_RANGE, endPoints[i].Y + DEL_RANGE}
		for j := 0; j < len(branchPoints); j++ {

			if branchPoints[j].X >= widthRange.X && branchPoints[j].X <= widthRange.Y && branchPoints[j].Y >= heightRange.X && branchPoints[j].Y <= heightRange.Y {
				branchList = append(branchList, endPoints[i])
				endList = append(endList, branchPoints[j])
			}
		}
	}

	return removeDouble(branchPoints, endList), removeDouble(endPoints, branchList)
}

// Возвращает список элементов, у которых нет одинакового в другом  списке
func removeDouble(points []Coord, compareList []Coord) []Coord {
	var z []Coord
	for _, pEl := range points {
		c := true
		for _, cEl := range compareList {
			if isCoordsEquals(pEl, cEl) {
				c = false
			}
		}
		if c {
			z = append(z, pEl)
		}
	}
	for _, cEl := range compareList {
		c := true
		for _, pEl := range points {
			if isCoordsEquals(cEl, pEl) {
				c = false
			}
		}
		if c {
			z = append(z, cEl)
		}
	}
	return z
}

// Функция подсчёта количество чёрных точек в округе
func getBlackArround(img [][]int, x int, y int) int {
	var c int = 0

	for j := y - 1; j < y+1; j++ {
		for i := x - 1; i < x+1; i++ {
			if i < 0 || i >= len(img) || j < 0 || j >= len(img[0]) {
				continue
			}

			if img[j][i] == 1 {
				c++
			}
		}
	}
	return c
}

// Функция составления списка особых точек
func findPoints(img [][]int) ([]Coord, []Coord) {
	var branchPoints []Coord
	var endPoints []Coord

	for h, vh := range img {
		for w, vw := range vh {
			if vw == 0 {
				var tmp int = getBlackArround(img, w, h)
				if tmp == 1 {
					endPoints = append(endPoints, Coord{w, h})
				}
				if tmp == 3 {
					branchPoints = append(branchPoints, Coord{w, h})
				}
			}
		}
	}

	return branchPoints, endPoints
}

// Функция бинаризации
func binarization(img [][]Pixel) [][]int {
	var bImg [][]int

	for _, row := range img {
		var p []int
		for _, col := range row {
			if col.A == 0 {
				p = append(p, 0)
				continue
			}

			pixel := int(float64(col.R)*BINARY_RATION_R + float64(col.G)*BINARY_RATION_G + float64(col.B)*BINARY_RATION_B)

			if pixel > 128 {
				pixel = 1 // Чёрный
			} else {
				pixel = 0 // Белый
			}

			p = append(p, pixel)
		}

		bImg = append(bImg, p)
	}

	return bImg
}

// Скелетизация
func skeletonization(img [][]int) {
	var count int = 1
	for count != 0 {
		count = deleteMain(img)
		if count > 0 {
			deleteNoise(img)
		}
	}
}

// Удаление пикселя по набору шумов
func deleteNoise(img [][]int) {
	for r := 1; r < len(img)-1; r++ {
		for c := 1; c < len(img[r])-1; c++ {
			if img[r][c] == 0 && fringe(getTripleVector(img, r, c)) {
				img[r][c] = 1
			}
		}
	}
}

// Удаление пикселя по основному набору
func deleteMain(img [][]int) int {
	var count int = 0
	for r := 1; r < len(img)-1; r++ {
		for c := 1; c < len(img[r])-1; c++ {
			if img[r][c] == 0 && check(getTripleVector(img, c, r)) {
				img[r][c] = 1
				count++
			}
		}
	}
	return count
}

// Проверка на удаление
func getTripleVector(img [][]int, x int, y int) []int {
	var a = make([]int, 9)
	for r := y - 1; r < y+1; r++ {
		for c := x - 1; c < x+1; c++ {
			if r < 0 || r >= len(img) || c < 0 || c >= len(img[0]) {
				a = append(a, 0)
			} else {
				a = append(a, img[r][c])
			}
		}
	}
	return a
}

// Функция сравения шаблонов с вектором
func check(vector []int) bool {

	// 4 шаблона
	if templateCompare(vector, []int{1, 1, 0, 0, 1, 0}, []int{1, 2, 3, 4, 5, 7}) {
		return true
	}
	if templateCompare(vector, []int{1, 1, 1, 0, 0, 0}, []int{0, 1, 3, 4, 5, 7}) {
		return true
	}
	if templateCompare(vector, []int{0, 1, 0, 0, 1, 1}, []int{1, 3, 4, 5, 6, 7}) {
		return true
	}
	if templateCompare(vector, []int{0, 0, 0, 1, 1, 1}, []int{1, 3, 4, 5, 7, 8}) {
		return true
	}

	// 4 шаблона
	if templateCompare(vector, []int{1, 1, 1, 0, 0, 0, 0}, []int{0, 1, 2, 3, 4, 5, 7}) {
		return true
	}
	if templateCompare(vector, []int{1, 0, 1, 0, 0, 1, 0}, []int{0, 1, 3, 4, 5, 6, 7}) {
		return true
	}
	if templateCompare(vector, []int{0, 0, 0, 0, 1, 1, 1}, []int{1, 3, 4, 5, 6, 7, 8}) {
		return true
	}
	if templateCompare(vector, []int{0, 1, 0, 0, 1, 0, 1}, []int{1, 2, 3, 4, 5, 7, 8}) {
		return true
	}

	return false
}

// Функция сравения шаблонов с вектором
func fringe(vector []int) bool {
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

	for _, target := range matrix {
		if intArrayEquals(vector, target) {
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

// Функция сравнения шаблона c вектором
func templateCompare(vector []int, template []int, allowed []int) bool {
	var count = 0

	for i, el := range vector {
		if i != allowed[count] {
			continue
		}

		if el != template[count] {
			return false
		}
		count++
	}

	return true
}
