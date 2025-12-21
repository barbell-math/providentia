#ifndef CGO_GLUE_ALGO
#define CGO_GLUE_ALGO

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
