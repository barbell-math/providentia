package dal

import (
	"context"
	"fmt"
	"iter"
	"strings"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	"github.com/jackc/pgx/v5"
)

type (
	AvailableTypes interface {
		types.Client | types.Exercise | types.Hyperparams // TODO - add more types here as created
	}

	CreateFunc[T AvailableTypes] func(
		ctxt context.Context,
		state *types.State,
		tx pgx.Tx,
		clients []T,
	) error

	CpyFromSlice[T any] struct {
		Data        []T
		ValueGetter func(v *T, res *[]any) error

		curIdx  int
		curVals []any
		err     error
	}
)

func (c *CpyFromSlice[T]) Next() bool { return c.curIdx < len(c.Data) }
func (c *CpyFromSlice[T]) Values() ([]any, error) {
	c.err = c.ValueGetter(&c.Data[c.curIdx], &c.curVals)
	c.curIdx++
	return c.curVals, c.err
}
func (c *CpyFromSlice[T]) Err() error { return c.err }

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

func defaultValuePlaceholder(num int) string {
	return fmt.Sprintf("$%d", num)
}

func defaultValuePlaceholders(n int) []string {
	rv := make([]string, n)
	for i := 1; i < n+1; i++ {
		rv[i-1] = fmt.Sprintf("$%d", i)
	}
	return rv
}

func defaultValuePlaceholdersJoined(n int) string {
	var sb strings.Builder
	for i := range n {
		sb.WriteString(defaultValuePlaceholder(i + 1))
		if i+1 < n {
			sb.WriteString(", ")
		}
	}
	return sb.String()
}
