package mat

import (
	sbtest "github.com/barbell-math/smoothbrain-test"
	"testing"
)

func TestAdd7x7By7x7(t *testing.T) {
	r := [7][7]float64{
		[7]float64{1, 1, 1, 1, 1, 1, 1},
		[7]float64{2, 2, 2, 2, 2, 2, 2},
		[7]float64{3, 3, 3, 3, 3, 3, 3},
		[7]float64{4, 4, 4, 4, 4, 4, 4},
		[7]float64{5, 5, 5, 5, 5, 5, 5},
		[7]float64{6, 6, 6, 6, 6, 6, 6},
		[7]float64{7, 7, 7, 7, 7, 7, 7},
	}
	l := [7][7]float64{
		[7]float64{1, 2, 3, 4, 5, 6, 7},
		[7]float64{1, 2, 3, 4, 5, 6, 7},
		[7]float64{1, 2, 3, 4, 5, 6, 7},
		[7]float64{1, 2, 3, 4, 5, 6, 7},
		[7]float64{1, 2, 3, 4, 5, 6, 7},
		[7]float64{1, 2, 3, 4, 5, 6, 7},
		[7]float64{1, 2, 3, 4, 5, 6, 7},
	}
	exp := [7][7]float64{
		[7]float64{2, 3, 4, 5, 6, 7, 8},
		[7]float64{3, 4, 5, 6, 7, 8, 9},
		[7]float64{4, 5, 6, 7, 8, 9, 10},
		[7]float64{5, 6, 7, 8, 9, 10, 11},
		[7]float64{6, 7, 8, 9, 10, 11, 12},
		[7]float64{7, 8, 9, 10, 11, 12, 13},
		[7]float64{8, 9, 10, 11, 12, 13, 14},
	}
	var res [7][7]float64
	Add7x7By7x7(l, r, &res)
	for i := 0; i < 7; i++ {
		for j := 0; j < 7; j++ {
			sbtest.EqFloat(t, exp[i][j], res[i][j], 0)
		}
	}
}

func TestAdd7x7By7x7ResIsL(t *testing.T) {
	r := [7][7]float64{
		[7]float64{1, 1, 1, 1, 1, 1, 1},
		[7]float64{2, 2, 2, 2, 2, 2, 2},
		[7]float64{3, 3, 3, 3, 3, 3, 3},
		[7]float64{4, 4, 4, 4, 4, 4, 4},
		[7]float64{5, 5, 5, 5, 5, 5, 5},
		[7]float64{6, 6, 6, 6, 6, 6, 6},
		[7]float64{7, 7, 7, 7, 7, 7, 7},
	}
	l := [7][7]float64{
		[7]float64{1, 2, 3, 4, 5, 6, 7},
		[7]float64{1, 2, 3, 4, 5, 6, 7},
		[7]float64{1, 2, 3, 4, 5, 6, 7},
		[7]float64{1, 2, 3, 4, 5, 6, 7},
		[7]float64{1, 2, 3, 4, 5, 6, 7},
		[7]float64{1, 2, 3, 4, 5, 6, 7},
		[7]float64{1, 2, 3, 4, 5, 6, 7},
	}
	exp := [7][7]float64{
		[7]float64{2, 3, 4, 5, 6, 7, 8},
		[7]float64{3, 4, 5, 6, 7, 8, 9},
		[7]float64{4, 5, 6, 7, 8, 9, 10},
		[7]float64{5, 6, 7, 8, 9, 10, 11},
		[7]float64{6, 7, 8, 9, 10, 11, 12},
		[7]float64{7, 8, 9, 10, 11, 12, 13},
		[7]float64{8, 9, 10, 11, 12, 13, 14},
	}
	Add7x7By7x7(l, r, &l)
	for i := 0; i < 7; i++ {
		for j := 0; j < 7; j++ {
			sbtest.EqFloat(t, exp[i][j], l[i][j], 0)
		}
	}
}

