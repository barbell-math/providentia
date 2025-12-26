package dal

import (
	"context"
	"fmt"
	"time"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
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
	WHERE providentia.client.email=$1
)`

	exerciseIdSelectSql = `
(
	SELECT providentia.exercise.id FROM providentia.exercise
	WHERE providentia.exercise.name=$2
)`

	deleteTrainingLogsByIdSql = `
DELETE FROM providentia.training_log
USING providentia.client
WHERE
	providentia.client.id = providentia.training_log.client_id AND
	providentia.client.email = $1 AND
	providentia.training_log.inter_session_cntr = $2 AND
	providentia.training_log.date_performed = $3;
`

	deleteTrainingLogsBetweenDatesSql = `
WITH deleted_training_logs AS (
	DELETE FROM providentia.training_log
	USING providentia.client
	WHERE
		providentia.client.id = providentia.training_log.client_id AND
		providentia.client.email = $1 AND
		providentia.training_log.date_performed >= $2 AND
		providentia.training_log.date_performed < $3
	RETURNING
		providentia.client.id,
		providentia.training_log.inter_session_cntr,
		providentia.training_log.date_performed
) SELECT COUNT(*) FROM deleted_training_logs
GROUP BY id, inter_session_cntr, date_performed;
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
				*res = make([]any, 9)
				(*res)[0] = v.Val.ClientEmail
				(*res)[1] = v.Val.ExerciseName
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

func deleteTrainingLogsById(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	ids []types.WorkoutId,
) error {
	for start, end := range batchIndexes(ids, int(state.Global.BatchSize)) {
		select {
		case <-ctxt.Done():
			return ctxt.Err()
		default:
		}

		b := pgx.Batch{}
		for i := start; i < end; i++ {
			b.Queue(
				deleteTrainingLogsByIdSql,
				ids[i].ClientEmail, ids[i].Session, ids[i].DatePerformed,
			)
		}
		results := tx.SendBatch(ctxt, &b)

		for i := start; i < end; i++ {
			if cmdTag, err := results.Exec(); err != nil {
				results.Close()
				return sberr.AppendError(
					types.CouldNotDeleteAllTrainingLogsErr, err,
				)
			} else if cmdTag.RowsAffected() == 0 {
				results.Close()
				return sberr.Wrap(
					types.CouldNotDeleteAllTrainingLogsErr,
					"Could not delete entry with id '%+v' (Does id exist?)",
					ids[i],
				)
			}
		}
		results.Close()

		state.Log.Log(
			ctxt, sblog.VLevel(3),
			"DAL: Deleted training_log entries",
			"NumWorkouts", end-start,
		)
	}

	return nil
}

func deleteTrainingLogsInDateRange(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	clientEmail string,
	start time.Time,
	end time.Time,
	res *int64,
) error {
	row := tx.QueryRow(
		ctxt, deleteTrainingLogsBetweenDatesSql, clientEmail, start, end,
	)
	if err := row.Scan(res); err != nil {
		return sberr.AppendError(types.CouldNotDeleteAllTrainingLogsErr, err)
	}
	state.Log.Log(
		ctxt, sblog.VLevel(3),
		fmt.Sprintf(
			"DAL: Deleted training_log entries in date range (%s, %s]",
			start, end,
		),
		"NumWorkouts", *res,
	)
	return nil
}
