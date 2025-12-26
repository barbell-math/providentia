package dal

import (
	"context"
	"fmt"
	"math"
	"time"
	"unsafe"

	"code.barbellmath.net/barbell-math/providentia/internal/util"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbarena "code.barbellmath.net/barbell-math/smoothbrain-arena"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
	"github.com/jackc/pgx/v5"
)

type (
	ReadNumWorkoutsForClientOpts struct {
		Email string
		Res   *int64
	}

	ReadWorkoutsByIdOpts struct {
		Ids []types.WorkoutId
		Res *[]types.Workout
	}

	FindWorkoutsByIdOpts struct {
		Ids []types.WorkoutId
		Res *[]types.Optional[types.Workout]
	}

	FindWorkoutsInDateRangeOpts struct {
		Email string
		Start time.Time
		End   time.Time
		Res   *[]types.Workout
	}

	DeleteWorkoutsInDateRangeOpts struct {
		Email string
		Start time.Time
		End   time.Time
		Res   *int64
	}

	readWorkoutSqlResult struct {
		ExerciseName string
		Weight       types.Kilogram
		Sets         float64
		CurSet       int
		Reps         int32
		Effort       types.RPE
		types.AbstractData
		types.PhysicsData
	}

	findworkoutBetweenDatesSqlResult struct {
		DatePerformed time.Time
		Session       uint16
		ExerciseName  string
		Weight        types.Kilogram
		Sets          float64
		CurSet        int
		Reps          int32
		Effort        types.RPE
		types.AbstractData
		types.PhysicsData
	}
)

