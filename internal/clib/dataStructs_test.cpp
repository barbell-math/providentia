#include <cstddef>
#include <stdbool.h>
#include "./tests.gen.h"
#include "./asserts.gen.h"
#include "./dataStructs.h"
#include "./errors.h"

extern "C" bool TestSlicePointerConstructor(void) {
	double data[3]={};
	for (int i=0; i<3; i++) {
		data[i]=i;
	}

	Slice<double> s(data, 3);
	EQ(s.Len(), (size_t)3)
	for (int i=0; i<3; i++) {
		EQ(&data[i], &s.data[i]);
	}

	s=Slice<double>(&data[1], 2);
	EQ(s.Len(), (size_t)2)
	for (int i=0; i<2; i++) {
		EQ(&data[i+1], &s.data[i]);
	}

	return true;
}

extern "C" bool TestSliceBracketOperator(void) {
	double data[3]={};
	for (int i=0; i<3; i++) {
		data[i]=i;
	}

	Slice<double> s(data, 3);
	EQ(s.Len(), (size_t)3)
	for (int i=0; i<3; i++) {
		EQ(data[i], s[i]);
		EQ(&data[i], &s[i]);
	}

	THROWS(Err::IndexOutOfBounds, s[-1];)
	THROWS(Err::IndexOutOfBounds, s[3];)

	return true;
}

extern "C" bool TestSliceSubslice(void) {
	double data[3]={};
	for (int i=0; i<3; i++) {
		data[i]=i;
	}

	Slice<double> s(data, 3);

	Slice<double> subSlice=s(0,1);
	EQ(subSlice.Len(), (size_t)1)
	for (size_t i=0; i<subSlice.Len(); i++) {
		EQ(data[i], s[i]);
		EQ(&data[i], &s[i]);
	}

	subSlice=s(1,2);
	EQ(subSlice.Len(), (size_t)1)
	for (size_t i=1; i<subSlice.Len(); i++) {
		EQ(data[i], s[i]);
		EQ(&data[i], &s[i]);
	}

	subSlice=s(2,3);
	EQ(subSlice.Len(), (size_t)1)
	for (size_t i=2; i<subSlice.Len(); i++) {
		EQ(data[i], s[i]);
		EQ(&data[i], &s[i]);
	}

	subSlice=s(1,3);
	EQ(subSlice.Len(), (size_t)2)
	for (size_t i=1; i<subSlice.Len(); i++) {
		EQ(data[i], s[i]);
		EQ(&data[i], &s[i]);
	}

	subSlice=s(0,3);
	EQ(subSlice.Len(), (size_t)3)
	for (size_t i=0; i<subSlice.Len(); i++) {
		EQ(data[i], s[i]);
		EQ(&data[i], &s[i]);
	}

	subSlice=s(1,1);
	EQ(subSlice.Len(), (size_t)0)
	for (size_t i=0; i<subSlice.Len(); i++) {
		EQ(data[i], s[i]);
		EQ(&data[i], &s[i]);
	}

	THROWS(Err::IndexOutOfBounds, subSlice=s(-1,3);)
	THROWS(Err::IndexOutOfBounds, subSlice=s(0,4);)
	THROWS(Err::StartAfterEnd<size_t>, subSlice=s(2,1);)

	return true;
}

extern "C" bool TestSliceIterator(void) {
	double data[3]={};
	for (int i=0; i<3; i++) {
		data[i]=i;
	}

	Slice<double> s(data, 3);

	size_t i=0;
	for (double& it : s) {
		EQ(data[i], s[i]);
		EQ(&data[i], &s[i]);
		EQ(it, s[i]);
		i++;
	}
	EQ(i, s.Len());

	i=0;
	for (const double& it : s) {
		EQ(data[i], s[i]);
		EQ(&data[i], &s[i]);
		EQ(it, s[i]);
		i++;
	}
	EQ(i, s.Len());

	return true;
}

extern "C" bool TestFixedSlicePointerConstructor(void) {
	double data[3]={};
	for (int i=0; i<3; i++) {
		data[i]=i;
	}

	FixedSlice<double, 3> a(data);
	EQ(a.Len(), (size_t)3)
	for (int i=0; i<3; i++) {
		EQ(&data[i], &a.data[i]);
	}

	FixedSlice<double, 2> a2(&data[1]);
	EQ(a2.Len(), (size_t)2)
	for (int i=0; i<2; i++) {
		EQ(&data[i+1], &a2.data[i]);
	}

	return true;
}

