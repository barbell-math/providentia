//go:build !prov_simd128 && !prov_simd256 && !prov_simd512

package mat

import "math"

// TODO - should these all accept all pointer arguments??

func Add7x7By7x7(l [7][7]float64, r [7][7]float64, res *[7][7]float64) {
	for row := 0; row < 7; row++ {
		for col := 0; col < 7; col++ {
			res[row][col] = l[row][col] + r[row][col]
		}
	}
}

func Mul7x7By7x1(l [7][7]float64, r [7][1]float64, res *[7][1]float64) {
	for row := 0; row < 7; row++ {
		res[row][0] = 0
		for col := 0; col < 7; col++ {
			res[row][0] += l[row][col] * r[col][0]
		}
	}
}

func TermMul7x1By7x1(l [7][1]float64, r [7][1]float64, res *[7][1]float64) {
	for row := 0; row < 7; row++ {
		res[row][0] = l[row][0] * r[row][0]
	}
}

func Mul7x1ByScalar(l [7][1]float64, r float64, res *[7][1]float64) {
	for row := 0; row < 7; row++ {
		res[row][0] = l[row][0] * r
	}
}

func CholeskyDecomp7x7(l [7][7]float64, res *[7][7]float64) {
	for i := 0; i < 7; i++ {
		for j := 0; j <= i; j++ {
			sum := float64(0)
			for k := 0; k < j; k++ {
				sum += res[i][k] * res[j][k]
			}

			if i == j {
				res[i][j] = math.Sqrt(l[i][i] - sum)
			} else {
				res[i][j] = (1.0 / res[j][j] * (l[i][j] - sum))
			}
		}
	}
}
