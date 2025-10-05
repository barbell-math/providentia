package types

import (
	"errors"
	"time"
)

type (
	// The set of all available hyperparameters.
	Hyperparams interface {
		BarPathCalcHyperparams |
			BarPathTrackerHyperparams
	}

	// Hyperparameters used by the algorithm that calculates physics data from
	// the bars position over time.
	BarPathCalcHyperparams struct {
		Version         int32
		MinNumSamples   uint64
		TimeDeltaEps    Second
		ApproxErr       ApproximationError
		NearZeroFilter  float64
		SmootherWeight1 float64
		SmootherWeight2 float64
		SmootherWeight3 float64
		SmootherWeight4 float64
		SmootherWeight5 float64
	}

	// Hyperparameters used by the algorithm that gets the bars position over
	// time from a video.
	BarPathTrackerHyperparams struct {
		Version     int32
		MinLength   Second
		MinFileSize uint64
		MaxFileSize uint64
	}

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
		Session       uint16    // The session of the workout
		DatePerformed time.Time // The date the workout was done one
	}

	// Represents a workout that will be uploaded to the database
	RawWorkout struct {
		WorkoutID                   // The ID for the workout
		Exercises []RawExerciseData // All raw data about the workout provided by the user
	}
	// Contains all raw data that can be provided by the user
	RawExerciseData struct {
		Name    string   // The unique name of the exercise
		Weight  Kilogram // The weight the exercise was performed with
		Sets    float64  // The number sets that were performed
		Reps    int32    // The number of reps that were performed
		Effort  RPE      //The effort the exercise was performed at
		BarPath []BarPathVariant
	}

	// Represents a workout with all calculated data that is in the database
	// The `BasicData` and `PhysData` fields are associated arrays where each
	// index represents an exercise.
	Workout struct {
		WorkoutID
		Exercises []ExerciseData
	}
	// Represents basic data about a exercise within a workout with some basic
	// calculated fields such as volume
	ExerciseData struct {
		Name         string
		Weight       Kilogram
		Sets         float64
		Reps         int32
		Effort       RPE
		Volume       Kilogram
		Exertion     RPE
		TotalReps    float64
		Time         [][]Second
		Position     [][]Vec2[Meter, Meter]
		Velocity     [][]Vec2[MeterPerSec, MeterPerSec]
		Acceleration [][]Vec2[MeterPerSec2, MeterPerSec2]
		Jerk         [][]Vec2[MeterPerSec3, MeterPerSec3]
		Force        [][]Vec2[Newton, Newton]
		Impulse      [][]Vec2[NewtonSec, NewtonSec]
		Work         [][]Joule
		Power        [][]Watt
		RepSplits    [][]Split
		MinVel       [][]PointInTime[Second, MeterPerSec]
		MaxVel       [][]PointInTime[Second, MeterPerSec]
		MinAcc       [][]PointInTime[Second, MeterPerSec2]
		MaxAcc       [][]PointInTime[Second, MeterPerSec2]
		MinForce     [][]PointInTime[Second, Newton]
		MaxForce     [][]PointInTime[Second, Newton]
		MinImpulse   [][]PointInTime[Second, NewtonSec]
		MaxImpulse   [][]PointInTime[Second, NewtonSec]
		AvgWork      [][]Joule
		MinWork      [][]PointInTime[Second, Joule]
		MaxWork      [][]PointInTime[Second, Joule]
		AvgPower     [][]Watt
		MinPower     [][]PointInTime[Second, Watt]
		MaxPower     [][]PointInTime[Second, Watt]
	}
)

var (
	InvalidCtxtErr = errors.New("Invalid context")

	PhysicsJobQueueErr        = errors.New("Could not process physics job")
	TimeSeriesDecreaseErr     = errors.New("Time series data must not decrease")
	TimeSeriesNotMonotonicErr = errors.New("Time series must increase mononically")

	InvalidClientErr                 = errors.New("Invalid client")
	CouldNotAddClientsErr            = errors.New("Could not add the requested clients")
	MissingFirstNameErr              = errors.New("First name must not be empty")
	MissingLastNameErr               = errors.New("Last name must not be empty")
	MissingEmailErr                  = errors.New("Email must not be empty")
	CouldNotGetNumClientsErr         = errors.New("Could not get num clients")
	CouldNotFindRequestedClientErr   = errors.New("Could not find requested client")
	CouldNotUpdateRequestedClientErr = errors.New("Could not update requested client")
	CouldNotDeleteRequestedClientErr = errors.New("Could not delete requested client")

	InvalidExerciseErr                 = errors.New("Invalid exercise")
	CouldNotAddExercisesErr            = errors.New("Could not add the requested exercises")
	MissingExerciseNameErr             = errors.New("Exercise name must not be empty")
	CouldNotGetNumExercisesErr         = errors.New("Could not get num exercises")
	CouldNotFindRequestedExerciseErr   = errors.New("Could not find requested exercise")
	CouldNotUpdateRequestedExerciseErr = errors.New("Could not update requested exercise")
	CouldNotDeleteRequestedExerciseErr = errors.New("Could not delete requested exercise")

	InvalidWorkoutErr                 = errors.New("Invalid workout")
	InvalidSessionErr                 = errors.New("Invalid session num")
	MalformedWorkoutExerciseErr       = errors.New("Malformed exercise")
	InvalidBarPathsLenErr             = errors.New("The bar paths list must either be empty or the same length as the ceiling of the number of sets")
	VideoPathDirNotFileErr            = errors.New("Expected a video file, got dir")
	TimePositionDataMismatchErr       = errors.New("The length of the time data and position data must match")
	TimeDataLenErr                    = errors.New("The minimum number of samples was not provided")
	CouldNotAddWorkoutErr             = errors.New("Could not add the requested workouts")
	CouldNotGetNumWorkoutsErr         = errors.New("Could not get num workouts")
	CouldNotGetTotalNumExercisesErr   = errors.New("Could not get total num exercises")
	CouldNotGetTotalNumPhysEntriesErr = errors.New("Could not get total num phys entries")
	CouldNotFindRequestedWorkoutErr   = errors.New("Could not find requested workout")
	CouldNotDeleteRequestedWorkoutErr = errors.New("Could not delete requested workout")
	InvalidDataDirErr                 = errors.New("Invalid data dir")

	InvalidHyperparamsErr               = errors.New("Invalid hyperparameters")
	EncodingJsonHyperparamsErr          = errors.New("An error occurred encoding hyperparameters to json")
	DecodingJsonHyperparamsErr          = errors.New("An error occurred decoding hyperparameters from json")
	InvalidBarPathCalcErr               = errors.New("Invalid bar path calc conf")
	InvalidMinNumSamplesErr             = errors.New("Invalid min num samples")
	InvalidTimeDeltaEpsErr              = errors.New("Invalid time delta eps")
	InvalidNearZeroFilterErr            = errors.New("Invalid near zero filter")
	InvalidBarPathTrackerErr            = errors.New("Invalid bar path tracker conf")
	InvalidMinLengthErr                 = errors.New("Invalid min length")
	InvalidMinFileSizeErr               = errors.New("Invalid min file size")
	InvalidMaxFileSizeErr               = errors.New("Invalid max file size")
	CouldNotAddNumHyperparamsErr        = errors.New("Could not add requested hyperparams")
	CouldNotGetNumHyperparamsErr        = errors.New("Could not get num hyperparams")
	CouldNotFindRequestedHyperparamsErr = errors.New("Could not find requested hyperparams")
	CouldNotDeleteHyperparamsErr        = errors.New("Could not delete requester hyperparams")
)
