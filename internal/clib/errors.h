#ifndef CGO_GLUE_ERRORS
#define CGO_GLUE_ERRORS

#include <cstddef>

namespace Err {

struct OOM {
	const char *Msg = "Out of memory";
	const char *Desc;
};

struct IndexOutOfBounds {
	const char *Msg = "Index out of bounds";
	const size_t Len;
	const size_t Idx;
};

template <typename T>
struct StartAfterEnd {
	const char *Msg = "Start value after end value";
	const T Start;
	const T End;
};

};

#endif
