package logic

import (
	"context"
	"testing"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbtest "code.barbellmath.net/barbell-math/smoothbrain-test"
)

// TODO - eventually look into running tests in parallel - will need multiple dbs
func TestWorkout(t *testing.T) {
	t.Run("failingNoWrites", workoutFailingNoWrites)
	// TODO - add test for adding duplicated workout
	// t.Run("workoutCreateRead", workoutCreateRead)
	// t.Run("transactionRollback", clientTransactionRollback)
	// t.Run("addGet", clientAddGet)
	// t.Run("addUpdateGet", clientAddUpdateGet)
	// t.Run("addDeleteGet", clientAddDeleteGet)
}

func workoutFailingNoWrites(t *testing.T) {
	ctxt, cleanup := resetDB(context.Background())
	t.Cleanup(cleanup)

	err := CreateClients(ctxt, types.Client{
		FirstName: "FName", LastName: "LName", Email: "email@email.com",
	})
	sbtest.Nil(t, err)
	numClients, err := ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, numClients)

	t.Run("invalidSession", workoutInvalidSession(ctxt))
	t.Run("invalidClient", workoutInvalidClient(ctxt))
	t.Run("videoAndPhysDataDiffLen", workoutVideoAndPhysDataDiffLen(ctxt))
	t.Run("inconsistentPhysData", workoutInconsistentPhysData(ctxt))
	t.Run("unknownExercise", workoutUnknownExercise(ctxt))
	t.Run("setDirInsteadOfVideoFile", workoutSetDirInsteadOfVideoFile(ctxt))
	t.Run("setTimeAndPosDataDiffLen", workoutSetTimeAndPosDiffLen(ctxt))
	t.Run("setNotEnoughSamples", workoutSetNotEnoughSamples(ctxt))
	t.Run("setBackwardsTime", workoutSetBackwardsTime(ctxt))
	t.Run("setDiffTimeDelta", workoutSetDiffTimeDelta(ctxt))

	// numClients, err := ReadNumClients(ctxt)
	// sbtest.Nil(t, err)
	// sbtest.Eq(t, 0, numClients)
}

func workoutInvalidSession(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateWorkouts(ctxt, types.Workout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "email@email.com", Session: 0,
			},
		})
		sbtest.ContainsError(t, types.InvalidWorkoutErr, err)
		sbtest.ContainsError(t, types.InvalidSessionErr, err)
	}
}

func workoutInvalidClient(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateWorkouts(ctxt, types.Workout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "bad@email.com",
				Session:     1,
			},
		})
		sbtest.ContainsError(t, types.InvalidWorkoutErr, err)
		sbtest.ContainsError(t, types.CouldNotFindRequestedClientErr, err)
	}
}

func workoutVideoAndPhysDataDiffLen(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateWorkouts(ctxt, types.Workout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "email@email.com",
				Session:     1,
			},
			Exercises: []types.RawData{
				{
					VideoPaths:   []string{"a", "b"},
					TimeData:     [][]float64{[]float64{}},
					PositionData: [][]float64{[]float64{}},
				},
			},
		})
		sbtest.ContainsError(
			t, types.InvalidWorkoutErr, err,
			"the length of the supplied time data",
		)
		sbtest.ContainsError(t, types.MalformedWorkoutExerciseErr, err)

		err = CreateWorkouts(ctxt, types.Workout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "email@email.com",
				Session:     1,
			},
			Exercises: []types.RawData{
				{
					VideoPaths:   []string{"a"},
					TimeData:     [][]float64{[]float64{}, []float64{}},
					PositionData: [][]float64{[]float64{}, []float64{}},
				},
			},
		})
		sbtest.ContainsError(
			t, types.InvalidWorkoutErr, err,
			"the length of the supplied time data",
		)
		sbtest.ContainsError(t, types.MalformedWorkoutExerciseErr, err)
	}
}

func workoutInconsistentPhysData(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateWorkouts(ctxt, types.Workout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "email@email.com",
				Session:     1,
			},
			Exercises: []types.RawData{
				{
					VideoPaths:   []string{"a", "b"},
					TimeData:     [][]float64{[]float64{}, []float64{}},
					PositionData: [][]float64{[]float64{}},
				},
			},
		})
		sbtest.ContainsError(
			t, types.InvalidWorkoutErr, err,
			"the length of the supplied position data",
		)
		sbtest.ContainsError(t, types.MalformedWorkoutExerciseErr, err)

		err = CreateWorkouts(ctxt, types.Workout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "email@email.com",
				Session:     1,
			},
			Exercises: []types.RawData{
				{
					VideoPaths:   []string{"a"},
					TimeData:     [][]float64{[]float64{}},
					PositionData: [][]float64{[]float64{}, []float64{}},
				},
			},
		})
		sbtest.ContainsError(
			t, types.InvalidWorkoutErr, err,
			"the length of the supplied position data",
		)
		sbtest.ContainsError(t, types.MalformedWorkoutExerciseErr, err)
	}
}

func workoutUnknownExercise(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateWorkouts(ctxt, types.Workout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "email@email.com",
				Session:     1,
			},
			Exercises: []types.RawData{
				{
					Name:         "badExercise",
					VideoPaths:   []string{"a", "b"},
					TimeData:     [][]float64{[]float64{}, []float64{}},
					PositionData: [][]float64{[]float64{}, []float64{}},
				},
			},
		})
		sbtest.ContainsError(t, types.InvalidWorkoutErr, err, "Unknown exercise")
		sbtest.ContainsError(t, types.MalformedWorkoutExerciseErr, err)
	}
}

