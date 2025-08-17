package logic

import (
	"context"
	"testing"
	"time"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbtest "code.barbellmath.net/barbell-math/smoothbrain-test"
)

// TODO - eventually look into running tests in parallel - will need multiple dbs
func TestRawWorkout(t *testing.T) {
	t.Run("failingNoWrites", workoutFailingNoWrites)
	t.Run("workoutCreateReadNoPhysicsData", workoutCreateReadNoPhysicsData)
	// TODO - add test for adding duplicated workout - will fail when inserting
	// TODO - test that phys data is not saved in db when some fail
	// t.Run("transactionRollback", clientTransactionRollback)
	// t.Run("addGet", clientAddGet)
	// t.Run("addUpdateGet", clientAddUpdateGet)
	// t.Run("addDeleteGet", clientAddDeleteGet)
}

func workoutFailingNoWrites(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
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
	t.Run("unknownExercise", workoutUnknownExercise(ctxt))
	t.Run("setTimeAndPosDataDiffLen", workoutSetTimeAndPosDiffLen(ctxt))
	t.Run("setNotEnoughSamples", workoutSetNotEnoughSamples(ctxt))
	t.Run("setBackwardsTime", workoutSetBackwardsTime(ctxt))
	t.Run("setDiffTimeDelta", workoutSetDiffTimeDelta(ctxt))
	t.Run("setDirInsteadOfVideoFile", workoutSetDirInsteadOfVideoFile(ctxt))
	t.Run("setInvalidVideoFile", workoutSetInvalidVideoFile(ctxt))
	t.Run("setFractionalSetsAndPhysDataLen", workoutSetFractionalSetsAndPhysDataLen(ctxt))

	numClients, err = ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, numClients)
	numExercises, err := ReadClientTotalNumExercises(ctxt, "email@email.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, 0, numExercises)
	numRawWorkouts, err := ReadClientNumWorkouts(ctxt, "email@email.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, 0, numRawWorkouts)
}

func workoutInvalidSession(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateWorkouts(ctxt, types.RawWorkout{
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
		err := CreateWorkouts(ctxt, types.RawWorkout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "bad@email.com",
				Session:     1,
			},
		})
		sbtest.ContainsError(t, types.InvalidWorkoutErr, err)
		sbtest.ContainsError(t, types.CouldNotFindRequestedClientErr, err)
	}
}

func workoutUnknownExercise(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateWorkouts(ctxt, types.RawWorkout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "email@email.com",
				Session:     1,
			},
			Exercises: []types.RawData{
				{
					Sets: 2,
					Name: "badExercise",
				},
			},
		})
		sbtest.ContainsError(t, types.InvalidWorkoutErr, err, "Unknown exercise")
		sbtest.ContainsError(t, types.MalformedWorkoutExerciseErr, err)
	}
}

func workoutSetTimeAndPosDiffLen(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateWorkouts(ctxt, types.RawWorkout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "email@email.com",
				Session:     1,
			},
			Exercises: []types.RawData{
				{
					Sets: 1,
					Name: "Squat",
					BarPath: []types.BarPathVariant{
						types.BarPathTimeSeriesData(
							types.RawTimeSeriesData{
								TimeData:     []float64{0, 1, 2},
								PositionData: []float64{0, 1, 2, 3},
							},
						),
					},
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
		err := CreateWorkouts(ctxt, types.RawWorkout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "email@email.com",
				Session:     1,
			},
			Exercises: []types.RawData{
				{
					Sets: 1,
					Name: "Squat",
					BarPath: []types.BarPathVariant{
						types.BarPathTimeSeriesData(
							types.RawTimeSeriesData{
								TimeData:     []float64{0, 1, 2, 3},
								PositionData: []float64{0, 1, 2, 3},
							},
						),
					},
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
		err := CreateWorkouts(ctxt, types.RawWorkout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "email@email.com",
				Session:     1,
			},
			Exercises: []types.RawData{
				{
					Sets: 1,
					Name: "Squat",
					BarPath: []types.BarPathVariant{
						types.BarPathTimeSeriesData(
							types.RawTimeSeriesData{
								TimeData:     []float64{1, 0, 2, 3, 4},
								PositionData: []float64{0, 1, 2, 3, 4},
							},
						),
					},
				},
			},
		})
		sbtest.ContainsError(
			t, types.CouldNotAddWorkoutErr, err,
			"Time samples must be increasing, got a delta of",
		)
		sbtest.ContainsError(t, types.PhysicsJobQueueErr, err)
		sbtest.ContainsError(t, types.TimeSeriesDecreaseErr, err)
	}
}

func workoutSetDiffTimeDelta(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateWorkouts(ctxt, types.RawWorkout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "email@email.com",
				Session:     1,
			},
			Exercises: []types.RawData{
				{
					Sets: 1,
					Name: "Squat",
					BarPath: []types.BarPathVariant{
						types.BarPathTimeSeriesData(
							types.RawTimeSeriesData{
								TimeData:     []float64{0, 1, 3, 3, 4},
								PositionData: []float64{0, 1, 2, 3, 4},
							},
						),
					},
				},
			},
		})
		sbtest.ContainsError(
			t, types.CouldNotAddWorkoutErr, err,
			"Time samples must all have the same delta",
		)
		sbtest.ContainsError(t, types.PhysicsJobQueueErr, err)
		sbtest.ContainsError(t, types.TimeSeriesNotMonotonicErr, err)
	}
}

