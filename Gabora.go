package main

import (
	"fmt"
	"math"
)

// Функция серизации
func grayscale(img [][]Pixel) [][]int {
	var bImg [][]int

	for _, row := range img {
		var p []int
		for _, col := range row {
			if col.A == 0 {
				p = append(p, 0)
				continue
			}

			pixel := int(float64(col.R)*BINARY_RATION_R + float64(col.G)*BINARY_RATION_G + float64(col.B)*BINARY_RATION_B)

			if pixel > 255 {
				pixel = 255
			}

			p = append(p, pixel)
		}

		bImg = append(bImg, p)
	}

	return bImg
}

// Фильтр Собеля
func sobelay(img [][]int) [][]float64 {
	var angleMatrix [][]float64

	Gx := func(z []int) float64 {
		return float64((z[6] + 2*z[7] + z[8]) - (z[0] + 2*z[1] + z[2]))
	}
	Gy := func(z []int) float64 {
		return float64((z[2] + 2*z[5] + z[8]) - (z[0] + 2*z[3] + z[6]))
	}

	for i, _ := range img {
		var angleVecor []float64
		for j, _ := range img[0] {

			template := getTripleVector(img, j, i)

			angleVecor = append(angleVecor, math.Atan(Gy(template)/Gx(template)))
		}

		angleMatrix = append(angleMatrix, angleVecor)
	}
	return angleMatrix
}

// Константы для фильтра Габора
//const lambda = float64(1)
//const theta = float64(2)
//const psi = float64(3)
//const sigma = float64(4)
//const gamma = float64(1)

// Функция Габора
func gabor(sigma float64, theta float64, lambda int, psi int, gamma float64) {
	sigma_x := sigma
	sigma_y := float64(sigma) / gamma

	nstds := 3

	xmax := float64(math.Ceil(math.Max(1, math.Max(math.Abs(float64(nstds)*sigma_x*math.Cos(theta)), math.Abs(float64(nstds)*sigma_y*math.Sin(theta))))))
	ymax := float64(math.Ceil(math.Max(1, math.Max(math.Abs(float64(nstds)*sigma_x*math.Sin(theta)), math.Abs(float64(nstds)*sigma_y*math.Cos(theta))))))
	xmin := -xmax
	ymin := -ymax

	var gridY []float64 //ymin : ymax + 1
	for i := ymin; i <= ymax+1; i++ {
		gridY = append(gridY, i)
	}
	fmt.Println("gridY")
	fmt.Println(gridY)

	var gridX []float64 //xmin : xmax + 1
	for i := xmin; i <= xmax+1; i++ {
		gridX = append(gridX, i)
	}

	x, y := meshgrid(gridX, gridY)

	fmt.Println(x)
	fmt.Println(y)

	x_theta := matrixOnMatrixPlus(matrixOnFigureMultiple(x, math.Cos(theta)), matrixOnFigureMultiple(y, math.Sin(theta)))
	y_theta := matrixOnMatrixPlus(matrixOnFigureMultiple(x, -math.Sin(theta)), matrixOnFigureMultiple(y, math.Cos(theta)))

	fmt.Println(x_theta)
	fmt.Println(y_theta)
}

func matrixOnFigureMultiple(matrix [][]float64, num float64) [][]float64 {
	answer := make([][]float64, len(matrix))
	for i := 0; i < len(answer); i++ {
		answer[i] = make([]float64, len(matrix[0]))
	}

	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[0]); j++ {
			answer[i][j] = matrix[i][j] * num
		}
	}
	return answer
}

func matrixOnMatrixPlus(matrix [][]float64, matrix2 [][]float64) [][]float64 {
	answer := make([][]float64, len(matrix))
	for i := 0; i < len(answer); i++ {
		answer[i] = make([]float64, len(matrix[0]))
	}

	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[0]); j++ {
			answer[i][j] = matrix[i][j] + matrix2[i][j]
		}
	}

	return answer
}

//https://translated.turbopages.org/proxy_u/en-ru.ru.30a9c78e-638b4225-8aef64e1-74722d776562/https/stackoverflow.com/questions/67049518/gabor-filter-parametrs-for-fingerprint-image-enhancement
func meshgrid(x []float64, y []float64) ([][]float64, [][]float64) {

	// Создаём динамический массив
	var arrX = make([][]float64, len(y))
	for i, _ := range arrX {
		arrX[i] = make([]float64, len(x))
	}
	var arrY = make([][]float64, len(y))
	for i, _ := range arrY {
		arrY[i] = make([]float64, len(x))
	}

	for j, _ := range arrX[0] {
		for i, _ := range arrX {
			arrX[i][j] = x[j]
		}
	}
	for i, _ := range arrY {
		for j, _ := range arrY[0] {
			arrY[i][j] = y[i]
		}
	}

	return arrX, arrY
}

/*
def gabor(sigma, theta, Lambda, psi, gamma):
    sigma_x = sigma
    sigma_y = float(sigma) / gamma
	# Bounding box
    nstds = 3  # Number of standard deviation sigma

    xmax = max(abs(nstds * sigma_x * np.cos(theta)), abs(nstds * sigma_y * np.sin(theta)))
    xmax = np.ceil(max(1, xmax))

    ymax = max(abs(nstds * sigma_x * np.sin(theta)), abs(nstds * sigma_y * np.cos(theta)))
    ymax = np.ceil(max(1, ymax))
    xmin = -xmax
    ymin = -ymax

    (y, x) = np.meshgrid(np.arange(ymin, ymax + 1), np.arange(xmin, xmax + 1))

    # Rotation
    x_theta = x * np.cos(theta) + y * np.sin(theta)
    y_theta = -x * np.sin(theta) + y * np.cos(theta)

    gb = np.exp(-.5 * (x_theta ** 2 / sigma_x ** 2 + y_theta ** 2 / sigma_y ** 2)) * np.cos(2 * np.pi / Lambda * x_theta + psi)
    return gb
*/