const (
	readNumWorkoutsForClientSql = `
SELECT COUNT(*) FROM (
	SELECT date_performed, inter_session_cntr
	FROM providentia.training_log
	JOIN providentia.client
		ON providentia.training_log.client_id = providentia.client.id
	WHERE providentia.client.email = $1
	GROUP BY date_performed, inter_session_cntr
) AS result;
`

	workoutByIdSql = `
SELECT
	providentia.exercise.name,
	providentia.training_log.weight,
	providentia.training_log.sets,
	COALESCE(providentia.training_log_to_physics_data.set_num+1, 0) AS cur_set,
	providentia.training_log.reps,
	providentia.training_log.effort,
	providentia.training_log.volume,
	providentia.training_log.exertion,
	providentia.training_log.total_reps,
	COALESCE(providentia.physics_data.bar_path_calc_id, 0),
	COALESCE(providentia.physics_data.bar_path_track_id, 0),
	COALESCE(providentia.physics_data.path, ''),
	providentia.physics_data.time,
	providentia.physics_data.position,
	providentia.physics_data.velocity,
	providentia.physics_data.acceleration,
	providentia.physics_data.jerk,
	providentia.physics_data.force,
	providentia.physics_data.impulse,
	providentia.physics_data.work,
	providentia.physics_data.power,
	providentia.physics_data.rep_splits,
	providentia.physics_data.min_vel,
	providentia.physics_data.max_vel,
	providentia.physics_data.max_acc,
	providentia.physics_data.max_acc,
	providentia.physics_data.max_force,
	providentia.physics_data.max_force,
	providentia.physics_data.max_impulse,
	providentia.physics_data.max_impulse,
	providentia.physics_data.avg_work,
	providentia.physics_data.max_work,
	providentia.physics_data.max_work,
	providentia.physics_data.avg_power,
	providentia.physics_data.max_power,
	providentia.physics_data.max_power
FROM providentia.training_log
JOIN providentia.client
	ON providentia.client.id = providentia.training_log.client_id
JOIN providentia.exercise
	ON providentia.exercise.id = providentia.training_log.exercise_id
LEFT JOIN providentia.training_log_to_physics_data
	ON providentia.training_log_to_physics_data.training_log_id = providentia.training_log.id
LEFT JOIN providentia.physics_data
	ON providentia.training_log_to_physics_data.physics_id = providentia.physics_data.id
WHERE
	email = $1 AND
	inter_session_cntr = $2 AND
	date_performed = $3
ORDER BY inter_workout_cntr, cur_set ASC;
`

	workoutsBetweenDatesSql = `
SELECT
	providentia.training_log.date_performed,
	providentia.training_log.inter_session_cntr,
	providentia.exercise.name,
	providentia.training_log.weight,
	providentia.training_log.sets,
	COALESCE(providentia.training_log_to_physics_data.set_num+1, 0) AS cur_set,
	providentia.training_log.reps,
	providentia.training_log.effort,
	providentia.training_log.volume,
	providentia.training_log.exertion,
	providentia.training_log.total_reps,
	COALESCE(providentia.physics_data.bar_path_calc_id, 0),
	COALESCE(providentia.physics_data.bar_path_track_id, 0),
	COALESCE(providentia.physics_data.path, ''),
	providentia.physics_data.time,
	providentia.physics_data.position,
	providentia.physics_data.velocity,
	providentia.physics_data.acceleration,
	providentia.physics_data.jerk,
	providentia.physics_data.force,
	providentia.physics_data.impulse,
	providentia.physics_data.work,
	providentia.physics_data.power,
	providentia.physics_data.rep_splits,
	providentia.physics_data.min_vel,
	providentia.physics_data.max_vel,
	providentia.physics_data.max_acc,
	providentia.physics_data.max_acc,
	providentia.physics_data.max_force,
	providentia.physics_data.max_force,
	providentia.physics_data.max_impulse,
	providentia.physics_data.max_impulse,
	providentia.physics_data.avg_work,
	providentia.physics_data.max_work,
	providentia.physics_data.max_work,
	providentia.physics_data.avg_power,
	providentia.physics_data.max_power,
	providentia.physics_data.max_power
FROM providentia.training_log
JOIN providentia.client
	ON providentia.client.id = providentia.training_log.client_id
JOIN providentia.exercise
	ON providentia.exercise.id = providentia.training_log.exercise_id
LEFT JOIN providentia.training_log_to_physics_data
	ON providentia.training_log_to_physics_data.training_log_id = providentia.training_log.id
LEFT JOIN providentia.physics_data
	ON providentia.training_log_to_physics_data.physics_id = providentia.physics_data.id
WHERE
	email = $1 AND
	date_performed >= $2 AND
	date_performed < $3
ORDER BY date_performed, inter_workout_cntr, cur_set ASC;
`
)

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

	for _, w := range workouts {
		select {
		case <-ctxt.Done():
			return ctxt.Err()
		default:
		}

		for interWorkoutCntr, e := range w.Exercises {
			iterTl := trainingLogArena.Alloc()
			*iterTl = trainingLogRes{
				Val: trainingLog{
					ClientEmail:      w.ClientEmail,
					ExerciseName:     e.Name,
					DatePerformed:    w.DatePerformed,
					InterSessionCntr: int16(w.Session),
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

	for _, c := range physicsArena.Chunks() {
		if err := createPhysicsDataReturningIds(ctxt, state, tx, c); err != nil {
			return sberr.AppendError(types.CouldNotCreateAllWorkoutsErr, err)
		}
	}

	for _, c := range trainingLogArena.Chunks() {
		if err := createTrainingLogsReturningIds(ctxt, state, tx, c); err != nil {
			return sberr.AppendError(types.CouldNotCreateAllWorkoutsErr, err)
		}
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

func ReadNumWorkoutsForClient(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts ReadNumWorkoutsForClientOpts,
) error {
	row := tx.QueryRow(ctxt, readNumWorkoutsForClientSql, opts.Email)
	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"DAL: Read total num workouts for client",
	)
	return row.Scan(opts.Res)
}

func ReadWorkoutsById(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts ReadWorkoutsByIdOpts,
) error {
	*opts.Res = util.SliceClamp(*opts.Res, len(opts.Ids))

	for i := range opts.Ids {
		select {
		case <-ctxt.Done():
			return ctxt.Err()
		default:
		}

		iterW := &(*opts.Res)[i]
		foundData, err := readSingleWorkout(ctxt, tx, &opts.Ids[i], iterW)
		if err != nil {
			return sberr.AppendError(types.CouldNotReadAllWorkoutsErr, err)
		}
		if !foundData {
			return sberr.Wrap(
				types.CouldNotReadAllWorkoutsErr,
				"Could not read entry with id '%+v' (Does id exist?)",
				opts.Ids[i],
			)
		}
		iterW.WorkoutId = opts.Ids[i]
	}

	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"DAL: Read workouts by WorkoutIds",
		"NumRows", len(opts.Ids),
	)
	return nil
}

func FindWorkoutsById(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts FindWorkoutsByIdOpts,
) error {
	*opts.Res = util.SliceClamp(*opts.Res, len(opts.Ids))

	found := 0
	for i := range opts.Ids {
		select {
		case <-ctxt.Done():
			return ctxt.Err()
		default:
		}

		iterW := &(*opts.Res)[i]
		foundData, err := readSingleWorkout(
			ctxt, tx, &opts.Ids[i], &iterW.Value,
		)
		if err != nil {
			return sberr.AppendError(types.CouldNotReadAllWorkoutsErr, err)
		}
		if !foundData {
			iterW.Present = false
			continue
		}
		found++
		iterW.Present = true
		iterW.Value.WorkoutId = opts.Ids[i]
	}

	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"DAL: Found workouts by WorkoutIds",
		"NumFound/NumRows", fmt.Sprintf("%d/%d", found, len(opts.Ids)),
	)
	return nil
}

func readSingleWorkout(
	ctxt context.Context,
	tx pgx.Tx,
	id *types.WorkoutId,
	iterW *types.Workout,
) (bool, error) {
	rows, err := tx.Query(
		ctxt,
		workoutByIdSql, id.ClientEmail, id.Session, id.DatePerformed,
	)
	if err != nil {
		return false, err
	}

	var iterE *types.ExerciseData
	for rows.Next() {
		iterResult := readWorkoutSqlResult{}
		if err := rows.Scan(
			&iterResult.ExerciseName,
			&iterResult.Weight,
			&iterResult.Sets,
			&iterResult.CurSet,
			&iterResult.Reps,
			&iterResult.Effort,
			&iterResult.Volume,
			&iterResult.Exertion,
			&iterResult.TotalReps,
			&iterResult.BarPathCalcVersion,
			&iterResult.BarPathTrackerVersion,
			&iterResult.VideoPath,
			&iterResult.Time,
			(*[]genericPoint)(unsafe.Pointer(&iterResult.Position)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.Velocity)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.Acceleration)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.Jerk)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.Force)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.Impulse)),
			&iterResult.Work,
			&iterResult.Power,
			(*[]genericPoint)(unsafe.Pointer(&iterResult.RepSplits)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.MinVel)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.MaxVel)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.MinAcc)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.MaxAcc)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.MinForce)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.MaxForce)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.MinImpulse)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.MaxImpulse)),
			&iterResult.AvgWork,
			(*[]genericPoint)(unsafe.Pointer(&iterResult.MinWork)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.MaxWork)),
			&iterResult.AvgPower,
			(*[]genericPoint)(unsafe.Pointer(&iterResult.MinPower)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.MaxPower)),
		); err != nil {
			rows.Close()
			return false, err
		}

		if iterE == nil || iterE.Name != iterResult.ExerciseName {
			iterW.Exercises = append(
				iterW.Exercises, types.ExerciseData{
					Name:   iterResult.ExerciseName,
					Weight: iterResult.Weight,
					Sets:   iterResult.Sets,
					Reps:   iterResult.Reps,
					Effort: iterResult.Effort,
					AbstractData: types.Optional[types.AbstractData]{
						Present: true,
						Value:   iterResult.AbstractData,
					},
					PhysData: make(
						[]types.Optional[types.PhysicsData],
						int(math.Ceil(iterResult.Sets)),
					),
				},
			)
			iterE = &iterW.Exercises[len(iterW.Exercises)-1]
		}

		if iterResult.CurSet > 0 {
			iterE.PhysData[iterResult.CurSet-1] = types.Optional[types.PhysicsData]{
				Present: len(iterResult.Time) > 0,
				Value:   iterResult.PhysicsData,
			}
		}
	}
	rows.Close()

	return iterE != nil, nil
}

