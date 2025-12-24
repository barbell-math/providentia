package util

import (
	"time"

	"golang.org/x/exp/constraints"
)

func SliceClamp[S []T, T any, U constraints.Integer](s S, _len U) S {
	if len(s) < int(_len) {
		s = make([]T, _len)
	}
	return s[:_len]
}

func DateEqual(date1 time.Time, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2 // Standard comparison of date components
}
