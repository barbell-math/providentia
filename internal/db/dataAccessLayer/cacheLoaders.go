package dal

import (
	"context"
)

type (
	// Used to map unique fields to database ids (primary keys). This is often
	// useful when creating a new row that has foreign keys that reference a
	// tables id field.
	//
	// Note: this cache is not safe for concurrent usage!
	IdCache[K comparable, V ~int32 | ~int64] struct {
		curSize int
		curIdx  int
		keys    []K
		vals    []V
		lookup  map[K]*V
		loader  func(*Queries, context.Context, K) (V, error)
	}
)

func NewClientIdCache(maxSize uint) IdCache[string, int64] {
	return IdCache[string, int64]{
		vals:   make([]int64, maxSize),
		keys:   make([]string, maxSize),
		lookup: map[string]*int64{},
		loader: Q.GetClientIdByEmail,
	}
}

func NewExerciseIdCache(maxSize uint) IdCache[string, int32] {
	return IdCache[string, int32]{
		vals:   make([]int32, maxSize),
		keys:   make([]string, maxSize),
		lookup: map[string]*int32{},
		loader: Q.GetExerciseIdByName,
	}
}

func (i *IdCache[K, V]) Get(
	ctxt context.Context,
	queries *SyncQueries,
	key K,
) (V, error) {
	if v, ok := i.lookup[key]; ok {
		return *v, nil
	}

	tmp, err := Query1x2(i.loader, queries, ctxt, key)
	if err != nil {
		return tmp, err
	}

	delete(i.lookup, i.keys[i.curIdx])
	i.keys[i.curIdx] = key
	i.vals[i.curIdx] = tmp
	i.lookup[key] = &i.vals[i.curIdx]
	i.curIdx = (i.curIdx + 1) % len(i.vals)
	return tmp, nil
}
