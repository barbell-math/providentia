#ifndef CGO_GLUE_DATASTRUCTS
#define CGO_GLUE_DATASTRUCTS

#include <cstddef>
#include <cstdio>
#include <cstdlib>
#include <iostream>
#include "common.h"
#include "errors.h"

template <typename  T>
struct Slice {
	T *data;
	size_t len;

	Slice(size_t len) {
		this->data=calloc(len, sizeof(T));	// TODO - need better allocation strat
		if (this->data==nullptr) {
			throw Err::OOM{ .Desc="Could not calloc slice" };
		}
		this->len=len;
	}
	Slice(T *data, int len): data(data), len(len) {}

	inline size_t Len() { return this->len; }

	T& operator[](size_t idx) {
		if (idx>=this->len) {
			throw Err::IndexOutOfBounds{ .Len=this->len, .Idx=idx };
		}
		return this->data[idx];
	}

	Slice<T> operator()(size_t start, size_t end) {
		if (start<0 || start>=this->len) {
			throw Err::IndexOutOfBounds{ .Len=this->len, .Idx=start };
		}
		if (end>this->len) {
			throw Err::IndexOutOfBounds{ .Len=this->len, .Idx=start };
		}
		if (end<start) {
			throw Err::StartAfterEnd<size_t>{ .Start=start, .End=end };
		}
		return Slice<T>(&this->data[start], end-start);
	}

	void Free() {
		Free(this->data);
		this->len=0;
		this->data=nullptr;
	}
};

template<typename T, size_t N>
struct FixedSlice {
	T *data;

	FixedSlice(T *data): data(data) {}
	FixedSlice(Slice<T> s): data(s.data) {
		if (s.len<N) {
			throw Err::IndexOutOfBounds{ .Len=s.len, .Idx=N };
		}
	}
	FixedSlice(Slice<T> s, size_t start): data(s(start, start+N).data) {}

	T& operator[](size_t idx) {
		if (idx>=N) {
			throw Err::IndexOutOfBounds{ .Len=N, .Idx=idx };
		}
		return this->data[idx];
	}

	inline size_t Len() { return N; }

	Slice<T> operator()(size_t start, size_t end) {
		if (start<0 || start>=N) {
			throw Err::IndexOutOfBounds{ .Len=N, .Idx=start };
		}
		if (end>N) {
			throw Err::IndexOutOfBounds{ .Len=N, .Idx=start };
		}
		if (end<start) {
			throw Err::StartAfterEnd<size_t>{ .Start=start, .End=end };
		}
		return Slice<T>(&this->data[start], end-start);
	}

	void Free() {
		Free(this->data);
		this->data=nullptr;
	}
};

template<typename T, size_t N>
struct FixedRing {
	T *data;
	size_t curIdx = N-1;

	FixedRing(T *data): data(data) {}
	FixedRing(FixedSlice<T, N> a): data(a.data) {}
	FixedRing(Slice<T> s): data(s.data) {
		if (s.Len()<N) {
			throw Err::IndexOutOfBounds{ .Len=s.Len(), .Idx=N };
		}
	}
	FixedRing(Slice<T> s, size_t start): data(s(start, start+N).data) {}

	inline size_t Len() { return N; }

	T& operator[](size_t idx) {
		if (idx>=N) {
			throw Err::IndexOutOfBounds{ .Len=N, .Idx=idx };
		}
		return this->data[(this->curIdx+1-N+idx)%N];
	}

	void Put(T v) { this->data[(++this->curIdx)%N]=v; }

	void Free() {
		Free(this->data);
		this->curIdx=0;
		this->data=nullptr;
	}
};

template <typename T>
std::ostream& operator<<(std::ostream& os, Slice<T> s) {
	os << "Slice[";
	for (size_t i=0; i<s.Len(); i++) {
		os << s[i];
		if (i+1<s.Len()) {
			os << ", ";
		}
	}
	os << "]";
    return os;
}

template <typename T, size_t N>
std::ostream& operator<<(std::ostream& os, FixedSlice<T, N> s) {
	os << "FixedSlice(" << s.Len() << ")[";
	for (size_t i=0; i<s.Len(); i++) {
		os << s[i];
		if (i+1<s.Len()) {
			os << ", ";
		}
	}
	os << "]";
    return os;
}

template <typename T, size_t N>
std::ostream& operator<<(std::ostream& os, FixedRing<T, N> s) {
	os << "FixedRing(" << s.Len() << ")[";
	for (size_t i=0; i<s.Len(); i++) {
		os << s[i];
		if (i+1<s.Len()) {
			os << ", ";
		}
	}
	os << "]";
    return os;
}

#endif
