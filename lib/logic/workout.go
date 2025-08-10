package logic

import (
	"context"

	dal "github.com/barbell-math/providentia/internal/db/dataAccessLayer"
	"github.com/barbell-math/providentia/internal/ops"
	"github.com/barbell-math/providentia/lib/types"
)

// TODO - add comment
func CreateWorkouts(
	ctxt context.Context,
	workouts ...types.Workout,
) (opErr error) {
	if len(workouts) == 0 {
		return
	}
	return runOp(ctxt, func(state *types.State, queries *dal.Queries) error {
		return ops.CreateWorkouts(ctxt, state, queries, workouts...)
	})
}