func workoutSetDirInsteadOfVideoFile(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateWorkouts(ctxt, types.RawWorkout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "email@email.com",
				Session:     1,
			},
			Exercises: []types.RawData{
				{
					Sets:    1,
					Name:    "Squat",
					BarPath: []types.BarPathVariant{types.BarPathVideo(".")},
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

func workoutSetInvalidVideoFile(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateWorkouts(ctxt, types.RawWorkout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "email@email.com",
				Session:     1,
			},
			Exercises: []types.RawData{
				{
					Sets: 1,
					Name: "Squat",
					BarPath: []types.BarPathVariant{
						types.BarPathVideo("./non-existant-dir"),
					},
				},
			},
		})
		sbtest.ContainsError(
			t, types.InvalidWorkoutErr, err,
			"no such file or directory",
		)
		sbtest.ContainsError(t, types.MalformedWorkoutExerciseErr, err)
	}
}

func workoutSetNotEnoughBarPathEntries(
	ctxt context.Context,
) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateWorkouts(ctxt, types.RawWorkout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "email@email.com",
				Session:     1,
			},
			Exercises: []types.RawData{
				{
					Name: "Squat",
					Sets: 3,
					BarPath: []types.BarPathVariant{
						types.BarPathTimeSeriesData(
							types.RawTimeSeriesData{
								TimeData:     []float64{0, 1, 2, 3, 4},
								PositionData: []float64{0, 1, 2, 3, 4},
							},
						),
						types.BarPathTimeSeriesData(
							types.RawTimeSeriesData{
								TimeData:     []float64{0, 1, 2, 3, 4},
								PositionData: []float64{0, 1, 2, 3, 4},
							},
						),
					},
				},
			},
		})
		sbtest.ContainsError(
			t, types.InvalidWorkoutErr, err,
			"the bar paths list must either be empty or the same length as the ceiling of the number of sets",
		)
		sbtest.ContainsError(t, types.MalformedWorkoutExerciseErr, err)
	}
}

func workoutSetFractionalSetsAndPhysDataLen(
	ctxt context.Context,
) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateWorkouts(ctxt, types.RawWorkout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "email@email.com",
				Session:     1,
			},
			Exercises: []types.RawData{
				{
					Name: "Squat",
					Sets: 2.5,
					BarPath: []types.BarPathVariant{
						types.BarPathTimeSeriesData(
							types.RawTimeSeriesData{
								TimeData:     []float64{0, 1, 2, 3, 4},
								PositionData: []float64{0, 1, 2, 3, 4},
							},
						),
						types.BarPathTimeSeriesData(
							types.RawTimeSeriesData{
								TimeData:     []float64{0, 1, 2, 3, 4},
								PositionData: []float64{0, 1, 2, 3, 4},
							},
						),
					},
				},
			},
		})
		sbtest.ContainsError(
			t, types.InvalidWorkoutErr, err,
			"the bar paths list must either be empty or the same length as the ceiling of the number of sets",
		)
		sbtest.ContainsError(t, types.MalformedWorkoutExerciseErr, err)
	}
}

// func workoutCreateRead(t *testing.T) {
// 	ctxt, cleanup := resetApp(context.Background())
// 	t.Cleanup(cleanup)
//
// 	err := CreateClients(ctxt, types.Client{
// 		FirstName: "FName", LastName: "LName", Email: "email@email.com",
// 	})
//
// 	exercises := [8]types.RawWorkoutExercise{}
// 	for i := range len(exercises) {
// 		exercises[i] = types.RawWorkoutExercise{
// 			Name:      "Squat",
// 			Weight:    float64(i * 3),
// 			Sets:      float64(i*3 + 1),
// 			Reps:      int32(i*3 + 1),
// 			Effort:    4,
// 			VideoPath: "",
// 		}
// 	}
//
// 	err = CreateRawWorkouts(ctxt, types.RawWorkout{
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

