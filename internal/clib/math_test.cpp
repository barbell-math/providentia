#include <stdbool.h>
#include "./tests.gen.h"
#include "./asserts.gen.h"
#include "./math.h"

extern "C" bool TestFirstDerivativeVec2SecondOrder(void) {
	Math::Vec2 data[3]={};

	for (int i=0; i<3; i++) {
		data[i].X=1;
		data[i].Y=i;
	}
	Math::Vec2 res=Math::FirstDerivative(FixedSlice<Math::Vec2,3>(data), 1);
	EQ(res.X, 0.0)
	EQ(res.Y, 1.0)
	res=Math::FirstDerivative(FixedSlice<Math::Vec2,3>(data), 0.5);
	EQ(res.X, 0.0)
	EQ(res.Y, 2.0)

	for (int i=0; i<3; i++) {
		data[i].X=i*i;
		data[i].Y=i*i*i;
	}
	res=Math::FirstDerivative(FixedSlice<Math::Vec2,3>(data), 1);
	EQ(res.X, 2.0)
	EQ(res.Y, 4.0)	// Not 3 because second order is not accurate enough

	return true;
}

extern "C" bool TestFirstDerivativeVec2FourthOrder(void) {
	Math::Vec2 data[5]={};

	for (int i=0; i<5; i++) {
		data[i].X=1;
		data[i].Y=i;
	}
	Math::Vec2 res=Math::FirstDerivative(FixedSlice<Math::Vec2,5>(data), 1);
	EQ(res.X, 0.0)
	EQ(res.Y, 1.0)
	res=Math::FirstDerivative(FixedSlice<Math::Vec2,5>(data), 0.5);
	EQ(res.X, 0.0)
	EQ(res.Y, 2.0)

	for (int i=0; i<5; i++) {
		data[i].X=i*i;
		data[i].Y=i*i*i;
	}
	res=Math::FirstDerivative(FixedSlice<Math::Vec2,5>(data), 1);
	EQ(res.X, 2.0 * 2.0) // 2x
	EQ(res.Y, 3.0 * 2.0 * 2.0) // 3x^2

	return true;
}

extern "C" bool TestSecondDerivativeVec2SecondOrder(void) {
	Math::Vec2 data[3]={};

	for (int i=0; i<3; i++) {
		data[i].X=1;
		data[i].Y=i;
	}
	Math::Vec2 res=Math::SecondDerivative(FixedSlice<Math::Vec2,3>(data), 1);
	EQ(res.X, 0.0)
	EQ(res.Y, 0.0)
	res=Math::SecondDerivative(FixedSlice<Math::Vec2,3>(data), 0.5);
	EQ(res.X, 0.0)
	EQ(res.Y, 0.0)

	for (int i=0; i<3; i++) {
		data[i].X=i*i;
		data[i].Y=i*i*i;
	}
	res=Math::SecondDerivative(FixedSlice<Math::Vec2,3>(data), 1);
	EQ(res.X, 2.0)
	EQ(res.Y, 6.0)
	res=Math::SecondDerivative(FixedSlice<Math::Vec2,3>(data), 0.5);
	EQ(res.X, 8.0)
	EQ(res.Y, 24.0)

	return true;
}

extern "C" bool TestSecondDerivativeVec2FourthOrder(void) {
	Math::Vec2 data[5]={};

	for (int i=0; i<5; i++) {
		data[i].X=1;
		data[i].Y=i;
	}
	Math::Vec2 res=Math::SecondDerivative(FixedSlice<Math::Vec2,5>(data), 1);
	EQ(res.X, 0.0)
	EQ(res.Y, 0.0)
	res=Math::SecondDerivative(FixedSlice<Math::Vec2,5>(data), 0.5);
	EQ(res.X, 0.0)
	EQ(res.Y, 0.0)

	for (int i=0; i<5; i++) {
		data[i].X=i*i;
		data[i].Y=i*i*i;
	}
	res=Math::SecondDerivative(FixedSlice<Math::Vec2,5>(data), 1);
	EQ(res.X, 2.0)
	EQ(res.Y, 12.0)
	res=Math::SecondDerivative(FixedSlice<Math::Vec2,5>(data), 0.5);
	EQ(res.X, 8.0)
	EQ(res.Y, 48.0)

	return true;
}

