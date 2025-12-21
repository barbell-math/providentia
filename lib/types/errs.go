package types

import "errors"

// [State] errors
var (
	InvalidCtxtErr = errors.New("Invalid context")

	InvalidGlobalConfErr = errors.New("Invalid global conf")
	InvalidBatchSizeErr  = errors.New("Invalid batch size")

	InvalidLoggerErr            = errors.New("Invalid logger")
	InvalidDBErr                = errors.New("Invalid database connection pool")
	InvalidPhysicsJobQueueErr   = errors.New("Invalid physics job queue")
	InvalidVideoJobQueue        = errors.New("Invalid video job queue")
	InvalidCSVLoaderJobQueueErr = errors.New("Invalid csv loader job queue")
	InvalidGPJobQueueErr        = errors.New("Invalid general purpose loader job queue")
)

// [ExerciseFocus] errors
var (
	CouldNotCreateAllExerciseFocusEntriesErr = errors.New("Could not create all exercise focus entries")
)

// [ExerciseKind] errors
var (
	CouldNotCreateAllExerciseKindEntriesErr = errors.New("Could not create all exercise kind entries")
)

// [Model] errors
var (
	CouldNotCreateAllModelsErr = errors.New("Could not create all models")
)

// [Client] errors
var (
	CouldNotCreateAllClientsErr = errors.New("Could not create all clients")
	CouldNotReadAllClientsErr   = errors.New("Could not read all clients")
	CouldNotUpdateAllClientsErr = errors.New("Could not update all clients")
	CouldNotDeleteAllClientsErr = errors.New("Could not delete all clients")
)

// [Exercise] errors
var (
	CouldNotCreateAllExercisesErr = errors.New("Could not create all exercises")
	CouldNotReadAllExercisesErr   = errors.New("Could not read all exercises")
	CouldNotUpdateAllExercisesErr = errors.New("Could not update all exercises")
	CouldNotDeleteAllExercisesErr = errors.New("Could not delete all exercises")
)

// [Hyperparams] errors
var (
	InvalidBarPathCalcErr    = errors.New("Invalid bar path calc conf")
	InvalidMinNumSamplesErr  = errors.New("Invalid min num samples")
	InvalidTimeDeltaEpsErr   = errors.New("Invalid time delta eps")
	InvalidNearZeroFilterErr = errors.New("Invalid near zero filter")
	InvalidNoiseFilterErr    = errors.New("Invalid noise filter")

	InvalidBarPathTrackerErr = errors.New("Invalid bar path tracker conf")
	InvalidMinLengthErr      = errors.New("Invalid min length")
	InvalidMaxFileSizeErr    = errors.New("Invalid max file size")

	CouldNotCreateAllHyperparamsErr = errors.New("Could not create all hyperparams")
	CouldNotReadAllHyperparamsErr   = errors.New("Could not read all hyperparams")
	CouldNotUpdateAllHyperparamsErr = errors.New("Could not update all hyperparams")
	CouldNotDeleteAllHyperparamsErr = errors.New("Could not delete all hyperparams")
)

// Job queue errors
var (
	CSVLoaderJobQueueErr = errors.New("Could not process csv loader job")

	PhysicsJobQueueErr        = errors.New("Could not process physics job")
	InvalidRawDataLenErr      = errors.New("Invalid raw data length")
	InvalidExpNumRepsErr      = errors.New("Invalid exp num reps")
	TimeSeriesDecreaseErr     = errors.New("Time series data must not decrease")
	TimeSeriesNotMonotonicErr = errors.New("Time series must increase mononically")
)
