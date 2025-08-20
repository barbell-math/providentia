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
func WithStateValue(
	ctxt context.Context,
	s *types.State,
) (newCtxt context.Context, opErr error) {
	newCtxt = context.WithValue(ctxt, stateCtxtKey, s)
	opErr = validateState(s)
	return
}

func validateState(s *types.State) error {
	if s.PhysicsData.MinNumSamples < 2 {
		return sberr.AppendError(
			types.InvalidPhysicsDataConfErr,
			sberr.Wrap(
				types.InvalidMinNumSamplesErr,
				"Must be >=2. Got: %d", s.PhysicsData.MinNumSamples,
			),
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
