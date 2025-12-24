package dal

import (
	"context"
	"time"

	"code.barbellmath.net/barbell-math/providentia/internal/util"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbarena "code.barbellmath.net/barbell-math/smoothbrain-arena"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	"github.com/jackc/pgx/v5"
	"golang.org/x/sync/errgroup"
)

// TODO - make returned errors match other errors...
func CreateWorkouts(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	workouts []types.Workout,
) error {
	type physDataRes = genericCreateReturningIdVal[*types.PhysicsData]
	type trainingLogRes = genericCreateReturningIdVal[trainingLog]

	physicsArena := sbarena.NewTypedArena[physDataRes](
		int(state.Global.BatchSize),
	)
	trainingLogArena := sbarena.NewTypedArena[trainingLogRes](
		int(state.Global.BatchSize),
	)
	tlToPdArena := sbarena.NewTypedArena[trainingLogToPhysicsData](
		int(state.Global.BatchSize),
	)

	curDate := time.Time{}
	interSessionCntr := int16(0)
	for _, w := range workouts {
		select {
		case <-ctxt.Done():
			return ctxt.Err()
		default:
		}

		if util.DateEqual(w.DatePerformed, curDate) {
			interSessionCntr = -1
			curDate = w.DatePerformed
		}
		interSessionCntr++

		for interWorkoutCntr, e := range w.Exercises {
			iterTl := trainingLogArena.Alloc()
			*iterTl = trainingLogRes{
				Val: trainingLog{
					ClientEmail:      w.ClientEmail,
					ExerciseName:     e.Name,
					DatePerformed:    w.DatePerformed,
					InterSessionCntr: interSessionCntr,
					InterWorkoutCntr: int16(interWorkoutCntr + 1),
					Weight:           e.Weight,
					Sets:             e.Sets,
					Reps:             e.Reps,
					Effort:           e.Effort,
				},
			}

			for setNum, p := range e.PhysData {
				if p.Present {
					iterPd := physicsArena.Alloc()
					*iterPd = physDataRes{
						Val: &p.Value,
					}

					iterTlToPd := tlToPdArena.Alloc()
					*iterTlToPd = trainingLogToPhysicsData{
						TrainingLogId: &iterTl.Id,
						PhysicsId:     &iterPd.Id,
						SetNum:        int32(setNum),
					}
				}
			}
		}
	}

	errG := errgroup.Group{}
	errG.Go(func() error {
		for _, c := range physicsArena.Chunks() {
			if err := createPhysicsDataReturningIds(
				ctxt, state, tx, c,
			); err != nil {
				return err
			}
		}
		return nil
	})

	errG.Go(func() error {
		for _, c := range trainingLogArena.Chunks() {
			if err := createTrainingLogsReturningIds(
				ctxt, state, tx, c,
			); err != nil {
				return err
			}
		}
		return nil
	})

	if err := errG.Wait(); err != nil {
		return sberr.AppendError(types.CouldNotCreateAllWorkoutsErr, err)
	}

	for _, c := range tlToPdArena.Chunks() {
		if err := createTrainingLogToPhysicsMappings(
			ctxt, state, tx, c,
		); err != nil {
			return sberr.AppendError(types.CouldNotCreateAllWorkoutsErr, err)
		}
	}

	return nil
}