extern "C" bool TestThirdDerivativeVec2SecondOrder(void) {
	Math::Vec2 data[5]={};

	for (int i=0; i<5; i++) {
		data[i].X=i*i;
		data[i].Y=i*i*i;
	}
	Math::Vec2 res=Math::ThirdDerivative(FixedSlice<Math::Vec2,5>(data), 1);
	EQ(res.X, 0.0)
	EQ(res.Y, 6.0)
	res=Math::ThirdDerivative(FixedSlice<Math::Vec2,5>(data), 0.5);
	EQ(res.X, 0.0)
	EQ(res.Y, 48.0)

	return true;
}

extern "C" bool TestThirdDerivativeVec2FourthOrder(void) {
	Math::Vec2 data[7]={};

	for (int i=0; i<7; i++) {
		data[i].X=i*i;
		data[i].Y=i*i*i;
	}
	Math::Vec2 res=Math::ThirdDerivative(FixedSlice<Math::Vec2,7>(data), 1);
	EQ(res.X, 0.0)
	EQ(res.Y, 6.0)
	res=Math::ThirdDerivative(FixedSlice<Math::Vec2,7>(data), 0.5);
	EQ(res.X, 0.0)
	EQ(res.Y, 48.0)

	return true;
}

extern "C" bool TestWeightedAvgVec2(void) {
	Math::Vec2 data[5]={
		Math::Vec2{.X=0, .Y=0},
		Math::Vec2{.X=1, .Y=1},
		Math::Vec2{.X=2, .Y=2},
		Math::Vec2{.X=1, .Y=1},
		Math::Vec2{.X=0, .Y=0},
	};
	double weights[5]={0,1,2,1,0};
	Math::Vec2 avg=Math::WeightedAvg(
		FixedSlice<Math::Vec2,5>(data), FixedSlice<double, 5>(weights)
	);
	EQ(avg.X, 1.5)
	EQ(avg.Y, 1.5)

	return true;
}

extern "C" bool TestWeightedAvgVec2WeightProvided(void) {
	Math::Vec2 data[5]={
		Math::Vec2{.X=0, .Y=0},
		Math::Vec2{.X=1, .Y=1},
		Math::Vec2{.X=2, .Y=2},
		Math::Vec2{.X=1, .Y=1},
		Math::Vec2{.X=0, .Y=0},
	};
	double weights[5]={0,1,2,1,0};
	Math::Vec2 avg=Math::WeightedAvg(
		FixedSlice<Math::Vec2,5>(data), FixedSlice<double, 5>(weights), 4
	);
	EQ(avg.X, 1.5)
	EQ(avg.Y, 1.5)

	return true;
}

extern "C" bool TestCenteredRollingWeightedAvgVec2WeightsSumToZero(void) {
	Math::Vec2 data[5]={
		Math::Vec2{.X=0, .Y=0},
		Math::Vec2{.X=1, .Y=1},
		Math::Vec2{.X=2, .Y=2},
		Math::Vec2{.X=1, .Y=1},
		Math::Vec2{.X=0, .Y=0},
	};
	Math::Vec2 tmps[2]={};
	double weights[3]={1,-2,1};
	Math::CenteredRollingWeightedAvg(
		Slice<Math::Vec2>(data, 5),
		FixedSlice<double, 3>(weights),
		FixedRing<Math::Vec2,2>(tmps)
	);

	EQ(data[0].X, 0.0)
	EQ(data[0].Y, 0.0)
	EQ(data[1].X, 1.0)
	EQ(data[1].Y, 1.0)
	EQ(data[2].X, 2.0)
	EQ(data[2].Y, 2.0)
	EQ(data[3].X, 1.0)
	EQ(data[3].Y, 1.0)
	EQ(data[4].X, 0.0)
	EQ(data[4].Y, 0.0)

	return true;
}

