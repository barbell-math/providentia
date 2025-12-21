package util

import "golang.org/x/exp/constraints"

func SliceClamp[S []T, T any, U constraints.Integer](s S, _len U) S {
	if len(s) < int(_len) {
		s = make([]T, _len)
	}
	return s[:_len]
}
