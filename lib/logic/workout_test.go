package logic

import (
	"context"
	"testing"
	"time"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbtest "code.barbellmath.net/barbell-math/smoothbrain-test"
)

func TestRawWorkout(t *testing.T) {
	t.Run("failingNoWrites", workoutFailingNoWrites)
	t.Run("duplicateWorkout", workoutDuplicateWorkout)
	t.Run("addGetNoPhysicsData", workoutAddGetNoPhysicsData)
	t.Run("addGetTimeSeriesPhysicsData", workoutAddGetTimeSeriesPhysicsData)
	// t.Run("addGetVideoPhysicsData", workoutAddGetVideoPhysicsData)
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
	t.Run(
		"setFractionalSetsAndExercisesLen",
		workoutSetFractionalSetsAndExercisesLen(ctxt),
	)

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
		sbtest.ContainsError(t, types.InvalidWorkoutErr, err, "Unknown Email")
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
			Exercises: []types.RawExerciseData{
				{
					Sets: 2,
					Name: "badExercise",
				},
			},
		})
		sbtest.ContainsError(t, types.InvalidWorkoutErr, err, "Unknown Exercise")
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
			Exercises: []types.RawExerciseData{
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
			Exercises: []types.RawExerciseData{
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
			Exercises: []types.RawExerciseData{
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
			Exercises: []types.RawExerciseData{
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
			Exercises: []types.RawExerciseData{
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
			Exercises: []types.RawExerciseData{
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
			Exercises: []types.RawExerciseData{
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

func workoutSetFractionalSetsAndExercisesLen(
	ctxt context.Context,
) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateWorkouts(ctxt, types.RawWorkout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "email@email.com",
				Session:     1,
			},
			Exercises: []types.RawExerciseData{
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

func workoutDuplicateWorkout(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)

	err := CreateClients(ctxt, types.Client{
		FirstName: "FName", LastName: "LName", Email: "email@email.com",
	})

	workouts := [1]types.RawWorkout{
		types.RawWorkout{
			WorkoutID: types.WorkoutID{
				ClientEmail:   "email@email.com",
				Session:       1,
				DatePerformed: sbtest.MustParseTime(time.DateOnly, "2025-01-02"),
			},
			Exercises: []types.RawExerciseData{
				types.RawExerciseData{
					Name:   "Squat",
					Weight: 355,
					Sets:   1,
					Reps:   5,
					Effort: 8.5,
					BarPath: []types.BarPathVariant{
						types.BarPathTimeSeriesData(types.RawTimeSeriesData{
							TimeData:     []float64{0, 1, 2, 3, 4, 5, 6},
							PositionData: []float64{0, 1, 2, 3, 4, 5, 6},
						}),
					},
				},
			},
		},
	}

	err = CreateWorkouts(ctxt, workouts[:]...)
	sbtest.Nil(t, err)
	numExercises, err := ReadClientTotalNumExercises(ctxt, "email@email.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, numExercises)
	numRawWorkouts, err := ReadClientNumWorkouts(ctxt, "email@email.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, numRawWorkouts)
	numPhysEntries, err := ReadClientTotalNumPhysEntries(ctxt, "email@email.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, numPhysEntries)

	res, err := ReadWorkoutsByID(ctxt, workouts[0].WorkoutID)
	rawWorkoutEqSavedWorkout(t, workouts[:], res)

	err = CreateWorkouts(ctxt, workouts[:]...)
	sbtest.ContainsError(
		t, types.CouldNotAddWorkoutErr, err,
		"duplicate key value violates unique constraint",
	)
	numExercises, err = ReadClientTotalNumExercises(ctxt, "email@email.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, numExercises)
	numRawWorkouts, err = ReadClientNumWorkouts(ctxt, "email@email.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, numRawWorkouts)
	numPhysEntries, err = ReadClientTotalNumPhysEntries(ctxt, "email@email.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, numPhysEntries)
}

func workoutAddGetNoPhysicsData(t *testing.T) {
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
			Exercises: []types.RawExerciseData{
				types.RawExerciseData{
					Name:   "Squat",
					Weight: 355,
					Sets:   5,
					Reps:   5,
					Effort: 8.5,
				},
				types.RawExerciseData{
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
			Exercises: []types.RawExerciseData{
				types.RawExerciseData{
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
	rawWorkoutEqSavedWorkout(t, workouts[:], res)
}

func workoutAddGetTimeSeriesPhysicsData(t *testing.T) {
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
			Exercises: []types.RawExerciseData{
				types.RawExerciseData{
					Name:   "Squat",
					Weight: 355,
					Sets:   2,
					Reps:   5,
					Effort: 8.5,
					BarPath: []types.BarPathVariant{
						types.BarPathTimeSeriesData(types.RawTimeSeriesData{
							TimeData:     []float64{0, 1, 2, 3, 4, 5, 6},
							PositionData: []float64{0, 1, 2, 3, 4, 5, 6},
						}),
						types.BarPathTimeSeriesData(types.RawTimeSeriesData{
							TimeData:     []float64{0, 1, 2, 3, 4, 5, 6},
							PositionData: []float64{0, 1, 2, 3, 4, 5, 6},
						}),
					},
				},
				types.RawExerciseData{
					Name:   "Bench",
					Weight: 135,
					Sets:   1,
					Reps:   8,
					Effort: 5,
					BarPath: []types.BarPathVariant{
						types.BarPathTimeSeriesData(types.RawTimeSeriesData{
							TimeData:     []float64{0, 1, 2, 3, 4, 5, 6},
							PositionData: []float64{0, 1, 2, 3, 4, 5, 6},
						}),
					},
				},
			},
		},
		types.RawWorkout{
			WorkoutID: types.WorkoutID{
				ClientEmail:   "email@email.com",
				Session:       1,
				DatePerformed: sbtest.MustParseTime(time.DateOnly, "2025-01-03"),
			},
			Exercises: []types.RawExerciseData{
				types.RawExerciseData{
					Name:   "Deadlift",
					Weight: 405,
					Sets:   2,
					Reps:   6,
					Effort: 7,
					BarPath: []types.BarPathVariant{
						types.BarPathTimeSeriesData(types.RawTimeSeriesData{
							TimeData:     []float64{0, 1, 2, 3, 4, 5, 6},
							PositionData: []float64{0, 1, 2, 3, 4, 5, 6},
						}),
						types.BarPathTimeSeriesData(types.RawTimeSeriesData{
							TimeData:     []float64{0, 1, 2, 3, 4, 5, 6},
							PositionData: []float64{0, 1, 2, 3, 4, 5, 6},
						}),
					},
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
	rawWorkoutEqSavedWorkout(t, workouts[:], res)
}
