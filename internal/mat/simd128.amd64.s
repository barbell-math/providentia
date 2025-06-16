//go:build amd64 && prov_simd128

#include "textflag.h"

// func MatMul7x7By7x1(
//		l [7][7]float64,
//		r [7][1]float64,
//		res *[7][1]float64,
//	)
//
// memory layout of the stack relative to FP
//  +0   				key					argument
//  +1  through +16 	flags				argument
//  +17 through +33		slotKeys			argument
//  +34 through +39		-					alignment padding
//  +40 through +41		potentialValues 	return value
//  +42 through +43		isEmpty 			return value
TEXT ·Mul7x7By7x1(SB),NOSPLIT,$0
	RET
