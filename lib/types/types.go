package types

import "time"

// Basic types
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
		NoiseFilter     uint64
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
		FirstName string `db:"first_name"` // The first name of the client
		LastName  string `db:"last_name"`  // The last name of the client
		Email     string `db:"email"`      // The clients email
	}

	// Represents an exercise from the database
	Exercise struct {
		Name    string        `db:"name"`     // The exercise name
		KindId  ExerciseKind  `db:"kind_id"`  // The kind of exercise
		FocusId ExerciseFocus `db:"focus_id"` // The focus of the exercise
	}
)

// Workout types
type (
	// A struct that is used to represent the bar path when it has been
	// calculated by an external source. The [TimeData] and [PositionData]
	// slices must be the same length.
	RawTimeSeriesData struct {
		TimeData     []Second             // The time data for the set
		PositionData []Vec2[Meter, Meter] // The position data for the set
	}

	// A tagged union that either contains a [RawTimeSeriesData] struct or a
	// path to a video file.
	//
	// A zero initialized BarPathVariant will hold neither a video path or time
	// series data and can be used to represent having no data.
	//
	// Use [logic.BarPathVariant] or [logic.BarPathTimeSeriesData] to initialize.
	BarPathVariant struct {
		Flag       BarPathFlag
		VideoPath  string
		TimeSeries RawTimeSeriesData
	}

	// Higher order data that can be calculated from the basic user provided data.
	AbstractData struct {
		Volume    Kilogram // Weight*sets*reps
		Exertion  RPE      // Effort*sets*reps
		TotalReps float64  // Sets*reps
	}

	// Physics data calculated from either [RawTimeSeriesData] or a video.
	PhysicsData struct {
		BarPathCalcVersion    int32
		BarPathTrackerVersion int32
		VideoPath             string
		Time                  []Second
		Position              []Vec2[Meter, Meter]
		Velocity              []Vec2[MeterPerSec, MeterPerSec]
		Acceleration          []Vec2[MeterPerSec2, MeterPerSec2]
		Jerk                  []Vec2[MeterPerSec3, MeterPerSec3]
		Force                 []Vec2[Newton, Newton]
		Impulse               []Vec2[NewtonSec, NewtonSec]
		Work                  []Joule
		Power                 []Watt
		RepSplits             []Split
		MinVel                []PointInTime[Second, MeterPerSec]
		MaxVel                []PointInTime[Second, MeterPerSec]
		MinAcc                []PointInTime[Second, MeterPerSec2]
		MaxAcc                []PointInTime[Second, MeterPerSec2]
		MinForce              []PointInTime[Second, Newton]
		MaxForce              []PointInTime[Second, Newton]
		MinImpulse            []PointInTime[Second, NewtonSec]
		MaxImpulse            []PointInTime[Second, NewtonSec]
		AvgWork               []Joule
		MinWork               []PointInTime[Second, Joule]
		MaxWork               []PointInTime[Second, Joule]
		AvgPower              []Watt
		MinPower              []PointInTime[Second, Watt]
		MaxPower              []PointInTime[Second, Watt]
	}

	// Holds all data that can be collected when a lifter performs an exercise.
	ExerciseData struct {
		Name   string   // The unique name of the exercise
		Weight Kilogram // The weight the exercise was performed with
		Sets   float64  // The number sets that were performed
		Reps   int32    // The number of reps that were performed
		Effort RPE      // The effort the exercise was performed at

		AbstractData Optional[AbstractData]  // Will be calculated by the database for consistency
		PhysData     []Optional[PhysicsData] // Can be calculated with [logic.CalcPhysicsData]
	}

	// A unique identifier for a workout in the database
	WorkoutId struct {
		ClientEmail   string    // The clients unique email
		Session       uint16    // The session of the workout
		DatePerformed time.Time // The date the workout was done one
	}

	// Represents a fill workout performed by a lifter with all provided and
	// calculated data.
	Workout struct {
		WorkoutId
		Exercises []ExerciseData
	}
)
