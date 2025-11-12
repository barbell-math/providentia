#ifndef CGO_GLUE_MATH
#define CGO_GLUE_MATH

#include <cmath>
#include <cstdlib>
#include <iostream>
#include "dataStructs.h"
#include "glue.h"

struct Vec2 {
	double X;
	double Y;
};

std::ostream& operator<<(std::ostream& os, Vec2 v) {
	os << "Vec2{X: " << v.X << ", Y: " << v.Y << "}";
    return os;
}

namespace Math {

// Operators -------------------------------------------------------------------
inline double Mag(Vec2 v) {
	return sqrt(v.X*v.X+v.Y*v.Y);
}

// Numerical Difference --------------------------------------------------------
// For an explanation of the formulas refer to here:
// http://code.barbellmath.net/barbell-math/providentia/wiki/Numerical-Difference-Methods

// Calculates the first derivative of data using numerical differentiation
// and places the result in res. h controls the delta between consecutive
// points. The accuracy of res will be proportional to h^2.
inline Vec2 FirstDerivative(FixedSlice<Vec2, 3> data, double h) {
	Vec2 res{};
	res.X=(-data[0].X+data[2].X)/(2*h);
	res.Y=(-data[0].Y+data[2].Y)/(2*h);
	return res;
}
// Calculates the first derivative of data using numerical differentiation
// and places the result in res. h controls the delta between consecutive
// points. The accuracy of res will be proportional to h^4.
inline Vec2 FirstDerivative(FixedSlice<Vec2, 5> data, double h) {
	Vec2 res{};
	res.X=(data[0].X -8*data[1].X +8*data[3].X -data[4].X)/(12*h);
	res.Y=(data[0].Y -8*data[1].Y +8*data[3].Y -data[4].Y)/(12*h);
	return res;
}

// Calculates the second derivative of data using numerical differentiation
// and places the result in res. h controls the delta between consecutive
// points. The accuracy of res will be proportional to h^2.
inline Vec2 SecondDerivative(FixedSlice<Vec2, 3> data, double h) {
	Vec2 res{};
	res.X=(data[0].X-2*data[1].X+data[2].X)/(h*h);
	res.Y=(data[0].Y-2*data[1].Y+data[2].Y)/(h*h);
	return res;
}
// Calculates the second derivative of data using numerical differentiation
// and places the result in res. h controls the delta between consecutive
// points. The accuracy of res will be proportional to h^4.
inline Vec2 SecondDerivative(FixedSlice<Vec2, 5> data, double h) {
	Vec2 res{};
	res.X=(-data[0].X +16*data[1].X -30*data[2].X +16*data[3].X -data[4].X)/(12*h*h);
	res.Y=(-data[0].Y +16*data[1].Y -30*data[2].Y +16*data[3].Y -data[4].Y)/(12*h*h);
	return res;
}

// Calculates the third derivative of data using numerical differentiation
// and places the result in res. h controls the delta between consecutive
// points. The accuracy of res will be proportional to h^2.
inline Vec2 ThirdDerivative(FixedSlice<Vec2, 5> data, double h) {
	Vec2 res{};
	res.X=(-data[0].X +2*data[1].X -2*data[3].X +data[4].X)/(2*h*h*h);
	res.Y=(-data[0].Y +2*data[1].Y -2*data[3].Y +data[4].Y)/(2*h*h*h);
	return res;
}
// Calculates the third derivative of data using numerical differentiation
// and places the result in res. h controls the delta between consecutive
// points. The accuracy of res will be proportional to h^4.
inline Vec2 ThirdDerivative(FixedSlice<Vec2, 7> data, double h) {
	Vec2 res{};
	res.X=(
		data[0].X -8*data[1].X +13*data[2].X -13*data[4].X +8*data[5].X -data[6].X
	)/(8*h*h*h);
	res.Y=(
		data[0].Y -8*data[1].Y +13*data[2].Y -13*data[4].Y +8*data[5].Y -data[6].Y
	)/(8*h*h*h);
	return res;
}

constexpr size_t SecondFirstOrderApprox = 5;
constexpr size_t FourthOrderApprox = 7;

// Calculates the first three derivatives of data, placing the results in first,
// second, and third respectively. h controls the delta between consecutive
// points. The accuracy will be determined by N.
//   - If N=5 the accuracy of all the derivatives will be proportional to h^2
//   - If N=7 the accuracy of all the derivatives will be proportional to h^4
//
// All three derivatives will be calculated using data to avoid accumulating
// error.
template <size_t N>
inline void CalcFirstThreeDerivatives(
	Slice<Vec2> data,
	Slice<Vec2> first,
	Slice<Vec2> second,
	Slice<Vec2> third,
	double h
) {
	for (size_t i=0; i<data.Len()-N+1; i++) {
		int middleIdx=i+N/2;
		first[middleIdx] = Math::FirstDerivative(FixedSlice<Vec2, N-2>(data,i+1), h);
		second[middleIdx] = Math::SecondDerivative(FixedSlice<Vec2, N-2>(data,i+1), h);
		third[middleIdx] = Math::ThirdDerivative(FixedSlice<Vec2, N>(data,i), h);
	}
}

template <size_t N>
inline Vec2 WeightedAverage(
	FixedSlice<Vec2, N> data,
	FixedSlice<double, N> weights
) {
	double wTot=0;
	for (size_t i=0; i<N; i++) {
		wTot+=weights[i];
	}
	if (wTot==0) {
		return Vec2{};
	}
	return WeightedAverage(data, weights, wTot);
}

template <size_t N>
inline Vec2 WeightedAverage(
	FixedSlice<Vec2, N> data,
	FixedSlice<double, N> weights,
	double wTot
) {
	Vec2 rv{};
	if (wTot==0) {
		return rv;
	}
	for (size_t i=0; i<N; i++) {
		rv.X+=(data[i].X*weights[i]);
		rv.Y+=(data[i].Y*weights[i]);
	}
	rv.X/=wTot;
	rv.Y/=wTot;
	return rv;
}

template <size_t N>
void RollingWeightedAverage(
	Slice<Vec2> data,
	FixedSlice<double, N> weights,
	FixedRing<Vec2, N/2+1> tmps
) {
	double wTot=0;
	for (size_t i=0; i<weights.Len(); i++) {
		wTot+=weights[i];
	}
	if (wTot==0) {
		return;
	}

	for (size_t i=0; i<data.Len()-weights.Len()+1; i++) {
		if (i>=N/2+1) {
			data[i-1]=tmps[0];
		}
		tmps.Put(Math::WeightedAverage(
			FixedSlice<Vec2, N>(data,i), weights, wTot
		));
	}
	size_t offset = data.Len()-weights.Len();
	for (size_t i=0; i<tmps.Len(); i++) {
		data[offset+i]=tmps[i];
	}
}

};

#endif
