package ops

import (
	"fmt"
	"iter"
	"strings"

	"github.com/barbell-math/providentia/lib/types"
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

func exerciseFocusHelp() string {
	names := types.ExerciseFocusNames()
	values := types.ExerciseFocusValues()

	var sb strings.Builder
	for i := range len(names) {
		sb.WriteString(fmt.Sprintf("%d (a.k.a. %s)", values[i], names[i]))
	}
	return sb.String()
}

func exerciseKindHelp() string {
	names := types.ExerciseKindNames()
	values := types.ExerciseKindValues()

	var sb strings.Builder
	for i := range len(names) {
		sb.WriteString(fmt.Sprintf("%d (a.k.a. %s)", values[i], names[i]))
	}
	return sb.String()
}
