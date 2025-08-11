package logic

import (
	"context"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/internal/ops"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
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
