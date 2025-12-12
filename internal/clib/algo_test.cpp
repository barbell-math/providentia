#include <stdbool.h>
#include "./tests.gen.h"
#include "./asserts.gen.h"
#include "./dataStructs.h"
#include "./algo.h"

extern "C" bool TestMaxHeap(void) {
	double data[10]={};
	for (int i=0; i<10; i++) {
		data[i]=i;
	}
	
	Slice<double> s(data, 10);
	Heap::Max(s);
	
	double tmp[10]={9, 8, 6, 7, 4, 5, 2, 0, 3, 1};
	for (int i=0; i<10; i++) {
		EQ(data[i], tmp[i]);
		EQ(s[i], tmp[i]);
	}

	return true;
}

extern "C" bool TestMinHeap(void) {
	double data[10]={};
	for (int i=0; i<10; i++) {
		data[i]=10-i;
	}
	
	Slice<double> s(data, 10);
	Heap::Min(s);

	double tmp[10]={1, 2, 4, 3, 6, 5, 8, 10, 7, 9};
	for (int i=0; i<10; i++) {
		EQ(data[i], tmp[i]);
		EQ(s[i], tmp[i]);
	}

	return true;
}