func FindWorkoutsInDateRange(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts FindWorkoutsInDateRangeOpts,
) error {
	if opts.End.Before(opts.Start) {
		return sberr.Wrap(
			types.CouldNotReadAllWorkoutsErr,
			"Start date (%s) must be before end date (%s)",
			opts.Start, opts.End,
		)
	}

	found := 0
	*opts.Res = (*opts.Res)[:0]

	rows, err := tx.Query(
		ctxt,
		workoutsBetweenDatesSql, opts.Email, opts.Start, opts.End,
	)
	if err != nil {
		return sberr.AppendError(types.CouldNotReadAllWorkoutsErr, err)
	}

	var iterW *types.Workout
	var iterE *types.ExerciseData
	for rows.Next() {
		iterResult := findworkoutBetweenDatesSqlResult{}
		if err := rows.Scan(
			&iterResult.DatePerformed,
			&iterResult.Session,
			&iterResult.ExerciseName,
			&iterResult.Weight,
			&iterResult.Sets,
			&iterResult.CurSet,
			&iterResult.Reps,
			&iterResult.Effort,
			&iterResult.Volume,
			&iterResult.Exertion,
			&iterResult.TotalReps,
			&iterResult.BarPathCalcVersion,
			&iterResult.BarPathTrackerVersion,
			&iterResult.VideoPath,
			&iterResult.Time,
			(*[]genericPoint)(unsafe.Pointer(&iterResult.Position)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.Velocity)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.Acceleration)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.Jerk)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.Force)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.Impulse)),
			&iterResult.Work,
			&iterResult.Power,
			(*[]genericPoint)(unsafe.Pointer(&iterResult.RepSplits)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.MinVel)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.MaxVel)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.MinAcc)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.MaxAcc)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.MinForce)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.MaxForce)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.MinImpulse)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.MaxImpulse)),
			&iterResult.AvgWork,
			(*[]genericPoint)(unsafe.Pointer(&iterResult.MinWork)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.MaxWork)),
			&iterResult.AvgPower,
			(*[]genericPoint)(unsafe.Pointer(&iterResult.MinPower)),
			(*[]genericPoint)(unsafe.Pointer(&iterResult.MaxPower)),
		); err != nil {
			rows.Close()
			return sberr.AppendError(types.CouldNotReadAllWorkoutsErr, err)
		}

		iterWorkoutId := types.WorkoutId{
			ClientEmail:   opts.Email,
			Session:       iterResult.Session,
			DatePerformed: iterResult.DatePerformed,
		}
		if iterW == nil || iterW.WorkoutId != iterWorkoutId {
			*opts.Res = append(*opts.Res, types.Workout{
				WorkoutId: iterWorkoutId,
			})
			iterW = &(*opts.Res)[len(*opts.Res)-1]
			found++
		}
		if iterE == nil || iterE.Name != iterResult.ExerciseName {
			iterW.Exercises = append(
				iterW.Exercises, types.ExerciseData{
					Name:   iterResult.ExerciseName,
					Weight: iterResult.Weight,
					Sets:   iterResult.Sets,
					Reps:   iterResult.Reps,
					Effort: iterResult.Effort,
					AbstractData: types.Optional[types.AbstractData]{
						Present: true,
						Value:   iterResult.AbstractData,
					},
					PhysData: make(
						[]types.Optional[types.PhysicsData],
						int(math.Ceil(iterResult.Sets)),
					),
				},
			)
			iterE = &iterW.Exercises[len(iterW.Exercises)-1]
		}

		if iterResult.CurSet > 0 {
			iterE.PhysData[iterResult.CurSet-1] = types.Optional[types.PhysicsData]{
				Present: len(iterResult.Time) > 0,
				Value:   iterResult.PhysicsData,
			}
		}
	}
	rows.Close()

	state.Log.Log(
		ctxt, sblog.VLevel(3),
		fmt.Sprintf(
			"DAL: Found workouts in date range (%s, %s]", opts.Start, opts.End,
		),
		"Found", found,
	)
	return nil
}