func workoutCreateReadNoPhysicsData(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)

	err := CreateClients(ctxt, types.Client{
		FirstName: "FName", LastName: "LName", Email: "email@email.com",
	})

	workouts := [2]types.RawWorkout{
		types.RawWorkout{
			WorkoutID: types.WorkoutID{
				ClientEmail:   "email@email.com",
				Session:       1,
				DatePerformed: sbtest.MustParseTime(time.DateOnly, "2025-01-02"),
			},
			Exercises: []types.RawData{
				types.RawData{
					Name:   "Squat",
					Weight: 355,
					Sets:   5,
					Reps:   5,
					Effort: 8.5,
				},
				types.RawData{
					Name:   "Bench",
					Weight: 135,
					Sets:   3,
					Reps:   8,
					Effort: 5,
				},
			},
		},
		types.RawWorkout{
			WorkoutID: types.WorkoutID{
				ClientEmail:   "email@email.com",
				Session:       1,
				DatePerformed: sbtest.MustParseTime(time.DateOnly, "2025-01-03"),
			},
			Exercises: []types.RawData{
				types.RawData{
					Name:   "Deadlift",
					Weight: 405,
					Sets:   6,
					Reps:   6,
					Effort: 7,
				},
			},
		},
	}

	err = CreateWorkouts(ctxt, workouts[:]...)
	sbtest.Nil(t, err)
	numExercises, err := ReadClientTotalNumExercises(ctxt, "email@email.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, 3, numExercises)
	numRawWorkouts, err := ReadClientNumWorkouts(ctxt, "email@email.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, 2, numRawWorkouts)

	res, err := ReadWorkoutsByID(ctxt, workouts[0].WorkoutID, workouts[1].WorkoutID)
	sbtest.Nil(t, err)
	sbtest.Eq(t, len(res), 2)
	sbtest.Eq(t, len(res[0].BasicData), 2)
	sbtest.Eq(t, len(res[0].PhysData), 2)
	sbtest.Eq(t, len(res[1].BasicData), 1)
	sbtest.Eq(t, len(res[1].PhysData), 1)

	sbtest.Eq(t, res[0].BasicData[0], types.BasicData{
		Name:      workouts[0].Exercises[0].Name,
		Weight:    workouts[0].Exercises[0].Weight,
		Sets:      workouts[0].Exercises[0].Sets,
		Reps:      workouts[0].Exercises[0].Reps,
		Effort:    workouts[0].Exercises[0].Effort,
		Volume:    workouts[0].Exercises[0].Sets * float64(workouts[0].Exercises[0].Reps) * workouts[0].Exercises[0].Weight,
		Exertion:  workouts[0].Exercises[0].Sets * float64(workouts[0].Exercises[0].Reps) * workouts[0].Exercises[0].Effort,
		TotalReps: workouts[0].Exercises[0].Sets * float64(workouts[0].Exercises[0].Reps),
	})
	sbtest.Eq(t, res[0].BasicData[1], types.BasicData{
		Name:      workouts[0].Exercises[1].Name,
		Weight:    workouts[0].Exercises[1].Weight,
		Sets:      workouts[0].Exercises[1].Sets,
		Reps:      workouts[0].Exercises[1].Reps,
		Effort:    workouts[0].Exercises[1].Effort,
		Volume:    workouts[0].Exercises[1].Sets * float64(workouts[0].Exercises[1].Reps) * workouts[0].Exercises[1].Weight,
		Exertion:  workouts[0].Exercises[1].Sets * float64(workouts[0].Exercises[1].Reps) * workouts[0].Exercises[1].Effort,
		TotalReps: workouts[0].Exercises[1].Sets * float64(workouts[0].Exercises[1].Reps),
	})
	sbtest.Eq(t, res[1].BasicData[0], types.BasicData{
		Name:      workouts[1].Exercises[0].Name,
		Weight:    workouts[1].Exercises[0].Weight,
		Sets:      workouts[1].Exercises[0].Sets,
		Reps:      workouts[1].Exercises[0].Reps,
		Effort:    workouts[1].Exercises[0].Effort,
		Volume:    workouts[1].Exercises[0].Sets * float64(workouts[1].Exercises[0].Reps) * workouts[1].Exercises[0].Weight,
		Exertion:  workouts[1].Exercises[0].Sets * float64(workouts[1].Exercises[0].Reps) * workouts[1].Exercises[0].Effort,
		TotalReps: workouts[1].Exercises[0].Sets * float64(workouts[1].Exercises[0].Reps),
	})
}