func TestAdd7x7By7x7ResIsR(t *testing.T) {
	r := [7][7]float64{
		[7]float64{1, 1, 1, 1, 1, 1, 1},
		[7]float64{2, 2, 2, 2, 2, 2, 2},
		[7]float64{3, 3, 3, 3, 3, 3, 3},
		[7]float64{4, 4, 4, 4, 4, 4, 4},
		[7]float64{5, 5, 5, 5, 5, 5, 5},
		[7]float64{6, 6, 6, 6, 6, 6, 6},
		[7]float64{7, 7, 7, 7, 7, 7, 7},
	}
	l := [7][7]float64{
		[7]float64{1, 2, 3, 4, 5, 6, 7},
		[7]float64{1, 2, 3, 4, 5, 6, 7},
		[7]float64{1, 2, 3, 4, 5, 6, 7},
		[7]float64{1, 2, 3, 4, 5, 6, 7},
		[7]float64{1, 2, 3, 4, 5, 6, 7},
		[7]float64{1, 2, 3, 4, 5, 6, 7},
		[7]float64{1, 2, 3, 4, 5, 6, 7},
	}
	exp := [7][7]float64{
		[7]float64{2, 3, 4, 5, 6, 7, 8},
		[7]float64{3, 4, 5, 6, 7, 8, 9},
		[7]float64{4, 5, 6, 7, 8, 9, 10},
		[7]float64{5, 6, 7, 8, 9, 10, 11},
		[7]float64{6, 7, 8, 9, 10, 11, 12},
		[7]float64{7, 8, 9, 10, 11, 12, 13},
		[7]float64{8, 9, 10, 11, 12, 13, 14},
	}
	Add7x7By7x7(l, r, &r)
	for i := 0; i < 7; i++ {
		for j := 0; j < 7; j++ {
			sbtest.EqFloat(t, exp[i][j], r[i][j], 0)
		}
	}
}

func TestMul7x7By7x1(t *testing.T) {
	r := [7][1]float64{
		[1]float64{1},
		[1]float64{2},
		[1]float64{3},
		[1]float64{4},
		[1]float64{5},
		[1]float64{6},
		[1]float64{7},
	}
	l := [7][7]float64{
		[7]float64{1, 2, 3, 4, 5, 6, 7},
		[7]float64{1, 2, 3, 4, 5, 6, 7},
		[7]float64{1, 2, 3, 4, 5, 6, 7},
		[7]float64{1, 2, 3, 4, 5, 6, 7},
		[7]float64{1, 2, 3, 4, 5, 6, 7},
		[7]float64{1, 2, 3, 4, 5, 6, 7},
		[7]float64{1, 2, 3, 4, 5, 6, 7},
	}
	var res [7][1]float64
	Mul7x7By7x1(l, r, &res)
	for i := 0; i < 7; i++ {
		sbtest.EqFloat(t, 140, res[i][0], 0)
	}

	r = [7][1]float64{
		[1]float64{1},
		[1]float64{2},
		[1]float64{3},
		[1]float64{4},
		[1]float64{5},
		[1]float64{6},
		[1]float64{7},
	}
	l = [7][7]float64{
		[7]float64{1, 1, 1, 1, 1, 1, 1},
		[7]float64{2, 2, 2, 2, 2, 2, 2},
		[7]float64{3, 3, 3, 3, 3, 3, 3},
		[7]float64{4, 4, 4, 4, 4, 4, 4},
		[7]float64{5, 5, 5, 5, 5, 5, 5},
		[7]float64{6, 6, 6, 6, 6, 6, 6},
		[7]float64{7, 7, 7, 7, 7, 7, 7},
	}
	exp := [7][1]float64{
		[1]float64{28},
		[1]float64{56},
		[1]float64{84},
		[1]float64{112},
		[1]float64{140},
		[1]float64{168},
		[1]float64{196},
	}
	Mul7x7By7x1(l, r, &res)
	for i := 0; i < 7; i++ {
		sbtest.EqFloat(t, exp[i][0], res[i][0], 0)
	}
}

func TestMul7x7By7x1ResIsR(t *testing.T) {
	r := [7][1]float64{
		[1]float64{1},
		[1]float64{2},
		[1]float64{3},
		[1]float64{4},
		[1]float64{5},
		[1]float64{6},
		[1]float64{7},
	}
	l := [7][7]float64{
		[7]float64{1, 1, 1, 1, 1, 1, 1},
		[7]float64{2, 2, 2, 2, 2, 2, 2},
		[7]float64{3, 3, 3, 3, 3, 3, 3},
		[7]float64{4, 4, 4, 4, 4, 4, 4},
		[7]float64{5, 5, 5, 5, 5, 5, 5},
		[7]float64{6, 6, 6, 6, 6, 6, 6},
		[7]float64{7, 7, 7, 7, 7, 7, 7},
	}
	exp := [7][1]float64{
		[1]float64{28},
		[1]float64{56},
		[1]float64{84},
		[1]float64{112},
		[1]float64{140},
		[1]float64{168},
		[1]float64{196},
	}
	Mul7x7By7x1(l, r, &r)
	for i := 0; i < 7; i++ {
		sbtest.EqFloat(t, exp[i][0], r[i][0], 0)
	}
}

