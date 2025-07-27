package types

import (
	"errors"
	"time"
)

type (
	// Represents a client from the database
	Client struct {
		FirstName string // The first name of the client
		LastName  string // The last name of the client
		Email     string // The clients email
	}

	// Represents an exercise from the database
	Exercise struct {
		Name    string
		KindID  ExerciseKind
		FocusID ExerciseFocus
	}

	// TODO - add comments about fields
	Workout struct {
		WorkoutID
		Exercises []WorkoutExercise
	}
	WorkoutID struct {
		ClientEmail   string
		Session       int32
		DatePerformed time.Time
	}
	WorkoutExercise struct {
		Name         string
		Weight       float64
		Sets         float64
		Reps         int32
		Effort       float64
		VideoPath    string
		PositionData []float64
	}
)

var (
	InvalidCtxtErr = errors.New("Invalid context")

	InvalidClientErr                 = errors.New("Invalid client")
	CouldNotAddClientsErr            = errors.New("Could not add the requested clients")
	CouldNotGetNumClientsErr         = errors.New("Could not get num clients")
	CouldNotFindRequestedClientErr   = errors.New("Could not find requested client")
	CouldNotUpdateRequestedClientErr = errors.New("Could not update requested client")
	CouldNotDeleteRequestedClientErr = errors.New("Could not delete requested client")

	InvalidExerciseErr                 = errors.New("Invalid exercise")
	CouldNotAddExercisesErr            = errors.New("Could not add the requested exercises")
	CouldNotGetNumExercisesErr         = errors.New("Could not get num exercises")
	CouldNotFindRequestedExerciseErr   = errors.New("Could not find requested exercise")
	CouldNotUpdateRequestedExerciseErr = errors.New("Could not update requested exercise")
	CouldNotDeleteRequestedExerciseErr = errors.New("Could not delete requested exercise")

	InvalidSessionErr     = errors.New("Invalid session num")
	InvalidVideoFileErr   = errors.New("The supplied video file was not valid")
	CouldNotAddWorkoutErr = errors.New("Could not add the requested workouts")
	// NewClientDataNotSortedErr             = errors.New("New client data must be sorted by date and session ascending")
	// CannotAppendDataBeforeExistingDataErr = errors.New("Cannot append training log data if date performed is before the date of the clients last workout")
	// CouldNotAddTrainingDataErr            = errors.New("Could not add the requested training data")
)
