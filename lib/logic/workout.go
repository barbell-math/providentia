package logic

import (
	"context"
	"time"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/internal/ops"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbcsv "code.barbellmath.net/barbell-math/smoothbrain-csv"
)

// Adds the supplied workouts to the database and calculates the physics
// information associated with any exercise that provides either a video path or
// time series data. The supplied workouts must have a valid workout ID and raw
// data. A valid workout ID must:
//
//   - Have a valid client email already present in the database
//   - Have a session number >0
//   - Have a valid date time stamp
//   - The workout ID must not already be present in the database
//
// Each raw data entry in the exercise list must:
//
//   - Have a valid exercise name already present in the database
//   - Have a weight >=0
//   - Have sets >=0
//   - Have reps >=0
//   - Have effort in the range [0, 10]
//   - Valid bar path data
//
// Bar path data can either be a path to a video file or time series data that
// represents the bars position over time. Valid time series bar path data must:
//
//   - Have time data and position data of the same length
//   - Have more than [state.PhysicsData.MinNumSamples] time samples
//   - The time data must be monotonically increasing, with a variance less than
//     [state.PhysicsData.TimeDeltaEps]
//
// Valid video bar path data must:
//
//   - Point to a valid path
//   - Point to a video longer than [state.BarPathTracker.MinLength]
//   - Point to a video file with a size greater than [state.BarPathTracker.MinSize]
//   - Point to a video file with a size less than [state.BarPathTracker.MaxSize]
//
// The context must have a [types.State] variable.
//
// Workouts will be uploaded in batches that respect the size set in the
// [State.BatchSize] variable.
//
// If any error occurs no changes will be made to the database.
func CreateWorkouts(
	ctxt context.Context,
	barPathCalcParams *types.BarPathCalcHyperparams,
	barTrackerCalcParams *types.BarPathTrackerHyperparams,
	workouts ...types.RawWorkout,
) (opErr error) {
	if len(workouts) == 0 {
		return
	}
	return runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) error {
			return ops.CreateWorkouts(
				ctxt, state, queries,
				barPathCalcParams, barTrackerCalcParams,
				workouts...,
			)
		},
	})
}

func EnsureWorkoutsExist(
	ctxt context.Context,
	barPathCalcParams *types.BarPathCalcHyperparams,
	barTrackerCalcParams *types.BarPathTrackerHyperparams,
	workouts ...types.RawWorkout,
) (opErr error) {
	if len(workouts) == 0 {
		return
	}
	return runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			return ops.CreateWorkouts(
				ctxt, state, queries,
				barPathCalcParams, barTrackerCalcParams,
				workouts...,
			)
		},
	})
}

// Adds the workouts supplied in the csv files to the database. Has the same
// behavior as [CreateWorkouts] other than getting the workouts from csv files.
//
// TODO - finish comment when done
// The csv files are expected to have column names on the first row and the
// following columns must be present as identified by the column name on the
// first row. More columns may be present, they will be ignored.
//   - FirstName (string): the first name of the client
//   - LastName (string): the last name of the client
//   - Email (string): the email of the client
//
// The `ReuseRecord` field on opts will be set to true before loading the csv
// file. All other options are left alone.
//
// The context must have a [types.State] variable.
//
// Clients will be uploaded in batches that respect the size set in the
// [State.BatchSize] variable.
//
// If any error occurs no changes will be made to the database.
func CreateWorkoutsFromCSV(
	ctxt context.Context,
	barPathCalcParams *types.BarPathCalcHyperparams,
	barTrackerCalcParams *types.BarPathTrackerHyperparams,
	opts sbcsv.Opts,
	files ...string,
) (opErr error) {
	if len(files) == 0 {
		return
	}
	return runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			return ops.CreateWorkoutsFromCSV(
				ctxt, state, queries,
				barPathCalcParams, barTrackerCalcParams,
				opts, files...,
			)
		},
	})
}

// Gets the total number of exercises across all workouts in the database for a
// given client.
//
// The context must have a [types.State] variable.
//
// No changes will be made to the database.
func ReadClientTotalNumTrainingLogEntries(
	ctxt context.Context,
	clientEmail string,
) (res int64, opErr error) {
	opErr = runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			res, err = ops.ReadClientTotalNumTrainingLogEntries(
				ctxt, state, queries, clientEmail,
			)
			return err
		},
	})
	return
}