func TestTermMul7x1By7x1(t *testing.T) {
	r := [7][1]float64{
		[1]float64{1},
		[1]float64{2},
		[1]float64{3},
		[1]float64{4},
		[1]float64{5},
		[1]float64{6},
		[1]float64{7},
	}
	l := [7][1]float64{
		[1]float64{7},
		[1]float64{6},
		[1]float64{5},
		[1]float64{4},
		[1]float64{3},
		[1]float64{2},
		[1]float64{1},
	}
	exp := [7][1]float64{
		[1]float64{7},
		[1]float64{12},
		[1]float64{15},
		[1]float64{16},
		[1]float64{15},
		[1]float64{12},
		[1]float64{7},
	}
	var res [7][1]float64
	TermMul7x1By7x1(l, r, &res)
	for i := 0; i < 7; i++ {
		sbtest.EqFloat(t, exp[i][0], res[i][0], 0)
	}
}

func TestTermMul7x1By7x1ResIsL(t *testing.T) {
	r := [7][1]float64{
		[1]float64{1},
		[1]float64{2},
		[1]float64{3},
		[1]float64{4},
		[1]float64{5},
		[1]float64{6},
		[1]float64{7},
	}
	l := [7][1]float64{
		[1]float64{7},
		[1]float64{6},
		[1]float64{5},
		[1]float64{4},
		[1]float64{3},
		[1]float64{2},
		[1]float64{1},
	}
	exp := [7][1]float64{
		[1]float64{7},
		[1]float64{12},
		[1]float64{15},
		[1]float64{16},
		[1]float64{15},
		[1]float64{12},
		[1]float64{7},
	}
	TermMul7x1By7x1(l, r, &l)
	for i := 0; i < 7; i++ {
		sbtest.EqFloat(t, exp[i][0], l[i][0], 0)
	}
}

func TestTermMul7x1By7x1ResIsR(t *testing.T) {
	r := [7][1]float64{
		[1]float64{1},
		[1]float64{2},
		[1]float64{3},
		[1]float64{4},
		[1]float64{5},
		[1]float64{6},
		[1]float64{7},
	}
	l := [7][1]float64{
		[1]float64{7},
		[1]float64{6},
		[1]float64{5},
		[1]float64{4},
		[1]float64{3},
		[1]float64{2},
		[1]float64{1},
	}
	exp := [7][1]float64{
		[1]float64{7},
		[1]float64{12},
		[1]float64{15},
		[1]float64{16},
		[1]float64{15},
		[1]float64{12},
		[1]float64{7},
	}
	TermMul7x1By7x1(l, r, &r)
	for i := 0; i < 7; i++ {
		sbtest.EqFloat(t, exp[i][0], r[i][0], 0)
	}
}

func TestMul7x1ByScalar(t *testing.T) {
	l := [7][1]float64{
		[1]float64{1},
		[1]float64{2},
		[1]float64{3},
		[1]float64{4},
		[1]float64{5},
		[1]float64{6},
		[1]float64{7},
	}
	r := float64(2)
	exp := [7][1]float64{
		[1]float64{2},
		[1]float64{4},
		[1]float64{6},
		[1]float64{8},
		[1]float64{10},
		[1]float64{12},
		[1]float64{14},
	}
	var res [7][1]float64
	Mul7x1ByScalar(l, r, &res)
	for i := 0; i < 7; i++ {
		sbtest.EqFloat(t, exp[i][0], res[i][0], 0)
	}
}

func TestMul7x1ByScalarResIsL(t *testing.T) {
	l := [7][1]float64{
		[1]float64{1},
		[1]float64{2},
		[1]float64{3},
		[1]float64{4},
		[1]float64{5},
		[1]float64{6},
		[1]float64{7},
	}
	r := float64(2)
	exp := [7][1]float64{
		[1]float64{2},
		[1]float64{4},
		[1]float64{6},
		[1]float64{8},
		[1]float64{10},
		[1]float64{12},
		[1]float64{14},
	}
	Mul7x1ByScalar(l, r, &l)
	for i := 0; i < 7; i++ {
		sbtest.EqFloat(t, exp[i][0], l[i][0], 0)
	}
}
