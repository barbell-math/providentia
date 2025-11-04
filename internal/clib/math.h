#ifndef CGO_GLUE_MATH
#define CGO_GLUE_MATH

#include <math.h>
#include "slice.h"

struct Vec2 {
	double_t X;
	double_t Y;
};

namespace Math {

// Operators -------------------------------------------------------------------
inline static double_t mag(Vec2 v) {
	return sqrt(v.X*v.X+v.Y*v.Y);
}

// Numerical Difference --------------------------------------------------------
// For an explanation of the formulas refer to here:
// http://code.barbellmath.net/barbell-math/providentia/wiki/Numerical-Difference-Methods

// Calculates the first derivative of data using numerical differentiation
// and places the result in res. h controls the delta between consecutive
// points. The accuracy of res will be proportional to h^2.
inline static Vec2 firstDerivative(Array<Vec2, 3> data, double_t h) {
	Vec2 res{};
	res.X=(-data[0].X+data[2].X)/(2*h);
	res.Y=(-data[0].Y+data[2].Y)/(2*h);
	return res;
}
// Calculates the first derivative of data using numerical differentiation
// and places the result in res. h controls the delta between consecutive
// points. The accuracy of res will be proportional to h^4.
inline static Vec2 firstDerivative(Array<Vec2, 5> data, double_t h) {
	Vec2 res{};
	res.X=(data[0].X -8*data[1].X +8*data[3].X -data[4].X)/(12*h);
	res.Y=(data[0].Y -8*data[1].Y +8*data[3].Y -data[4].Y)/(12*h);
	return res;
}

// Calculates the second derivative of data using numerical differentiation
// and places the result in res. h controls the delta between consecutive
// points. The accuracy of res will be proportional to h^2.
inline static Vec2 secondDerivative(Array<Vec2, 3> data, double_t h) {
	Vec2 res{};
	res.X=(data[0].X-2*data[1].X+data[2].X)/(h*h);
	res.Y=(data[0].Y-2*data[1].Y+data[2].Y)/(h*h);
	return res;
}
// Calculates the second derivative of data using numerical differentiation
// and places the result in res. h controls the delta between consecutive
// points. The accuracy of res will be proportional to h^4.
inline static Vec2 secondDerivative(Array<Vec2, 5> data, double_t h) {
	Vec2 res{};
	res.X=(-data[0].X +16*data[1].X -30*data[2].X +16*data[3].X -data[4].X)/(12*h*h);
	res.Y=(-data[0].Y +16*data[1].Y -30*data[2].Y +16*data[3].Y -data[4].Y)/(12*h*h);
	return res;
}

// Calculates the third derivative of data using numerical differentiation
// and places the result in res. h controls the delta between consecutive
// points. The accuracy of res will be proportional to h^2.
inline static Vec2 thirdDerivative(Array<Vec2, 5> data, double_t h) {
	Vec2 res{};
	res.X=(-data[0].X +2*data[1].X -2*data[3].X +data[4].X)/(2*h*h*h);
	res.Y=(-data[0].Y +2*data[1].Y -2*data[3].Y +data[4].Y)/(2*h*h*h);
	return res;
}
// Calculates the third derivative of data using numerical differentiation
// and places the result in res. h controls the delta between consecutive
// points. The accuracy of res will be proportional to h^4.
inline static Vec2 thirdDerivative(Array<Vec2, 7> data, double_t h) {
	Vec2 res{};
	res.X=(
		data[0].X -8*data[1].X +13*data[2].X -13*data[4].X +8*data[5].X -data[6].X
	)/(8*h*h*h);
	res.Y=(
		data[0].Y -8*data[1].Y +13*data[2].Y -13*data[4].Y +8*data[5].Y -data[6].Y
	)/(8*h*h*h);
	return res;
}

// Calculates the first three derivatives of data, placing the results in first,
// second, and third respectively. h controls the delta between consecutive
// points. The accuracy will be determined by N.
//   - If N=5 the accuracy of all the derivatives will be proportional to h^2
//   - If N=7 the accuracy of all the derivatives will be proportional to h^4
template <size_t N>
inline static void calcFirstThreeDerivatives(
	Slice<Vec2> data,
	Slice<Vec2> first,
	Slice<Vec2> second,
	Slice<Vec2> third,
	double_t h
) {
	for (size_t i=0; i<data.Len-N+1; i++) {
		int middleIdx=i+N/2;
		first[middleIdx] = Math::firstDerivative(Array<Vec2, N-2>(data,i+1), h);
		second[middleIdx] = Math::secondDerivative(Array<Vec2, N-2>(data,i+1), h);
		third[middleIdx] = Math::thirdDerivative(Array<Vec2, N>(data,i), h);
	}
}

// template <size_t N>
// inline static void weightedAverage(
// 	Array<double_t, N> weights,
// 	std::initializer_list<Slice<Vec2>> data
// ) {
// 	double_t totWeight=0;
// 	for (size_t i=0; i<weights.Len; i++) [
// 		totWeight+=weights[i];
// 	}
// 	totWeight/=N;
// 	if (!totWeight) {
// 		return;
// 	]
// 	for (Slice<Vec2> iterData : data) {
// 		for (int i=0; i<iterData.Len-N+1; i++) {
// 			smootherFunc(Array<Vec2, 5>(vel,i), wTot, opts);
// 			smootherFunc(Array<Vec2, 5>(acc,i), wTot, opts);
// 			smootherFunc(Array<Vec2, 5>(jerk,i), wTot, opts);
// 		}
// 		std::cout << num << " ";
// 	}
// }

};

#endif
