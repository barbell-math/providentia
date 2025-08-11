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
		Name    string        // The exercise name
		KindID  ExerciseKind  // The kind of exercise
		FocusID ExerciseFocus // The focus of the exercise
	}

	// A unique identifier for a workout in the database
	WorkoutID struct {
		ClientEmail   string    // The clients unique email
		Session       uint32    // The session of the workout
		DatePerformed time.Time // The date the workout was done one
	}

	// Represents a workout that will be uploaded to the database
	Workout struct {
		WorkoutID           // The ID for the workout
		Exercises []RawData // All raw data about the workout provided by the user
	}
	// Contains all raw data that can be provided by the user
	RawData struct {
		Name   string  // The unique name of the exercise
		Weight float64 // The weight the exercise was performed with
		Sets   float64 // The number sets that were performed
		Reps   int32   // The number of reps that were performed
		Effort float64 //The effort the exercise was performed at
		// The VideoPath and TimeData,PositionData fields are mutually exclusive
		// for each set. For each set, either supply a video path *or* time data
		// and position data. The time and position data will be gathered from
		// the video if it is supplied.
		VideoPaths   []string    // The video path for each set
		TimeData     [][]float64 // The time data for each set
		PositionData [][]float64 // The position data for each set
	}

	CalculatedData struct {
		WorkoutID
		BasicData []BasicData
		PhysData  []PhysicsData
	}
	BasicData struct {
		Name      string
		Weight    float64
		Sets      float64
		Reps      int32
		Effort    float64
		Volume    float64
		Exertion  float64
		TotalReps float64
	}
	PhysicsData struct {
		Time         [][]float64
		Position     [][]float64
		Velocity     [][]float64
		Acceleration [][]float64
		Jerk         [][]float64
		Force        [][]float64
		Impulse      [][]float64
		Work         [][]float64
	}
)

var (
	InvalidCtxtErr          = errors.New("Invalid context")
	InvalidMinNumSamplesErr = errors.New("Invalid min num samples")

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

	InvalidWorkoutErr           = errors.New("Invalid workout")
	InvalidSessionErr           = errors.New("Invalid session num")
	MalformedWorkoutExerciseErr = errors.New("Malformed exercise")
	CouldNotAddWorkoutErr       = errors.New("Could not add the requested workouts")
	// NewClientDataNotSortedErr             = errors.New("New client data must be sorted by date and session ascending")
	// CannotAppendDataBeforeExistingDataErr = errors.New("Cannot append training log data if date performed is before the date of the clients last workout")
	// CouldNotAddTrainingDataErr            = errors.New("Could not add the requested training data")
)
