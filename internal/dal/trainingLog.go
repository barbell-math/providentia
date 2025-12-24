package dal

import (
	"context"
	"time"

	"code.barbellmath.net/barbell-math/providentia/internal/util"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	"github.com/jackc/pgx/v5"
)

type (
	trainingLog struct {
		ClientEmail      string
		ExerciseName     string
		DatePerformed    time.Time      `db:"date_performed"`
		InterSessionCntr int16          `db:"inter_session_cntr"`
		InterWorkoutCntr int16          `db:"inter_workout_cntr"`
		Weight           types.Kilogram `db:"weight"`
		Sets             float64        `db:"sets"`
		Reps             int32          `db:"reps"`
		Effort           types.RPE      `db:"effort"`
	}
)

const (
	trainingLogTableName = "training_log"

	clientIdSelectSql = `
(
	SELECT providentia.client.id FROM providentia.client
	WHERE providentia.client.email='$1'
)
`

	exerciseIdSelectSql = `
(
	SELECT providentia.exercise.id FROM providentia.exercise
	WHERE providentia.exercise.name='$2'
)
`
)

func createTrainingLogsReturningIds(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	data []genericCreateReturningIdVal[trainingLog],
) error {
	return genericCreateReturningId(
		ctxt, state, tx, &genericCreateReturningIdOpts[trainingLog]{
			TableName: trainingLogTableName,
			Columns: []string{
				"client_id", "exercise_id",
				"inter_session_cntr", "inter_workout_cntr",
				"date_performed", "weight", "sets", "reps", "effort",
			},
			ValueGetter: func(
				v *genericCreateReturningIdVal[trainingLog],
				res *[]any,
			) error {
				*res = util.SliceClamp(*res, 9)
				(*res)[0] = v.Val.ExerciseName
				(*res)[1] = v.Val.ClientEmail
				(*res)[2] = v.Val.InterSessionCntr
				(*res)[3] = v.Val.InterWorkoutCntr
				(*res)[4] = v.Val.DatePerformed
				(*res)[5] = v.Val.Weight
				(*res)[6] = v.Val.Sets
				(*res)[7] = v.Val.Reps
				(*res)[8] = v.Val.Effort
				return nil
			},
			ModifyValuePlaceholders: func(placeholders []string) []string {
				placeholders[0] = clientIdSelectSql
				placeholders[1] = exerciseIdSelectSql
				return placeholders
			},
			Data: data,
			Err:  types.CouldNotCreateAllTrainingLogsErr,
		},
	)
}
