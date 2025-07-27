package logic

import (
	"context"
	"testing"
	"time"

	"github.com/barbell-math/providentia/lib/types"
	sbtest "github.com/barbell-math/smoothbrain-test"
)

// TODO - eventually look into running tests in parallel - will need multiple dbs
func TestWorkout(t *testing.T) {
	t.Run("failingNoWrites", workoutFailingNoWrites)
	t.Run("workoutCreateRead", workoutCreateRead)
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

	t.Run("workoutInvalidClient", workoutInvalidClient(ctxt))
	t.Run("workoutInvalidSession", workoutInvalidSession(ctxt))
	t.Run("workoutInvalidVideoFile", workoutInvalidVideoFile(ctxt))

	// numClients, err := ReadNumClients(ctxt)
	// sbtest.Nil(t, err)
	// sbtest.Eq(t, 0, numClients)
}

func workoutInvalidClient(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateWorkouts(ctxt, types.Workout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "bad@email.com",
				Session:     1,
			},
		})
		sbtest.ContainsError(t, types.CouldNotFindRequestedClientErr, err)
	}
}

func workoutInvalidSession(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateWorkouts(ctxt, types.Workout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "email@email.com", Session: 0,
			},
		})
		sbtest.ContainsError(t, types.InvalidSessionErr, err)
	}
}

func workoutInvalidVideoFile(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateWorkouts(ctxt, types.Workout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "email@email.com", Session: 1,
			},
			Exercises: []types.WorkoutExercise{
				{VideoPath: "./non/existant/path"},
			},
		})
		sbtest.ContainsError(t, types.InvalidVideoFileErr, err)

		err = CreateWorkouts(ctxt, types.Workout{
			WorkoutID: types.WorkoutID{
				ClientEmail: "email@email.com",
				Session:     1,
			},
			Exercises: []types.WorkoutExercise{
				{VideoPath: "../../bs/"},
			},
		})
		sbtest.ContainsError(t, types.InvalidVideoFileErr, err)
	}
}

func workoutCreateRead(t *testing.T) {
	ctxt, cleanup := resetDB(context.Background())
	t.Cleanup(cleanup)

	err := CreateClients(ctxt, types.Client{
		FirstName: "FName", LastName: "LName", Email: "email@email.com",
	})

	exercises := [8]types.WorkoutExercise{}
	for i := range len(exercises) {
		exercises[i] = types.WorkoutExercise{
			Name:      "Squat",
			Weight:    float64(i * 3),
			Sets:      float64(i*3 + 1),
			Reps:      int32(i*3 + 1),
			Effort:    4,
			VideoPath: "",
		}
	}

	err = CreateWorkouts(ctxt, types.Workout{
		WorkoutID: types.WorkoutID{
			ClientEmail:   "email@email.com",
			Session:       1,
			DatePerformed: time.Now(),
		},
		Exercises: exercises[:],
	})
	sbtest.Nil(t, err)

	// TODO - test num training logs?
}
