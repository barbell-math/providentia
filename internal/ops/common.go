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