func workoutSetDirInsteadOfVideoFile(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateWorkouts(ctxt, types.Workout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "email@email.com",
				Session:     1,
			},
			Exercises: []types.RawData{
				{
					Name:         "Squat",
					VideoPaths:   []string{"./testData"},
					TimeData:     [][]float64{[]float64{}},
					PositionData: [][]float64{[]float64{}},
				},
			},
		})
		sbtest.ContainsError(
			t, types.InvalidWorkoutErr, err,
			"expected a video file, got dir",
		)
		sbtest.ContainsError(t, types.MalformedWorkoutExerciseErr, err)
	}
}

func workoutSetTimeAndPosDiffLen(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateWorkouts(ctxt, types.Workout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "email@email.com",
				Session:     1,
			},
			Exercises: []types.RawData{
				{
					Name:         "Squat",
					VideoPaths:   []string{""},
					TimeData:     [][]float64{[]float64{0, 1, 2, 3, 4, 5}},
					PositionData: [][]float64{[]float64{6, 7, 8, 9, 10}},
				},
			},
		})
		sbtest.ContainsError(
			t, types.InvalidWorkoutErr, err,
			"the length of the time data",
		)
		sbtest.ContainsError(t, types.MalformedWorkoutExerciseErr, err)
	}
}

func workoutSetNotEnoughSamples(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateWorkouts(ctxt, types.Workout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "email@email.com",
				Session:     1,
			},
			Exercises: []types.RawData{
				{
					Name:         "Squat",
					VideoPaths:   []string{""},
					TimeData:     [][]float64{[]float64{0, 1, 2}},
					PositionData: [][]float64{[]float64{3, 4, 5}},
				},
			},
		})
		sbtest.ContainsError(
			t, types.InvalidWorkoutErr, err,
			"the minimum number of samples",
		)
		sbtest.ContainsError(t, types.MalformedWorkoutExerciseErr, err)
	}
}

func workoutSetBackwardsTime(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateWorkouts(ctxt, types.Workout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "email@email.com",
				Session:     1,
			},
			Exercises: []types.RawData{
				{
					Name:         "Squat",
					VideoPaths:   []string{""},
					TimeData:     [][]float64{[]float64{1, 0, 2, 3, 4}},
					PositionData: [][]float64{[]float64{5, 6, 7, 8, 9}},
				},
			},
		})
		sbtest.ContainsError(
			t, types.InvalidWorkoutErr, err,
			"time samples must be increasing",
		)
		sbtest.ContainsError(t, types.MalformedWorkoutExerciseErr, err)
	}
}

func workoutSetDiffTimeDelta(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateWorkouts(ctxt, types.Workout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "email@email.com",
				Session:     1,
			},
			Exercises: []types.RawData{
				{
					Name:         "Squat",
					VideoPaths:   []string{""},
					TimeData:     [][]float64{[]float64{0, 1, 3, 3, 4}},
					PositionData: [][]float64{[]float64{5, 6, 7, 8, 9}},
				},
			},
		})
		sbtest.ContainsError(
			t, types.InvalidWorkoutErr, err,
			"time samples must all have the same delta",
		)
		sbtest.ContainsError(t, types.MalformedWorkoutExerciseErr, err)
	}
}

// func workoutInvalidVideoFile(ctxt context.Context) func(t *testing.T) {
// 	return func(t *testing.T) {
// 		err := CreateWorkouts(ctxt, types.Workout{
// 			WorkoutID: types.WorkoutID{
// 				ClientEmail: "email@email.com", Session: 1,
// 			},
// 			Exercises: []types.WorkoutExercise{
// 				{VideoPath: "./non/existant/path"},
// 			},
// 		})
// 		sbtest.ContainsError(t, types.InvalidVideoFileErr, err)
//
// 		err = CreateWorkouts(ctxt, types.Workout{
// 			WorkoutID: types.WorkoutID{
// 				ClientEmail: "email@email.com",
// 				Session:     1,
// 			},
// 			Exercises: []types.WorkoutExercise{
// 				{VideoPath: "../../bs/"},
// 			},
// 		})
// 		sbtest.ContainsError(t, types.InvalidVideoFileErr, err)
// 	}
// }

// func workoutCreateRead(t *testing.T) {
// 	ctxt, cleanup := resetDB(context.Background())
// 	t.Cleanup(cleanup)
//
// 	err := CreateClients(ctxt, types.Client{
// 		FirstName: "FName", LastName: "LName", Email: "email@email.com",
// 	})
//
// 	exercises := [8]types.WorkoutExercise{}
// 	for i := range len(exercises) {
// 		exercises[i] = types.WorkoutExercise{
// 			Name:      "Squat",
// 			Weight:    float64(i * 3),
// 			Sets:      float64(i*3 + 1),
// 			Reps:      int32(i*3 + 1),
// 			Effort:    4,
// 			VideoPath: "",
// 		}
// 	}
//
// 	err = CreateWorkouts(ctxt, types.Workout{
// 		WorkoutID: types.WorkoutID{
// 			ClientEmail:   "email@email.com",
// 			Session:       1,
// 			DatePerformed: time.Now(),
// 		},
// 		Exercises: exercises[:],
// 	})
// 	sbtest.Nil(t, err)
//
// 	// TODO - test num training logs?
// }
