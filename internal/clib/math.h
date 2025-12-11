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

	friend std::ostream& operator<<(std::ostream& os, Vec2 v) {
		os << "Vec2{X: " << v.X << ", Y: " << v.Y << "}";
	    return os;
	}
};

struct Vec2XOps : Vec2 {
	friend bool operator>(const Vec2XOps l, const Vec2XOps r) {
		return l.X > r.X;
	}
	friend bool operator<(const Vec2XOps l, const Vec2XOps r) {
		return l.X < r.X;
	}
	friend bool operator>=(const Vec2XOps l, const Vec2XOps r) {
		return l.X >= r.X;
	}
	friend bool operator<=(const Vec2XOps l, const Vec2XOps r) {
		return l.X <= r.X;
	}
	friend bool operator!=(const Vec2XOps l, const Vec2XOps r) {
		return l.X != r.X;
	}
	friend bool operator==(const Vec2XOps l, const Vec2XOps r) {
		return l.X == r.X;
	}
};

struct Vec2YOps : Vec2 {
	friend bool operator>(const Vec2YOps l, const Vec2YOps r) {
		return l.Y > r.Y;
	}
	friend bool operator<(const Vec2YOps l, const Vec2YOps r) {
		return l.Y < r.Y;
	}
	friend bool operator>=(const Vec2YOps l, const Vec2YOps r) {
		return l.Y >= r.Y;
	}
	friend bool operator<=(const Vec2YOps l, const Vec2YOps r) {
		return l.Y <= r.Y;
	}
	friend bool operator!=(const Vec2YOps l, const Vec2YOps r) {
		return l.Y != r.Y;
	}
	friend bool operator==(const Vec2YOps l, const Vec2YOps r) {
		return l.Y == r.Y;
	}
};

