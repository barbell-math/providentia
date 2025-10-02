package logic

import (
	"context"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
)

type (
	ctxtKey struct{}
)

var (
	stateCtxtKey ctxtKey
)

// Returns a [State] from the supplied context or nil if it was not present. The
// boolean flag indicates if the [State] value was present.
func StateFromContext(ctxt context.Context) (*types.State, bool) {
	s, ok := ctxt.Value(stateCtxtKey).(*types.State)
	return s, ok
}

// Adds the supplied state value to the supplied context, returning a new
// context with the state value.
//
// Most other library functions require the supplied context to hold a
// [types.State] value, which will require calling this function.
//
// The state value must be valid or other library functions will encounter
// unexpected errors. If a [types.State] struct is manually created rather than
// using [ConfToState] the [ValidateState] function should be called before
// making any other library calls.
func WithStateValue(ctxt context.Context, s *types.State) context.Context {
	return context.WithValue(ctxt, stateCtxtKey, s)
}

// Validates the supplied state. If this function does not return `nil` then
// other library functions are likely to error or have unexpected results.
//
// If a [types.State] struct is manually created rather than using [ConfToState]
// this function should be called before making any other library calls.
func ValidateState(s *types.State) error {
	if err := checkStateGlobalConf(s); err != nil {
		return err
	}
	if s.Log == nil {
		return sberr.Wrap(types.InvalidLoggerErr, "The Log field must not be nil")
	}
	if s.DB == nil {
		return sberr.Wrap(types.InvalidDBErr, "The DB field must not be nil")
	}
	if s.PhysicsJobQueue == nil {
		return sberr.Wrap(
			types.InvalidPhysicsJobQueueErr,
			"The PhysicsJobQueue field must not be nil",
		)
	}
	if s.VideoJobQueue == nil {
		return sberr.Wrap(
			types.InvalidVideoJobQueue,
			"The VideoJobQueue field must not be nil",
		)
	}
	if s.CSVLoaderJobQueue == nil {
		return sberr.Wrap(
			types.InvalidCSVLoaderJobQueueErr,
			"The CSVLoaderJobQueue field must not be nil",
		)
	}
	return nil
}

func checkStateGlobalConf(state *types.State) error {
	if state.Global.BatchSize == 0 {
		return sberr.AppendError(
			types.InvalidGlobalErr,
			sberr.Wrap(
				types.InvalidBatchSizeErr,
				"Must be >0. Got: %d", state.Global.BatchSize,
			),
		)
	}
	if state.Global.BatchSize > 1e5 {
		state.Log.Warn(
			"Large batch sizes can lead to OOM errors and will limit the " +
				"number of concurrent jobs running at once",
		)
	}
	return nil
}

// Cleans up the resources in the supplied state.
func CleanupState(s *types.State) {
	if s.DB != nil {
		s.DB.Close()
	}
}