extern "C" bool TestFixedSliceWholeSliceConstructor(void) {
	double data[3]={};
	for (int i=0; i<3; i++) {
		data[i]=i;
	}
	Slice<double> s(data, 3);

	FixedSlice<double, 3> a(s);
	EQ(a.Len(), (size_t)3)
	for (int i=0; i<3; i++) {
		EQ(&data[i], &a.data[i]);
	}

	FixedSlice<double, 2> a2(s(1,3));
	EQ(a2.Len(), (size_t)2)
	for (int i=0; i<2; i++) {
		EQ(&data[i+1], &a2.data[i]);
	}

	THROWS(Err::IndexOutOfBounds, FixedSlice<double, 4>tmp(s);)

	return true;
}

extern "C" bool TestFixedSliceSubSliceConstructor(void) {
	double data[3]={};
	for (int i=0; i<3; i++) {
		data[i]=i;
	}
	Slice<double> s(data, 3);

	FixedSlice<double, 2> a(s, 0);
	EQ(a.Len(), (size_t)2)
	for (int i=0; i<2; i++) {
		EQ(&data[i], &a.data[i]);
	}

	FixedSlice<double, 2> a2(s, 1);
	EQ(a2.Len(), (size_t)2)
	for (int i=1; i<3; i++) {
		EQ(&data[i+1], &a2.data[i]);
	}

	THROWS(Err::IndexOutOfBounds, FixedSlice<double, 2>tmp(s, 2);)
	THROWS(Err::IndexOutOfBounds, FixedSlice<double, 4>tmp(s, 0);)

	return true;
}

extern "C" bool TestFixedSliceBracketOperator(void) {
	double data[3]={};
	for (int i=0; i<3; i++) {
		data[i]=i;
	}

	FixedSlice<double, 3> a(data);
	EQ(a.Len(), (size_t)3)
	for (int i=0; i<3; i++) {
		EQ(data[i], a[i]);
		EQ(&data[i], &a[i]);
	}

	THROWS(Err::IndexOutOfBounds, a[-1];)
	THROWS(Err::IndexOutOfBounds, a[3];)

	return true;
}

extern "C" bool TestFixedSliceSubslice(void) {
	double data[3]={};
	for (int i=0; i<3; i++) {
		data[i]=i;
	}

	FixedSlice<double, 3> a(data);

	Slice<double> subSlice=a(0,1);
	EQ(subSlice.Len(), (size_t)1)
	for (size_t i=0; i<subSlice.Len(); i++) {
		EQ(data[i], a[i]);
		EQ(&data[i], &a[i]);
	}

	subSlice=a(1,2);
	EQ(subSlice.Len(), (size_t)1)
	for (size_t i=1; i<subSlice.Len(); i++) {
		EQ(data[i], a[i]);
		EQ(&data[i], &a[i]);
	}

	subSlice=a(2,3);
	EQ(subSlice.Len(), (size_t)1)
	for (size_t i=2; i<subSlice.Len(); i++) {
		EQ(data[i], a[i]);
		EQ(&data[i], &a[i]);
	}

	subSlice=a(1,3);
	EQ(subSlice.Len(), (size_t)2)
	for (size_t i=1; i<subSlice.Len(); i++) {
		EQ(data[i], a[i]);
		EQ(&data[i], &a[i]);
	}

	subSlice=a(0,3);
	EQ(subSlice.Len(), (size_t)3)
	for (size_t i=0; i<subSlice.Len(); i++) {
		EQ(data[i], a[i]);
		EQ(&data[i], &a[i]);
	}

	subSlice=a(1,1);
	EQ(subSlice.Len(), (size_t)0)
	for (size_t i=0; i<subSlice.Len(); i++) {
		EQ(data[i], a[i]);
		EQ(&data[i], &a[i]);
	}

	THROWS(Err::IndexOutOfBounds, subSlice=a(-1,3);)
	THROWS(Err::IndexOutOfBounds, subSlice=a(0,4);)
	THROWS(Err::StartAfterEnd<size_t>, subSlice=a(2,1);)

	return true;
}

