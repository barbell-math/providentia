package ops

import (
	"iter"
)

func batchIndexes[S ~[]E, E any](s S, step int) iter.Seq2[int, int] {
	return func(yield func(int, int) bool) {
		for i := 0; i < len(s); i += step {
			end := i + step
			if !yield(i, min(len(s), end)) {
				break
			}
		}
	}
}

func flatten[S ~[]E, E any](s ...S) iter.Seq2[int, E] {
	return func(yield func(int, E) bool) {
		for i := range len(s) {
			for j := range len(s[i]) {
				if !yield(j, s[i][j]) {
					break
				}
			}
		}
	}
}
