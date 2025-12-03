#ifndef CGO_GLUE_DATASTRUCTS
#define CGO_GLUE_DATASTRUCTS

#include <cstddef>
#include <cstdio>
#include <cstdlib>
#include <iostream>
#include "common.h"
#include "errors.h"

template <typename T>
struct Slice {
	T *data;
	size_t len;

	Slice(size_t len) {
		this->data=(T*)calloc(len, sizeof(T));
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

	void Zero() { memset(this->data, 0, this->len); }
	void Fill(T val) {
		for (size_t i=0; i<this->len; i++) {
			this->data[i]=val;
		}
	}
	void Reverse() {
		for (size_t i=0; i<this->len/2; i++) {
			std::swap(this->data[i], this->data[this->len-i-1]);
		}
	}

	void Free() {
		free(this->data);
		this->len=0;
		this->data=nullptr;
	}

	friend std::ostream& operator<<(std::ostream& os, Slice<T> s) {
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
		free(this->data);
		this->data=nullptr;
	}

	friend std::ostream& operator<<(std::ostream& os, FixedSlice<T, N> s) {
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
		free(this->data);
		this->curIdx=0;
		this->data=nullptr;
	}

	friend std::ostream& operator<<(std::ostream& os, FixedRing<T, N> s) {
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
};

template <typename T, typename U>
struct AssociatedSlices {
	Slice<T> first;
	Slice<U> second;

	struct Elems {
		T& First;
		U& Second;

		Elems& operator=(const Elems& other) {
			this->First=other.First;
			this->Second=other.Second;
			return *this;
		}

		friend void swap(Elems l, Elems r) { 
			std::swap(l.First, r.First);
			std::swap(l.Second, r.Second);
		}

		friend bool operator>(const Elems l, const Elems r) {
			if (l.First>r.First) { return true; }
			if (l.First==r.First && l.Second>r.Second) { return true; }
			return false;
		}
		friend bool operator<(const Elems l, const Elems r) {
			if (l.First<r.First) { return true; }
			if (l.First==r.First && l.Second<r.Second) { return true; }
			return false;
		}
		friend bool operator>=(const Elems l, const Elems r) {
			if (l.First>r.First) { return true; }
			if (l.First==r.First && l.Second>=r.Second) { return true; }
			return false;
		}
		friend bool operator<=(const Elems l, const Elems r) {
			if (l.First<r.First) { return true; }
			if (l.First==r.First && l.Second<=r.Second) { return true; }
			return false;
		}
		friend bool operator!=(const Elems l, const Elems r) {
			return l.First != r.First && l.Second!=r.Second;
		}
		friend bool operator==(const Elems l, const Elems r) {
			return l.First == r.First && l.Second==r.Second;
		}

		friend std::ostream& operator<<(std::ostream& os, Elems e) {
			os << "{First: " << e.First << ", Second: " << e.Second << "}";
			return os;
		}
	};

	AssociatedSlices(Slice<T> one, Slice<U> two): first(one), second(two) {
		if (one.Len()!=two.Len()) {
			throw Err::ValuesDidNotMatch{
				.Desc="Associated slice lengths must match",
				.First=one.Len(), .Second=two.Len(),
			};
		}
	}

	inline size_t Len() { return first.Len(); }

	Elems operator[](size_t idx) {
		if (idx>=first.Len()) {
			throw Err::IndexOutOfBounds{ .Len=first.Len(), .Idx=idx };
		}
		return Elems{ .First=this->first[idx], .Second=this->second[idx], };
	}

	friend std::ostream& operator<<(std::ostream& os, AssociatedSlices<T, U> a) {
		os << "AssociatedSlices(" << a.first.Len() << ")[";
		for (size_t i=0; i<a.first.Len(); i++) {
			os << a[i];
			if (i+1<a.first.Len()) {
				os << ", ";
			}
		}
		os << "]";
	    return os;
	}
};

namespace Heap {

template <typename T, typename U>
void heapHelper(T s, size_t curIdx, bool(*cmp)(U l, U r)) {
	size_t largest=curIdx;
	size_t left=2*curIdx+1;
	size_t right=2*curIdx+2;

	if (left<s.Len() && cmp(s[left], s[largest])) {
		largest=left;
	}
	if (right<s.Len() && cmp(s[right], s[largest])) {
		largest=right;
	}

	if (largest!=curIdx) {
		using std::swap;
		swap(s[curIdx], s[largest]);
		heapHelper<T, U>(s, largest, cmp);
	}
}

template <typename T, typename U>
void Max(T s) {
	int startNode=(s.Len()/2)-1;
	auto op=[](U a, U b) { return a>b; };
	for (size_t i=startNode; i>0; i--) {
		heapHelper<T, U>(s, i, op);
	}
	heapHelper<T, U>(s, 0, op);
}

template <typename T, typename U>
void Min(T s) {
	int startNode=(s.Len()/2)-1;
	auto op=[](U a, U b) { return a<b; };
	for (size_t i=startNode; i>0; i--) {
		heapHelper<T, U>(s, i, op);
	}
	heapHelper<T, U>(s, 0, op);
}

template <typename T>
void Max(Slice<T> s) { Max<Slice<T>, T>(s); }

template <typename T>
void Min(Slice<T> s) { Min<Slice<T>, T>(s); }


};

#endif
