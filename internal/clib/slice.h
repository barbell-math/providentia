#ifndef CGO_GLUE_SLICE
#define CGO_GLUE_SLICE

#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include "common.h"

template <typename  T>
struct Slice {
	T *Data;
	size_t Len;

	Slice(void *data, int len): Data((T*)data), Len(len) {}

	T& operator[](size_t idx) {
		if (idx<0 || idx>=this->Len) {
			char err[73];
			sprintf(err, "Index out of bounds: Len: %lu Idx: %lu", this->Len, idx);
			Funcs::panic(err);
		}
		return this->Data[idx];
	}

	Slice<T> operator()(size_t start, size_t end) {
		if (start<0 || start>=this->Len) {
			char err[73];
			sprintf(err, "Index out of bounds: Len: %lu Idx: %lu", this->Len, start);
			Funcs::panic(err);
		}
		if (end>this->Len) {
			end=this->Len;
		}
		return Slice<T>(&this->Data[start], end-start);
	}

	void free() {
		free(this->Data);
		this->Len=0;
	}
};

template<typename  T, size_t N>
struct Array : Slice<T> {
	Array(void *data): Slice<T>(data, N) {}
	Array(Slice<T> s): Slice<T>(s.Data, N) {
		if (s.Len<N) {
			char err[73];
			sprintf(err, "Index out of bounds: Len: %lu Idx: %lu", s.Len, N);
			Funcs::panic(err);
		}
		this->Len=N;
		this->Data=s.Data;
	}
	Array(Slice<T> s, size_t start): Slice<T>(s.Data, N) {
		if (s.Len-start<N) {
			char err[73];
			sprintf(err, "Index out of bounds: Len: %lu Idx: %lu", s.Len-start, N);
			Funcs::panic(err);
		}
		this->Len=N;
		this->Data=&s.Data[start];
	}
};

#endif
