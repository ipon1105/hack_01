package main

import (
	"math"
)

//y' или yTheta
func y1(x float64, y float64, theta float64) float64 {
	return -x*math.Sin(theta) + y*math.Cos(theta)
}

//x' или xTheta
func x1(x float64, y float64, theta float64) float64 {
	return x*math.Cos(theta) + y*math.Sin(theta)
}

func filter(x float64, y float64, lambda float64, theta float64, psi float64, sigma float64, gamma float64) float64 {
	left := math.Exp(-0.5 * (math.Pow(x1(x, y, theta), 2)/math.Pow(sigma, 2) + math.Pow(y1(x, y, theta), 2)/math.Pow(sigma/gamma, 2)))
	right := math.Cos(2*math.Pi/lambda*x1(x, y, theta) + psi)
	return left * right
}

func filterFull(x float64, y float64, lambda float64, theta float64, psi float64, sigma float64, gamma float64) float64 {
	left := math.Exp((-math.Pow(x1(x, y, theta), 2) + math.Pow(gamma, 2)*math.Pow(y1(x, y, theta), 2)) / (2 * math.Pow(sigma, 2)))
	right := math.Cos(2*math.Pi*(x1(x, y, theta)/lambda) + psi)
	return left * right
}

///////////

func name(branches []Coord, ends []Coord, size_x int, size_y int) [][]int {
	var lines []int
	for i := 0; i < size_y; i++ {
		for j := 0; j < size_x; j++ {
			sum_c := 0
			sum_d := 0
			var z int = 0
			for _, el := range branches {
				if (size_y-i == el.Y) && (size_x-j == el.X) {
					z = 1
				}
			}
			for _, el := range ends {
				if (size_y-i == el.Y) && (size_x-j == el.X) {
					z = 1
				}
			}

			/*
				if(z==0)
					for c = (1:size(cores,1))
						sum_c = sum_c + atan(((size_y-a)-cores(c,2))/(b-cores(c,1)));
					end
					for c = (1:size(deltas,1))
						sum_d = sum_d + atan(((size_y-a)-deltas(c,2))/(b-deltas(c,1)));
					end
					line(b) = (sum_d - sum_c)/2;
				else
			*/
			if z == 0 {
				for _, el := range branches {
					sum_c += int(math.Atan(float64(((size_y - i) - el.Y) / (j - el.X))))
				}
				for _, el := range ends {
					sum_d += int(math.Atan(float64(((size_y - i) - el.Y) / (j - el.X))))
				}
				lines = append(lines, (sum_d-sum_c)/2)

			}
		}
	}

	/*
		 function [ img ] = Sh_M_orientation( cores, deltas, size_x, size_y )
		for a=(1:size_y)
		for b = (1:size_x)
		sum_c = 0;
		sum_d = 0;
		z = 0;
			for c = (1:size(cores,1))
				if(((size_y-a)==cores(c,2))&(b==cores(c,1)))
					z=1;
				end
			end

		for c = (1:size(deltas,1))
			if(((size_y-a)==deltas(c,2))&(b==deltas(c,1)))
				z=1;
			end
		end

		//-

		line(b)=0;

		end

		end

		if (a==1)

		img = line;

		else

		img = cat(1, img, line);

		end

		end

		end
	*/

	return nil
}
