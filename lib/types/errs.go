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

// Job queue errors
var (
	CSVLoaderJobQueueErr = errors.New("Could not process csv loader job")
)
