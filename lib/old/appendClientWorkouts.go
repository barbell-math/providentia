package provlib

import (
	"context"
	"errors"
	"os"
	"time"

	dal "github.com/barbell-math/providentia/internal/db/dataAccessLayer"
	sberr "github.com/barbell-math/smoothbrain-errs"
	sblog "github.com/barbell-math/smoothbrain-logging"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type (
	TrainingLogData struct {
		Exercise      string
		DatePerformed time.Time
		Session       int32
		Weight        float64
		Sets          float64
		Reps          int32
		Effort        float64
		VideoPath     string
	}
)

var (
	InvalidSessionErr   = errors.New("Sesison must be >0")
	InvalidVideoFileErr = errors.New(
		"The supplied video file was not valid",
	)
	CouldNotFindRequestedExerciseErr = errors.New(
		"Could not find the requested exercise",
	)
	NewClientDataNotSortedErr = errors.New(
		"New client data must be sorted by date and session ascending",
	)
	CannotAppendDataBeforeExistingDataErr = errors.New(
		"Cannot append training log data if date performed is before the date of the clients last workout",
	)
	CouldNotAddTrainingDataErr = errors.New(
		"Could not add the requested training data",
	)
)

// Appends the supplied training log data for the supplied client to the
// database. All training log data must have a date *after* the supplied clients
// last workout. The training log data must also be sorted by date then session.
//
// The order of the supplied data will be preserved, meaning the order of the
// data given to this function will define the order of the exercises in the
// workout. The ordering of exercises in a workout is important to make any
// modeling as accurate as possible.
//
// The context must have a [State] variable.
//
// Training log data will be uploaded in batches that respect the size set in
// the [State.BatchSize] variable.
//
// If any error occurs no changes will be made to the database.
func AppendClientWorkouts(
	ctxt context.Context,
	clientEmail string,
	tlData ...TrainingLogData,
) (opErr error) {
	if len(tlData) == 0 {
		return
	}

	state, queries, cleanup, err := dbOpInfo(ctxt)
	if err != nil {
		opErr = err
		return
	}
	defer cleanup(&opErr)

	var clientID int64
	clientID, opErr = queries.GetClientIDFromEmail(ctxt, clientEmail)
	if opErr != nil {
		opErr = sberr.AppendError(CouldNotFindRequestedClientErr, opErr)
		return
	}

	if opErr = validateTrainingLogData(
		ctxt, queries, clientID, tlData,
	); opErr != nil {
		return
	}

	var curYear int
	var curMonth time.Month
	var curDay int
	var curSessionCntr int32
	var interWorkoutCntr int32
	exerciseMap := map[string]dal.GetExerciseIDsRow{}

	for start, end := range batchIndexes(tlData, int(state.Conf.BatchSize)) {
		// TODO - is there a way to pre-fetch all ids at once?
		// yes - SELECT * FROM Exercises WHERE Name IN (name1, name2);
		// but would make the error message more vague for potentially very
		// little performance gain
		ids := make([]dal.GetExerciseIDsRow, end-start)
		for i := 0; i < end-start; i++ {
			iterRawData := tlData[start+i]
			if iterIds, ok := exerciseMap[iterRawData.Exercise]; ok {
				ids[i] = iterIds
				continue
			}
			var iterIds dal.GetExerciseIDsRow
			iterIds, opErr = queries.GetExerciseIDs(ctxt, iterRawData.Exercise)
			if opErr != nil {
				opErr = sberr.AppendError(
					sberr.Wrap(
						CouldNotFindRequestedExerciseErr,
						"Missing exercise: %s", iterRawData.Exercise,
					),
					opErr,
				)
				return
			}
			exerciseMap[iterRawData.Exercise] = iterIds
			ids[i] = iterIds
		}

		// TODO - precalc things like video id (which will require some arbitrary
		// computaiton) and store the resulting ids in a tlData index->videoid map
		// TODO - can things like video data be computed in parallel here?? What if
		// the ctxt had a "allowed threads" param? Or threads per op??

		trainingLogs := make([]dal.BulkCreateTraingLogParams, end-start)

		for i := 0; i < end-start; i++ {
			iterRawData := tlData[start+i]
			rawYear, rawMonth, rawDay := iterRawData.DatePerformed.Date()
			if rawDay != curDay || rawMonth != curMonth || rawYear != curYear {
				interWorkoutCntr = 0
				curDay = rawDay
				curMonth = rawMonth
				curYear = rawYear
			} else if curSessionCntr != iterRawData.Session {
				interWorkoutCntr = 0
				curSessionCntr = iterRawData.Session
			}

			interWorkoutCntr++
			trainingLogs[i] = dal.BulkCreateTraingLogParams{
				ExerciseID:      ids[i].ExerciseID,
				ExerciseKindID:  ids[i].KindID,
				ExerciseFocusID: ids[i].FocusID,
				ClientID:        clientID,
				// TODO -Videoid ???

				DatePerformed: pgtype.Date{
					Time:             iterRawData.DatePerformed,
					InfinityModifier: pgtype.Finite,
					Valid:            true,
				},
				Weight: iterRawData.Weight,
				Sets:   iterRawData.Sets,
				Reps:   iterRawData.Reps,
				Effort: iterRawData.Effort,

				InterSessionCntr: iterRawData.Session,
				InterWorkoutCntr: interWorkoutCntr,
			}
		}

		var numRows int64
		numRows, opErr = queries.BulkCreateTraingLog(ctxt, trainingLogs)
		if opErr != nil {
			opErr = sberr.AppendError(CouldNotAddTrainingDataErr, opErr)
			return
		}
		state.Log.Log(
			ctxt, sblog.VLevel(2),
			"Appended training log entries",
			"ClientID", clientID,
			"NumRows", numRows,
		)
	}

	// TODO - run things to generate computed data (like model data) (what about
	// video data though?? - na, closer to database representation if it is above
	// here, the tl depends on the video data, it does not depend on generated
	// model data) (should I use a flag based system?? - na? save some space?)

	return
}

func validateTrainingLogData(
	ctxt context.Context,
	queries *dal.Queries,
	clientID int64,
	tlData []TrainingLogData,
) error {
	for i := 0; i < len(tlData); i++ {
		a := &tlData[i]
		if a.Session <= 0 {
			return sberr.Wrap(InvalidSessionErr, "Got: %d", a.Session)
		}

		if a.VideoPath != "" {
			fi, err := os.Stat(a.VideoPath)
			if err != nil {
				return sberr.AppendError(InvalidVideoFileErr, err)
			} else if fi.IsDir() {
				return sberr.Wrap(
					InvalidVideoFileErr,
					"Supplied path was a dir not a file",
				)
			}
		}

		if i > 0 {
			a := &tlData[i-1]
			b := &tlData[i]
			if a.DatePerformed.After(b.DatePerformed) {
				return NewClientDataNotSortedErr
			} else if a.Session > b.Session {
				return NewClientDataNotSortedErr
			}
		}
	}

	clientLastWorkoutDate, err := queries.ClientLastWorkoutDate(ctxt, clientID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		if clientLastWorkoutDate.Time.After(tlData[0].DatePerformed) {
			return sberr.Wrap(
				CannotAppendDataBeforeExistingDataErr,
				"Date of last workout: %v\nOldest requested date: %v",
				clientLastWorkoutDate, tlData[0].DatePerformed,
			)
		}
	}

	return nil
}

// TODO
func InsertClientWorkout(ctxt context.Context) error {
	panic("IMPLEMENT ME")
}
