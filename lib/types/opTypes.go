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
	RawWorkout struct {
		WorkoutID           // The ID for the workout
		Exercises []RawData // All raw data about the workout provided by the user
	}
	// Contains all raw data that can be provided by the user
	RawData struct {
		Name    string  // The unique name of the exercise
		Weight  float64 // The weight the exercise was performed with
		Sets    float64 // The number sets that were performed
		Reps    int32   // The number of reps that were performed
		Effort  float64 //The effort the exercise was performed at
		BarPath []BarPathVariant
	}

	// Represents a workout with all calculated data that is in the database
	// The `BasicData` and `PhysData` fields are associated arrays where each
	// index represents an exercise.
	Workout struct {
		WorkoutID
		BasicData []BasicData
		PhysData  []PhysicsData
	}
	// Represents basic data about a exercise within a workout with some basic
	// calculated fields such as volume
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
	// Represents the calculated physics data for each set and rep of a given
	// exercsise
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
	InvalidCtxtErr = errors.New("Invalid context")

	PhysicsJobQueueErr        = errors.New("Could not process physics job")
	TimeSeriesDecreaseErr     = errors.New("Time series data must not decrease")
	TimeSeriesNotMonotonicErr = errors.New("Time series must increase mononically")

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

	InvalidWorkoutErr               = errors.New("Invalid workout")
	InvalidSessionErr               = errors.New("Invalid session num")
	MalformedWorkoutExerciseErr     = errors.New("Malformed exercise")
	CouldNotAddWorkoutErr           = errors.New("Could not add the requested workouts")
	CouldNotGetNumWorkoutsErr       = errors.New("Could not get num workouts")
	CouldNotGetTotalNumExercisesErr = errors.New("Could not get total num exercises")
	CouldNotFindRequestedWorkoutErr = errors.New("Could not find requested workout")
)
