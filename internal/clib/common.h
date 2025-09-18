#ifndef CGO_GLUE_COMMON
#define CGO_GLUE_COMMON

#include <math.h>

struct TimestampedVal {
	int idx;
	double_t time;
	double_t value;

	static bool sortByTime(const TimestampedVal& a, const TimestampedVal& b) {
		return a.time<b.time;
	}

	static bool sortByValue(const TimestampedVal& a, const TimestampedVal& b) {
		return a.value<b.value;
	}
};

struct Vec2 {
	static float mag(double_t x, double_t y) {
		return sqrtf(x*x+y*y);
	}
};


#endif