extern "C" bool TestFixedSliceIterator(void) {
	double data[3]={};
	for (int i=0; i<3; i++) {
		data[i]=i;
	}

	FixedSlice<double, 3> s(data);

	size_t i=0;
	for (double& it : s) {
		EQ(data[i], s[i]);
		EQ(&data[i], &s[i]);
		EQ(it, s[i]);
		i++;
	}
	EQ(i, s.Len());

	i=0;
	for (const double& it : s) {
		EQ(data[i], s[i]);
		EQ(&data[i], &s[i]);
		EQ(it, s[i]);
		i++;
	}
	EQ(i, s.Len());

	return true;
}

extern "C" bool TestFixedRingPointerConstructor(void) {
	double data[3]={};
	for (int i=0; i<3; i++) {
		data[i]=i;
	}

	FixedRing<double, 3> r(data);
	EQ(r.Len(), (size_t)3)
	for (int i=0; i<3; i++) {
		EQ(&data[i], &r.data[i]);
	}

	return true;
}

extern "C" bool TestFixedRingFixedSliceConstructor(void) {
	double data[3]={};
	for (int i=0; i<3; i++) {
		data[i]=i;
	}
	FixedSlice<double, 3> s(data);

	FixedRing<double, 3> r(s);
	EQ(r.Len(), (size_t)3)
	for (int i=0; i<3; i++) {
		EQ(&data[i], &r.data[i]);
	}

	return true;
}

extern "C" bool TestFixedRingSliceConstructor(void) {
	double data[3]={};
	for (int i=0; i<3; i++) {
		data[i]=i;
	}
	Slice<double> s(data, 3);

	FixedRing<double, 3> r(s);
	EQ(r.Len(), (size_t)3)
	for (int i=0; i<3; i++) {
		EQ(&data[i], &r.data[i]);
	}

	return true;
}

extern "C" bool TestFixedRingSubSliceConstructor(void) {
	double data[3]={};
	for (int i=0; i<3; i++) {
		data[i]=i;
	}
	Slice<double> s(data, 3);

	FixedRing<double, 2> a(s, 0);
	EQ(a.Len(), (size_t)2)
	for (int i=0; i<2; i++) {
		EQ(&data[i], &a.data[i]);
	}

	FixedRing<double, 2> a2(s, 1);
	EQ(a2.Len(), (size_t)2)
	for (int i=1; i<3; i++) {
		EQ(&data[i+1], &a2.data[i]);
	}

	THROWS(Err::IndexOutOfBounds, FixedRing<double, 2>tmp(s, 2);)
	THROWS(Err::IndexOutOfBounds, FixedRing<double, 4>tmp(s, 0);)

	return true;
}

extern "C" bool TestFixedRingPut(void) {
	double data[3]={};
	for (int i=0; i<3; i++) {
		data[i]=i;
	}

	FixedRing<double, 3> r(data);
	for (double i=3; i<6; i++) {
		r.Put(i);
	}
	for (double i=3; i<6; i++) {
		EQ(data[(size_t)i-3], i);
		EQ(r[(size_t)i-3], i);
	}

	for (double i=6; i<9; i++) {
		r.Put(i);
	}
	for (double i=6; i<9; i++) {
		EQ(data[(size_t)i-6], i);
		EQ(r[(size_t)i-6], i);
	}

	return true;
}

extern "C" bool TestFixedRingIterator(void) {
	double data[3]={};
	for (int i=0; i<3; i++) {
		data[i]=i;
	}

	FixedRing<double, 3> r(data);

	size_t i=0;
	for (double& it : r) {
		EQ(data[i], r[i]);
		EQ(&data[i], &r[i]);
		EQ(it, r[i]);
		i++;
	}
	EQ(i, r.Len());

	i=0;
	for (const double& it : r) {
		EQ(data[i], r[i]);
		EQ(&data[i], &r[i]);
		EQ(it, r[i]);
		i++;
	}
	EQ(i, r.Len());

	return true;
}

extern "C" bool TestAssociatedSlicesConstructor(void) {
	double data[10]={};
	int data2[10]={};
	for (int i=0; i<10; i++) {
		data[i]=i;
		data2[i]=i;
	}
	
	Slice<double> s1(data, 10);
	Slice<int> s2(data2, 10);
	AssociatedSlices<double, int> a(s1, s2);

	for (int i=0; i<10; i++) {
		EQ(a[i].First, (double)i);
		EQ(a[i].Second, i);
		EQ(&a[i].First, &data[i]);
		EQ(&a[i].Second, &data2[i]);
	}

	a[0].First=10;
	a[0].Second=10;
	EQ(a[0].First, data[0]);
	EQ(a[0].Second, data2[0]);

	return true;
}
