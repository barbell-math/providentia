package jobs

import (
	"context"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbjobqueue "code.barbellmath.net/barbell-math/smoothbrain-jobQueue"
)

type (
	GP[T any] struct {
		S *types.State
		Q *dal.SyncQueries
		B *sbjobqueue.Batch
		V T
		F func(context.Context, *types.State, *dal.SyncQueries, T) error
	}
)

func (g *GP[T]) JobType(_ types.GeneralPurposeJob) {}

func (g *GP[T]) Batch() *sbjobqueue.Batch {
	return g.B
}

func (g *GP[T]) Run(ctxt context.Context) error {
	return g.F(ctxt, g.S, g.Q, g.V)
}
