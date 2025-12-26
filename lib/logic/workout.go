package logic

import (
	"context"
	"time"

	"code.barbellmath.net/barbell-math/providentia/internal/dal"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
)

// Adds the supplied workouts to the database. The supplied workouts must have a
// valid workout ID. A valid workout ID must:
//
//   - Have a valid client email already present in the database
//   - Have a session number >0
//   - Have a valid date time stamp
//   - The workout ID must not already be present in the database
//
// Each exercise data in the exercise list must:
//
//   - Have a valid exercise name already present in the database
//   - Have a weight >=0
//   - Have sets >=0
//   - Have reps >=0
//   - Have effort in the range [0, 10]
//
// The context must have a [types.State] variable.
//
// Workouts will be uploaded in batches that respect the size set in the
// [State.BatchSize] variable.
//
// If any error occurs no changes will be made to the database.
func CreateWorkouts(
	ctxt context.Context,
	workouts ...types.Workout,
) (opErr error) {
	if len(workouts) == 0 {
		return
	}
	return runOp(ctxt, dal.CreateWorkouts, workouts)
}

// // Adds the clients supplied in the csv files to the database. Has the same
// // behavior as [CreateClients] other than getting the clients from csv files.
// // The csv files are expected to have column names on the first row and the
// // following columns must be present as identified by the column name on the
// // first row. More columns may be present, they will be ignored.
// //   - FirstName (string): the first name of the client
// //   - LastName (string): the last name of the client
// //   - Email (string): the email of the client
// //
// // For performance it is recommended to set the `ReuseRecord` variable to `true`
// // in the [sbcsv.Opts] struct. This will reduce the number of allocations made.
// //
// // The context must have a [types.State] variable.
// //
// // Clients will be uploaded in batches that respect the size set in the
// // [State.BatchSize] variable.
// //
// // If any error occurs no changes will be made to the database.
// func CreateClientsFromCSV(
// 	ctxt context.Context,
// 	opts *sbcsv.Opts,
// 	files ...string,
// ) (opErr error) {
// 	if len(files) == 0 {
// 		return
// 	}
// 	return runOp(ctxt, jobs.RunCSVLoaderJobs, jobs.CSVLoaderOpts[types.Client]{
// 		Opts:    opts,
// 		Files:   files,
// 		Creator: dal.CreateClients,
// 	})
// }

// Gets the total number of workouts in the database for a given client.
//
// The context must have a [types.State] variable.
//
// No changes will be made to the database.
func ReadNumWorkoutsForClient(
	ctxt context.Context,
	email string,
) (res int64, opErr error) {
	opErr = runOp(
		ctxt, dal.ReadNumWorkoutsForClient, dal.ReadNumWorkoutsForClientOpts{
			Email: email,
			Res:   &res,
		},
	)
	return
}

// Gets the workout data associated with the supplied ids if they exist. If they
// do not exist an error will be returned. The order of the returned workouts
// will match the order of the supplied workotu ids.
//
// The context must have a [types.State] variable.
//
// No changes will be made to the database.
func ReadWorkoutsById(
	ctxt context.Context,
	ids ...types.WorkoutId,
) (res []types.Workout, opErr error) {
	if len(ids) == 0 {
		return
	}
	opErr = runOp(
		ctxt, dal.ReadWorkoutsById, dal.ReadWorkoutsByIdOpts{
			Ids: ids,
			Res: &res,
		},
	)
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
func FindWorkoutsById(
	ctxt context.Context,
	ids ...types.WorkoutId,
) (res []types.Optional[types.Workout], opErr error) {
	if len(ids) == 0 {
		return
	}
	opErr = runOp(
		ctxt, dal.FindWorkoutsById, dal.FindWorkoutsByIdOpts{
			Ids: ids,
			Res: &res,
		},
	)
	return
}

// Gets the workouts for the supplied client in the supplied date range. If the
// supplied client does not exist no workouts will be returned and an error
// will be returned. If `start` is after `end` no workouts will be returned and
// an error will be returned.
//
// `start` is inclusive and `end` is exclusive.
//
// The context must have a [types.State] variable.
//
// No changes will be made to the database.
func FindWorkoutsInDateRange(
	ctxt context.Context,
	clientEmail string,
	start time.Time,
	end time.Time,
) (res []types.Workout, opErr error) {
	opErr = runOp(
		ctxt, dal.FindWorkoutsInDateRange, dal.FindWorkoutsInDateRangeOpts{
			Email: clientEmail,
			Start: start,
			End:   end,
			Res:   &res,
		},
	)
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
	ids ...types.WorkoutId,
) (opErr error) {
	if len(ids) == 0 {
		return
	}
	return runOp(ctxt, dal.DeleteWorkouts, ids)
}

// Deletes the workouts for the supplied client in the supplied date range
// returning the number of deleted workouts. If the supplied client does not
// exist no workouts will be deleted and an error will be returned. If `start`
// is after `end` no workouts will be deleted and an error will be returned. If
// no workouts exists between `start` and `end` an error will be returned.
//
// `start` is inclusive and `end` is exclusive.
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
	opErr = runOp(
		ctxt, dal.DeleteWorkoutsInDateRange, dal.DeleteWorkoutsInDateRangeOpts{
			Email: clientEmail,
			Start: start,
			End:   end,
			Res:   &res,
		},
	)
	return
}
