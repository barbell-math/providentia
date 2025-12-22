package dal

import (
	"context"
	"time"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	"github.com/jackc/pgx/v5"
)

type (
	trainingLog struct {
		ClientEmail      string    `db:"client_email"`
		Session          uint16    `db:"session"`
		DatePerformed    time.Time `db:"date_performed"`
		InterSessionCntr int16
		InterWorkoutCntr int16
		ExerciseName     string
		Weight           types.Kilogram `db:"weight"`
		Sets             float64        `db:"sets"`
		Reps             int32          `db:"reps"`
		Effort           types.RPE      `db:"effort"`
	}
)

const (
	trainingLogTableName = "training_log"
)

func createTrainingLogsReturningIds(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	data []genericCreateReturningIdVal[trainingLog],
) error {
	return nil
}