namespace Math {

// Operators -------------------------------------------------------------------
inline double Mag(Vec2 v) {
	return sqrt(v.X*v.X+v.Y*v.Y);
}

// Numerical Difference --------------------------------------------------------
// For an explanation of the formulas refer to here:
// http://code.barbellmath.net/barbell-math/providentia/wiki/Numerical-Difference-Methods

// Calculates the first derivative of data using numerical differentiation
// returning the result. h controls the delta between consecutive points. The
// accuracy of res will be proportional to h^2.
inline Vec2 FirstDerivative(FixedSlice<Vec2, 3> data, double h) {
	Vec2 res{};
	res.X=(-data[0].X+data[2].X)/(2*h);
	res.Y=(-data[0].Y+data[2].Y)/(2*h);
	return res;
}
// Calculates the first derivative of data using numerical differentiation
// returning the result. h controls the delta between consecutive points. The
// accuracy of res will be proportional to h^4.
inline Vec2 FirstDerivative(FixedSlice<Vec2, 5> data, double h) {
	Vec2 res{};
	res.X=(data[0].X -8*data[1].X +8*data[3].X -data[4].X)/(12*h);
	res.Y=(data[0].Y -8*data[1].Y +8*data[3].Y -data[4].Y)/(12*h);
	return res;
}

// Calculates the second derivative of data using numerical differentiation
// returning the result. h controls the delta between consecutive points. The
// accuracy of res will be proportional to h^2.
inline Vec2 SecondDerivative(FixedSlice<Vec2, 3> data, double h) {
	Vec2 res{};
	res.X=(data[0].X-2*data[1].X+data[2].X)/(h*h);
	res.Y=(data[0].Y-2*data[1].Y+data[2].Y)/(h*h);
	return res;
}
// Calculates the second derivative of data using numerical differentiation
// returning the result. h controls the delta between consecutive points. The
// accuracy of res will be proportional to h^4.
inline Vec2 SecondDerivative(FixedSlice<Vec2, 5> data, double h) {
	Vec2 res{};
	res.X=(-data[0].X +16*data[1].X -30*data[2].X +16*data[3].X -data[4].X)/(12*h*h);
	res.Y=(-data[0].Y +16*data[1].Y -30*data[2].Y +16*data[3].Y -data[4].Y)/(12*h*h);
	return res;
}

// Calculates the third derivative of data using numerical differentiation
// returning the result. h controls the delta between consecutive points. The
// accuracy of res will be proportional to h^2.
inline Vec2 ThirdDerivative(FixedSlice<Vec2, 5> data, double h) {
	Vec2 res{};
	res.X=(-data[0].X +2*data[1].X -2*data[3].X +data[4].X)/(2*h*h*h);
	res.Y=(-data[0].Y +2*data[1].Y -2*data[3].Y +data[4].Y)/(2*h*h*h);
	return res;
}
// Calculates the third derivative of data using numerical differentiation
// returning the result. h controls the delta between consecutive points. The
// accuracy of res will be proportional to h^4.
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

// Can be passed as N to [CalcFirstThreeDerivatives]
constexpr size_t SecondFirstOrderApprox = 5;
// Can be passed as N to [CalcFirstThreeDerivatives]
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
void CalcFirstThreeDerivatives(
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
	
	// Smear edges to the ends of the results rather than computing forward and
	// backward difference formulas. Running those calculations would provide
	// little benefit while significantly increasing complexity and maintenance
	for (size_t i=0; i<N/2 && i<data.Len(); i++) {
		first[i]=first[N/2];
		second[i]=second[N/2];
		third[i]=third[N/2];
	}
	for (size_t i=data.Len()-N/2; i<data.Len(); i++) {
		first[i]=first[data.Len()-N/2-1];
		second[i]=second[data.Len()-N/2-1];
		third[i]=third[data.Len()-N/2-1];
	}
}

// Calculates the weighted average of the supplied data using the supplied
// weights. `wTot` represents the sum of all the weights. No check is performed
// that `wTot` and `weights` match. This function is mainly used in a scenario
// where a weighted average is calculated many times with the same weights,
// allowing for improved performance if `wTot` is cached. If the sum of all the
// weights is 0 a zero-valued Vec2 is returned.
template <size_t N>
inline Vec2 WeightedAvg(
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

// Calculates the weighted average of the supplied data using the supplied
// weights. If the sum of all the weights is 0 a zero-valued Vec2 is returned.
template <size_t N>
inline Vec2 WeightedAvg(
	FixedSlice<Vec2, N> data,
	FixedSlice<double, N> weights
) {
	double wTot=0;
	for (size_t i=0; i<N; i++) {
		wTot+=weights[i];
	}
	return WeightedAvg(data, weights, wTot);
}

// Calculates a centered rolling weighted average. The weights are slid across
// the data in a window like fashion and the center data point is updated to be
// the calculated average.
template <size_t N>
void CenteredRollingWeightedAvg(
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
		tmps.Put(Math::WeightedAvg(FixedSlice<Vec2, N>(data,i), weights, wTot));
	}
	size_t offset = data.Len()-weights.Len();
	for (size_t i=0; i<tmps.Len(); i++) {
		data[offset+i]=tmps[i];
	}
}

// Finds the N smallest minimums in the supplied data, returning the number of
// minimums that were found (capped at `mins.Len()`). The `mins` slice will be
// populated with the indexes of the minimums in `data`. This is intended to be
// used with real-world data and should not be used to solve for the minimum of
// an equation.
//
// A minimum is defined to be any point that has `radius` num neighbors on both
// sides that are all decreasing up to the central point. `radius` can be set
// higher to filter out irrelevant minimums in noisy data.
//
// The type T must have the standard comparison operators defined.
//
// `maxVal` must be the largest value possible for the type T. If it is not and
// there are minimums larger than `maxVal` then those minimums will not be
// found. This could act as a filter to find the N smallest minimums less than
// `maxVal` if desired.
//
// If there are less than `mins.Len()` minimums found then the first N values
// in mins will be populated where N=the return value. The indexes in `mins` are
// are not sorted by there associated minimum values in any way.
template <typename T>
size_t NSmallestMinimums(
	Slice<T> data,
	Slice<size_t> mins,
	const T& maxVal,
	const size_t radius=1
) {
	size_t numMins=0;
	Slice<T> tmpVals(mins.Len());
	tmpVals.Fill(maxVal);
	AssociatedSlices<T, size_t> heap(tmpVals, mins);

	for (size_t i=radius; i<data.Len()-radius; ) {
		for (size_t j=i-radius+1; j<=i; j++) {
			if (data[j]>=data[j-1]) {
				i+=j-(i-radius);
				goto outerLoopEnd;
			}
		}
		for (size_t j=i+1; j<i+radius+1; j++) {
			if (data[j]<=data[j-1]) {
				i=j;
				goto outerLoopEnd;
			}
		}

		if (data[i]<heap[0].First) {
			heap[0].First=data[i];
			heap[0].Second=i;
			Heap::Max<
				AssociatedSlices<T, size_t>,
				typename AssociatedSlices<T, size_t>::Elems
			>(heap);
			numMins++;
			i+=radius+1;
		}

	outerLoopEnd:
	}

	tmpVals.Free();
	if (numMins<mins.Len()) { mins.Reverse(); }
	return std::min(numMins, mins.Len());
}

// TODO - finish
template <typename T>
size_t NLargestMaximums(
	Slice<T> data,
	Slice<size_t> maxes,
	const T& minVal,
	const size_t radius=1
) {
	size_t numMaxes=0;
	Slice<T> tmpVals(maxes.Len());
	tmpVals.Fill(minVal);
	AssociatedSlices<T, size_t> heap(tmpVals, maxes);

	for (size_t i=radius; i<data.Len()-radius; ) {
		for (size_t j=i-radius+1; j<=i; j++) {
			if (data[j]<=data[j-1]) {
				i+=j-(i-radius);
				goto outerLoopEnd;
			}
		}
		for (size_t j=i+1; j<i+radius+1; j++) {
			if (data[j]>=data[j-1]) {
				i=j;
				goto outerLoopEnd;
			}
		}

		if (data[i]>heap[0].First) {
			heap[0].First=data[i];
			heap[0].Second=i;
			Heap::Min<
				AssociatedSlices<T, size_t>,
				typename AssociatedSlices<T, size_t>::Elems
			>(heap);
			numMaxes++;
			i+=radius+1;
		}

	outerLoopEnd:
	}

	tmpVals.Free();
	if (numMaxes<maxes.Len()) { maxes.Reverse(); }
	return std::min(numMaxes, maxes.Len());
}

// TODO
// size_t LeftRoot(Slice<T> data, start idx)
// size_t RightRoot(Slice<T> data, start idx)
// [2]size_t Roots(Slice<T> data, start idx)

};

#endif