// Gets the total number of physics entries across all workouts in the database
// for a given client. Each exercise with physics data will correspond to a
// single entry in the physics table.
//
// The context must have a [types.State] variable.
//
// No changes will be made to the database.
func ReadClientTotalNumPhysEntries(
	ctxt context.Context,
	clientEmail string,
) (res int64, opErr error) {
	opErr = runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			res, err = ops.ReadClientTotalNumPhysEntries(
				ctxt, state, queries, clientEmail,
			)
			return err
		},
	})
	return
}

// Gets the total number of workouts in the database for a given client.
//
// The context must have a [types.State] variable.
//
// No changes will be made to the database.
func ReadClientNumWorkouts(
	ctxt context.Context,
	clientEmail string,
) (res int64, opErr error) {
	opErr = runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			res, err = ops.ReadClientNumWorkouts(ctxt, state, queries, clientEmail)
			return err
		},
	})
	return
}

// Gets the workout data associated with the supplied ids if they exist. If they
// do not exist an error will be returned. The order of the returned workouts
// will match the order of the supplied workotu ids.
//
// The context must have a [types.State] variable.
//
// No changes will be made to the database.
func ReadWorkoutsByID(
	ctxt context.Context,
	ids ...types.WorkoutID,
) (res []types.Workout, opErr error) {
	if len(ids) == 0 {
		return
	}
	opErr = runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			res, err = ops.ReadWorkoutsByID(ctxt, state, queries, ids...)
			return
		},
	})
	return
}

// Gets the workout data associated with the supplied workout ids if they exist.
// If a workout exists it will be put in the returned slice and the found flag
// will be set to true. If a workout does not exist the value in the slice will
// be a zero initialized workout and the found flag will be set to false. No
// error will be returned if a workout does not exist. The order of the returned
// workouts will match the order of the supplied workout id's.
//
// The context must have a [types.State] variable.
//
// No changes will be made to the database.
func FindWorkoutsByID(
	ctxt context.Context,
	ids ...types.WorkoutID,
) (res []types.Found[types.Workout], opErr error) {
	if len(ids) == 0 {
		return
	}
	opErr = runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			res, err = ops.FindWorkoutsByID(ctxt, state, queries, ids...)
			return
		},
	})
	return
}

// Gets the workouts for the supplied client in the supplied date range. If the
// supplied client does not exist no workouts will be returned and an error
// will be returned. If `start` is after `end` no workouts will be returned and
// an error will be returned. If no workouts exists between `start` and `end` an
// error will be returned.
//
// Both `start` and `end` are inclusive.
//
// The context must have a [types.State] variable.
//
// No changes will be made to the database.
func ReadWorkoutsInDateRange(
	ctxt context.Context,
	clientEmail string,
	start time.Time,
	end time.Time,
) (res []types.Workout, opErr error) {
	opErr = runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			res, err = ops.ReadWorkoutsInDateRange(
				ctxt, state, queries, clientEmail, start, end,
			)
			return err
		},
	})
	return
}

// Deletes the workout data associated with the supplied ids if they exist. If
// they do not exist an error will be returned.
//
// The context must have a [types.State] variable.
//
// If any error occurs no changes will be made to the database.
func DeleteWorkouts(
	ctxt context.Context,
	ids ...types.WorkoutID,
) (opErr error) {
	if len(ids) == 0 {
		return
	}
	opErr = runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			err = ops.DeleteWorkouts(ctxt, state, queries, ids...)
			return err
		},
	})
	return
}

// Deletes the workouts for the supplied client in the supplied date range
// returning the number of deleted workouts. If the supplied client does not
// exist no workouts will be deleted and an error will be returned. If `start`
// is after `end` no workouts will be deleted and an error will be returned. If
// no workouts exists between `start` and `end` an error will be returned.
//
// Both `start` and `end` are inclusive.
//
// The context must have a [types.State] variable.
//
// If any error occurs no changes will be made to the database.
func DeleteWorkoutsInDateRange(
	ctxt context.Context,
	clientEmail string,
	start time.Time,
	end time.Time,
) (res int64, opErr error) {
	opErr = runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) (err error) {
			res, err = ops.DeleteWorkoutsInDateRange(
				ctxt, state, queries, clientEmail, start, end,
			)
			return err
		},
	})
	return
}
