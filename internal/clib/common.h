#ifndef CGO_GLUE_COMMON
#define CGO_GLUE_COMMON

#include <math.h>
#include <stdio.h>
#include <stdlib.h>

namespace Funcs {
inline void panic(const char* message) {
	printf("PANIC: %s\n", message);
	abort();
}
}

struct TimestampedVal {
	int Idx;
	double_t Time;
	double_t Value;

	static bool sortByTime(const TimestampedVal& a, const TimestampedVal& b) {
		return a.Time<b.Time;
	}

	static bool sortByValue(const TimestampedVal& a, const TimestampedVal& b) {
		return a.Value<b.Value;
	}
};

struct PointInTime {
	double_t Time;
	double_t Value;
};

#endif
