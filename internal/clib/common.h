#ifndef CGO_GLUE_COMMON
#define CGO_GLUE_COMMON

#include <stdio.h>
#include <stdlib.h>

// TODO - DELETE!!!
struct TimestampedVal {
	int Idx;
	double Time;
	double Value;

	static bool sortByTime(const TimestampedVal& a, const TimestampedVal& b) {
		return a.Time<b.Time;
	}

	static bool sortByValue(const TimestampedVal& a, const TimestampedVal& b) {
		return a.Value<b.Value;
	}
};

// TODO - DELETE!!!
struct PointInTime {
	double Time;
	double Value;
};

#endif
