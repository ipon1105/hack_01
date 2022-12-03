package main

import (
	"fmt"
	"math"
)

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

	xmax := int(math.Ceil(math.Max(1, math.Max(math.Abs(float64(nstds)*sigma_x*math.Cos(theta)), math.Abs(float64(nstds)*sigma_y*math.Sin(theta))))))
	ymax := int(math.Ceil(math.Max(1, math.Max(math.Abs(float64(nstds)*sigma_x*math.Sin(theta)), math.Abs(float64(nstds)*sigma_y*math.Cos(theta))))))
	xmin := -xmax
	ymin := -ymax

	var gridY []int //ymin : ymax + 1
	for i := ymin; i <= ymax+1; i++ {
		gridY = append(gridY, i)
	}
	fmt.Println("gridY")
	fmt.Println(gridY)

	var gridX []int //xmin : xmax + 1
	for i := xmin; i <= xmax+1; i++ {
		gridX = append(gridX, i)
	}
	fmt.Println("gridX")
	fmt.Println(gridX)

	//a, b := meshgrid(gridX, gridY)

	//fmt.Println(a + b)
	//
	//x_theta := a*math.Cos(theta) + b*math.Sin(theta)
	//y_theta := -a*math.Sin(theta) + b*math.Cos(theta)
}

func meshgrid(x []int, y []int) ([][]int, [][]int) {

	// Создаём динамический массив
	var arrX = make([][]int, len(y))
	for i, _ := range arrX {
		arrX[i] = make([]int, len(x))
	}
	var arrY = make([][]int, len(y))
	for i, _ := range arrY {
		arrY[i] = make([]int, len(x))
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
