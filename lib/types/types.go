package types

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
