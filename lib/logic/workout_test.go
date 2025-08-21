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
	// TODO - add test for adding duplicated workout - will fail when inserting
	// t.Run("duplicateWorkout", workoutDuplicateWorkout)
	// TODO - test that phys data is not saved in db when some fail
	// t.Run("transactionRollback", clientTransactionRollback)
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
	sbtest.Nil(t, err)
	sbtest.Eq(t, len(res), 2)
	sbtest.Eq(t, len(res[0].Exercises), 2)
	sbtest.Eq(t, len(res[1].Exercises), 1)

	sbtest.Eq(t, res[0].Exercises[0].Name, workouts[0].Exercises[0].Name)
	sbtest.Eq(t, res[0].Exercises[0].Weight, workouts[0].Exercises[0].Weight)
	sbtest.Eq(t, res[0].Exercises[0].Sets, workouts[0].Exercises[0].Sets)
	sbtest.Eq(t, res[0].Exercises[0].Reps, workouts[0].Exercises[0].Reps)
	sbtest.Eq(t, res[0].Exercises[0].Effort, workouts[0].Exercises[0].Effort)
	sbtest.Eq(
		t, res[0].Exercises[0].Volume,
		workouts[0].Exercises[0].Sets*float64(workouts[0].Exercises[0].Reps)*workouts[0].Exercises[0].Weight,
	)
	sbtest.Eq(
		t, res[0].Exercises[0].Exertion,
		workouts[0].Exercises[0].Sets*float64(workouts[0].Exercises[0].Reps)*workouts[0].Exercises[0].Effort,
	)
	sbtest.Eq(
		t, res[0].Exercises[0].TotalReps,
		workouts[0].Exercises[0].Sets*float64(workouts[0].Exercises[0].Reps),
	)

	sbtest.Eq(t, res[0].Exercises[1].Name, workouts[0].Exercises[1].Name)
	sbtest.Eq(t, res[0].Exercises[1].Weight, workouts[0].Exercises[1].Weight)
	sbtest.Eq(t, res[0].Exercises[1].Sets, workouts[0].Exercises[1].Sets)
	sbtest.Eq(t, res[0].Exercises[1].Reps, workouts[0].Exercises[1].Reps)
	sbtest.Eq(t, res[0].Exercises[1].Effort, workouts[0].Exercises[1].Effort)
	sbtest.Eq(
		t, res[0].Exercises[1].Volume,
		workouts[0].Exercises[1].Sets*float64(workouts[0].Exercises[1].Reps)*workouts[0].Exercises[1].Weight,
	)
	sbtest.Eq(
		t, res[0].Exercises[1].Exertion,
		workouts[0].Exercises[1].Sets*float64(workouts[0].Exercises[1].Reps)*workouts[0].Exercises[1].Effort,
	)
	sbtest.Eq(
		t, res[0].Exercises[1].TotalReps,
		workouts[0].Exercises[1].Sets*float64(workouts[0].Exercises[1].Reps),
	)

	sbtest.Eq(t, res[1].Exercises[0].Name, workouts[1].Exercises[0].Name)
	sbtest.Eq(t, res[1].Exercises[0].Weight, workouts[1].Exercises[0].Weight)
	sbtest.Eq(t, res[1].Exercises[0].Sets, workouts[1].Exercises[0].Sets)
	sbtest.Eq(t, res[1].Exercises[0].Reps, workouts[1].Exercises[0].Reps)
	sbtest.Eq(t, res[1].Exercises[0].Effort, workouts[1].Exercises[0].Effort)
	sbtest.Eq(
		t, res[1].Exercises[0].Volume,
		workouts[1].Exercises[0].Sets*float64(workouts[1].Exercises[0].Reps)*workouts[1].Exercises[0].Weight,
	)
	sbtest.Eq(
		t, res[1].Exercises[0].Exertion,
		workouts[1].Exercises[0].Sets*float64(workouts[1].Exercises[0].Reps)*workouts[1].Exercises[0].Effort,
	)
	sbtest.Eq(
		t, res[1].Exercises[0].TotalReps,
		workouts[1].Exercises[0].Sets*float64(workouts[1].Exercises[0].Reps),
	)

	sbtest.Eq(t, len(res[0].Exercises[0].Time), 0)
	sbtest.Eq(t, len(res[0].Exercises[0].Position), 0)
	sbtest.Eq(t, len(res[0].Exercises[1].Time), 0)
	sbtest.Eq(t, len(res[0].Exercises[1].Position), 0)
	sbtest.Eq(t, len(res[1].Exercises[0].Time), 0)
	sbtest.Eq(t, len(res[1].Exercises[0].Position), 0)
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
	// TODO - see if there is a better way to test if workouts are the same
	// This is some madness
	sbtest.Nil(t, err)
	sbtest.Eq(t, len(res), 2)
	sbtest.Eq(t, len(res[0].Exercises), 2)
	sbtest.Eq(t, len(res[1].Exercises), 1)

	sbtest.Eq(t, res[0].Exercises[0].Name, workouts[0].Exercises[0].Name)
	sbtest.Eq(t, res[0].Exercises[0].Weight, workouts[0].Exercises[0].Weight)
	sbtest.Eq(t, res[0].Exercises[0].Sets, workouts[0].Exercises[0].Sets)
	sbtest.Eq(t, res[0].Exercises[0].Reps, workouts[0].Exercises[0].Reps)
	sbtest.Eq(t, res[0].Exercises[0].Effort, workouts[0].Exercises[0].Effort)
	sbtest.Eq(
		t, res[0].Exercises[0].Volume,
		workouts[0].Exercises[0].Sets*float64(workouts[0].Exercises[0].Reps)*workouts[0].Exercises[0].Weight,
	)
	sbtest.Eq(
		t, res[0].Exercises[0].Exertion,
		workouts[0].Exercises[0].Sets*float64(workouts[0].Exercises[0].Reps)*workouts[0].Exercises[0].Effort,
	)
	sbtest.Eq(
		t, res[0].Exercises[0].TotalReps,
		workouts[0].Exercises[0].Sets*float64(workouts[0].Exercises[0].Reps),
	)

	sbtest.Eq(t, res[0].Exercises[1].Name, workouts[0].Exercises[1].Name)
	sbtest.Eq(t, res[0].Exercises[1].Weight, workouts[0].Exercises[1].Weight)
	sbtest.Eq(t, res[0].Exercises[1].Sets, workouts[0].Exercises[1].Sets)
	sbtest.Eq(t, res[0].Exercises[1].Reps, workouts[0].Exercises[1].Reps)
	sbtest.Eq(t, res[0].Exercises[1].Effort, workouts[0].Exercises[1].Effort)
	sbtest.Eq(
		t, res[0].Exercises[1].Volume,
		workouts[0].Exercises[1].Sets*float64(workouts[0].Exercises[1].Reps)*workouts[0].Exercises[1].Weight,
	)
	sbtest.Eq(
		t, res[0].Exercises[1].Exertion,
		workouts[0].Exercises[1].Sets*float64(workouts[0].Exercises[1].Reps)*workouts[0].Exercises[1].Effort,
	)
	sbtest.Eq(
		t, res[0].Exercises[1].TotalReps,
		workouts[0].Exercises[1].Sets*float64(workouts[0].Exercises[1].Reps),
	)

	sbtest.Eq(t, res[1].Exercises[0].Name, workouts[1].Exercises[0].Name)
	sbtest.Eq(t, res[1].Exercises[0].Weight, workouts[1].Exercises[0].Weight)
	sbtest.Eq(t, res[1].Exercises[0].Sets, workouts[1].Exercises[0].Sets)
	sbtest.Eq(t, res[1].Exercises[0].Reps, workouts[1].Exercises[0].Reps)
	sbtest.Eq(t, res[1].Exercises[0].Effort, workouts[1].Exercises[0].Effort)
	sbtest.Eq(
		t, res[1].Exercises[0].Volume,
		workouts[1].Exercises[0].Sets*float64(workouts[1].Exercises[0].Reps)*workouts[1].Exercises[0].Weight,
	)
	sbtest.Eq(
		t, res[1].Exercises[0].Exertion,
		workouts[1].Exercises[0].Sets*float64(workouts[1].Exercises[0].Reps)*workouts[1].Exercises[0].Effort,
	)
	sbtest.Eq(
		t, res[1].Exercises[0].TotalReps,
		workouts[1].Exercises[0].Sets*float64(workouts[1].Exercises[0].Reps),
	)

	sbtest.Eq(t, len(res[0].Exercises[0].Time), 2)
	sbtest.Eq(t, len(res[0].Exercises[0].Time[0]), 7)
	sbtest.Eq(t, len(res[0].Exercises[0].Time[1]), 7)
	sbtest.Eq(t, len(res[0].Exercises[0].Position), 2)
	sbtest.Eq(t, len(res[0].Exercises[0].Position[0]), 7)
	sbtest.Eq(t, len(res[0].Exercises[0].Position[1]), 7)
	sbtest.Eq(t, len(res[0].Exercises[1].Time), 1)
	sbtest.Eq(t, len(res[0].Exercises[1].Time[0]), 7)
	sbtest.Eq(t, len(res[0].Exercises[1].Position), 1)
	sbtest.Eq(t, len(res[0].Exercises[1].Position[0]), 7)
	sbtest.Eq(t, len(res[1].Exercises[0].Time), 2)
	sbtest.Eq(t, len(res[1].Exercises[0].Time[0]), 7)
	sbtest.Eq(t, len(res[1].Exercises[0].Time[1]), 7)
	sbtest.Eq(t, len(res[1].Exercises[0].Position), 2)
	sbtest.Eq(t, len(res[1].Exercises[0].Position[0]), 7)
	sbtest.Eq(t, len(res[1].Exercises[0].Position[1]), 7)
}