extern "C" bool TestCenteredRollingWeightedAvgVec2(void) {
	Math::Vec2 data[7]={
		Math::Vec2{.X=0, .Y=0},
		Math::Vec2{.X=1, .Y=1},
		Math::Vec2{.X=2, .Y=2},
		Math::Vec2{.X=1, .Y=1},
		Math::Vec2{.X=0, .Y=0},
		Math::Vec2{.X=1, .Y=1},
		Math::Vec2{.X=2, .Y=2},
	};
	Math::Vec2 tmps[2]={};
	double weights[3]={1,2,1};
	Math::CenteredRollingWeightedAvg(
		Slice<Math::Vec2>(data, 7),
		FixedSlice<double, 3>(weights),
		FixedRing<Math::Vec2,2>(tmps)
	);

	EQ(data[0].X, 0.0)
	EQ(data[0].Y, 0.0)
	EQ(data[1].X, 1.0)
	EQ(data[1].Y, 1.0)
	EQ(data[2].X, 1.5)
	EQ(data[2].Y, 1.5)
	EQ(data[3].X, 1.0)
	EQ(data[3].Y, 1.0)
	EQ(data[4].X, 0.5)
	EQ(data[4].Y, 0.5)
	EQ(data[5].X, 1.0)
	EQ(data[5].Y, 1.0)
	EQ(data[6].X, 2.0)
	EQ(data[6].Y, 2.0)

	return true;
}

extern "C" bool TestNSmallestMinimumsMoreExpectedThanPresent(void) {
	Math::Vec2YOps data[11]={
		(Math::Vec2YOps)Math::Vec2{.X=0, .Y=0},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=2, .Y=2},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=0, .Y=0},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=2, .Y=2},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=0, .Y=0},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=-1},
		(Math::Vec2YOps)Math::Vec2{.X=2, .Y=-2},
	};
	size_t results[2]={};
	size_t numMins = Math::NSmallestMinimums(
		Slice<Math::Vec2YOps>(data, 11),
		Slice<size_t>(results, 2),
		(Math::Vec2YOps)Math::Vec2{
			.X=std::numeric_limits<double>::infinity(),
			.Y=std::numeric_limits<double>::infinity(),
		}
	);

	EQ(numMins, (size_t)1);
	EQ(results[0], (size_t)4);
	EQ(results[1], (size_t)0);

	return true;
}

extern "C" bool TestNSmallestMinimumsMorePresentThanExpected(void) {
	Math::Vec2YOps data[12]={
		(Math::Vec2YOps)Math::Vec2{.X=0, .Y=2},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=2, .Y=-1},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=0, .Y=2},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=2, .Y=0},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=0, .Y=2},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=2, .Y=-2},
		(Math::Vec2YOps)Math::Vec2{.X=2, .Y=1},
	};
	size_t results[2]={};
	size_t numMins = Math::NSmallestMinimums(
		Slice<Math::Vec2YOps>(data, 12),
		Slice<size_t>(results, 2),
		(Math::Vec2YOps)Math::Vec2{
			.X=std::numeric_limits<double>::infinity(),
			.Y=std::numeric_limits<double>::infinity(),
		}
	);

	EQ(numMins, (size_t)2);
	EQ(results[0], (size_t)2);
	EQ(results[1], (size_t)10);

	return true;
}

