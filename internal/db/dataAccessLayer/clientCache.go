package dal

import (
	"context"
	"unsafe"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	"github.com/maypok86/otter/v2"
)

func NewClientCacheLoader(
	queries *SyncQueries,
) otter.LoaderFunc[string, types.IdWrapper[int64, types.Client]] {
	return func(
		ctxt context.Context,
		key string,
	) (res types.IdWrapper[int64, types.Client], err error) {
		queries.Run(func(q *Queries) {
			// TODO - some kind of check would be nice...
			// _=types.IdWrapper[uint64, types.Client](Client{})
			var tmp Client
			tmp, err = q.GetFullClientByEmail(ctxt, key)
			res = *(*types.IdWrapper[int64, types.Client])(unsafe.Pointer(&tmp))
		})
		return
	}
}

func NewExerciseCacheLoader(
	queries *SyncQueries,
) otter.LoaderFunc[string, types.IdWrapper[int32, types.Exercise]] {
	return func(
		ctxt context.Context,
		key string,
	) (res types.IdWrapper[int32, types.Exercise], err error) {
		queries.Run(func(q *Queries) {
			// TODO - some kind of check would be nice...
			// _=types.IdWrapper[uint32, types.Exercise](Client{})
			var tmp Exercise
			tmp, err = q.GetFullExerciseByName(ctxt, key)
			res = *(*types.IdWrapper[int32, types.Exercise])(unsafe.Pointer(&tmp))
		})
		return
	}
}
