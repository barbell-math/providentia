package logic

import (
	"context"
	"reflect"
	"testing"
	"time"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbtest "code.barbellmath.net/barbell-math/smoothbrain-test"
)

func TestCancelation(t *testing.T) {
	timeoutCtxt, cancel := context.WithTimeout(
		context.Background(), 500*time.Millisecond,
	)
	defer cancel()

	ctxt, cleanup := resetApp(timeoutCtxt)
	t.Cleanup(cleanup)

	err := runOp(ctxt, opCalls{
		op: func(state *types.State, queries *dal.SyncQueries) error {
			time.Sleep(5 * time.Second)
			return nil
		},
	})
	sbtest.ContainsError(t, context.DeadlineExceeded, err)
}

// TODO - this needs to move to sbtest and check for things like embeded field
// equivalency and generic value equivalencacy
func structsEquivalent[T any, U any](t *testing.T) {
	tRef := reflect.TypeFor[T]()
	uRef := reflect.TypeFor[T]()

	sbtest.Eq(t, tRef.Kind(), reflect.Struct)
	sbtest.Eq(t, uRef.Kind(), reflect.Struct)
	sbtest.Eq(t, tRef.NumField(), uRef.NumField())

	for i := range tRef.NumField() {
		sbtest.Eq(t, tRef.Field(i).Type.Kind(), uRef.Field(i).Type.Kind())
	}
}