func DeleteWorkouts(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	ids []types.WorkoutId,
) error {
	// Note - the order is important here. The physics data has to be deleted
	// first otherwise it will not be able to be found using the provided ids.
	if err := deletePhysicsDataById(ctxt, state, tx, ids); err != nil {
		return sberr.AppendError(types.CouldNotDeleteAllWorkoutsErr, err)
	}
	if err := deleteTrainingLogsById(ctxt, state, tx, ids); err != nil {
		return sberr.AppendError(types.CouldNotDeleteAllWorkoutsErr, err)
	}

	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"DAL: Deleted workout entries",
		"NumRows", len(ids),
	)
	return nil
}

func DeleteWorkoutsInDateRange(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts DeleteWorkoutsInDateRangeOpts,
) error {
	if opts.End.Before(opts.Start) {
		return sberr.Wrap(
			types.CouldNotDeleteAllWorkoutsErr,
			"Start date (%s) must be before end date (%s)",
			opts.Start, opts.End,
		)
	}

	// Note - the order is important here. The physics data has to be deleted
	// first otherwise it will not be able to be found using the provided
	// options.
	var tmp int64
	if err := deletePhysicsDataInDateRange(
		ctxt, state, tx, opts.Email, opts.Start, opts.End, &tmp,
	); err != nil {
		return sberr.AppendError(types.CouldNotDeleteAllWorkoutsErr, err)
	}
	if err := deleteTrainingLogsInDateRange(
		ctxt, state, tx, opts.Email, opts.Start, opts.End, opts.Res,
	); err != nil {
		return sberr.AppendError(types.CouldNotDeleteAllWorkoutsErr, err)
	}

	state.Log.Log(
		ctxt, sblog.VLevel(3),
		fmt.Sprintf(
			"DAL: Deleted workout entries in date range (%s, %s]",
			opts.Start, opts.End,
		),
		"NumRows", *opts.Res,
	)
	return nil
}
