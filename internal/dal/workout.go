package dal

import (
	"context"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	"github.com/jackc/pgx/v5"
)

func CreateWorkouts(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	workouts []types.Workout,
) error {
	// trainingLogToPhysicsIdMapping struct list
	// createPhysicsRetIds idVal list -> point to id mapping list
	// createTrainingLogsRetIds idVal list -> point to id mapping list

	// patch up pointers (is there a way to know allocation size before hand?)
	// how to patch up pointers????

	// create all physics and training log entries to have id vals set
	// create all mapping entries

	return nil
}
