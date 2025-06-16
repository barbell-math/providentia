//go:build prov_simd128 || prov_simd256 || prov_simd512

package mat

func Add7x7By7x7(l [7][7]float64, r [7][7]float64, res *[7][7]float64)
func Mul7x7By7x1(l [7][7]float64, r [7][1]float64, res *[7][1]float64)
func TermMul7x1By7x1(l [7][1]float64, r [7][1]float64, res *[7][1]float64)
func Mul7x1ByScalar(l [7][1]float64, r float64, res *[7][1]float64)