extern "C" bool TestNSmallestMinimumsLargerWindow(void) {
	Math::Vec2YOps data[12]={
		(Math::Vec2YOps)Math::Vec2{.X=0, .Y=2},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=2, .Y=-1},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=0, .Y=2},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=2, .Y=0},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=-1},
		(Math::Vec2YOps)Math::Vec2{.X=0, .Y=0},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=2, .Y=2},
		(Math::Vec2YOps)Math::Vec2{.X=2, .Y=1},
	};
	size_t results[2]={};
	size_t numMins = Math::NSmallestMinimums(
		Slice<Math::Vec2YOps>(data, 12),
		Slice<size_t>(results, 2),
		(Math::Vec2YOps)Math::Vec2{
			.X=std::numeric_limits<double>::infinity(),
			.Y=std::numeric_limits<double>::infinity(),
		},
		3
	);

	EQ(numMins, (size_t)1);
	EQ(results[0], (size_t)7);
	EQ(results[1], (size_t)0);

	return true;
}

extern "C" bool TestNLargestMaximumsMoreExpectedThanPresent(void) {
	Math::Vec2YOps data[11]={
		(Math::Vec2YOps)Math::Vec2{.X=0, .Y=2},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=2, .Y=0},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=0, .Y=2},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=2, .Y=0},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=0, .Y=2},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=3},
		(Math::Vec2YOps)Math::Vec2{.X=2, .Y=4},
	};
	size_t results[2]={};
	size_t numMaxes = Math::NLargestMaximums(
		Slice<Math::Vec2YOps>(data, 11),
		Slice<size_t>(results, 2),
		(Math::Vec2YOps)Math::Vec2{
			.X=-1*std::numeric_limits<double>::infinity(),
			.Y=-1*std::numeric_limits<double>::infinity(),
		}
	);

	EQ(numMaxes, (size_t)1);
	EQ(results[0], (size_t)4);
	EQ(results[1], (size_t)0);

	return true;
}

extern "C" bool TestNLargestMaximumssMorePresentThanExpected(void) {
	Math::Vec2YOps data[12]={
		(Math::Vec2YOps)Math::Vec2{.X=0, .Y=0},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=2, .Y=3},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=0, .Y=0},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=2, .Y=2},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=0, .Y=0},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=2, .Y=4},
		(Math::Vec2YOps)Math::Vec2{.X=2, .Y=1},
	};
	size_t results[2]={};
	size_t numMaxes = Math::NLargestMaximums(
		Slice<Math::Vec2YOps>(data, 12),
		Slice<size_t>(results, 2),
		(Math::Vec2YOps)Math::Vec2{
			.X=-1*std::numeric_limits<double>::infinity(),
			.Y=-1*std::numeric_limits<double>::infinity(),
		}
	);

	EQ(numMaxes, (size_t)2);
	EQ(results[0], (size_t)2);
	EQ(results[1], (size_t)10);

	return true;
}

extern "C" bool TestNLargestMaximumsLargerWindow(void) {
	Math::Vec2YOps data[12]={
		(Math::Vec2YOps)Math::Vec2{.X=0, .Y=0},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=2, .Y=2},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=0, .Y=0},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=2, .Y=2},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=3},
		(Math::Vec2YOps)Math::Vec2{.X=0, .Y=2},
		(Math::Vec2YOps)Math::Vec2{.X=1, .Y=1},
		(Math::Vec2YOps)Math::Vec2{.X=2, .Y=0},
		(Math::Vec2YOps)Math::Vec2{.X=2, .Y=0},
	};
	size_t results[2]={};
	size_t numMaxes = Math::NLargestMaximums(
		Slice<Math::Vec2YOps>(data, 12),
		Slice<size_t>(results, 2),
		(Math::Vec2YOps)Math::Vec2{
			.X=-1*std::numeric_limits<double>::infinity(),
			.Y=-1*std::numeric_limits<double>::infinity(),
		},
		3
	);

	EQ(numMaxes, (size_t)1);
	EQ(results[0], (size_t)7);
	EQ(results[1], (size_t)0);

	return true;
}
