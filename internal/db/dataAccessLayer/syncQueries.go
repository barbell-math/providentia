package dal

import (
	"context"
	"sync"
)

type (
	SyncQueries struct {
		mut     sync.Mutex
		queries *Queries
	}

	Q = *Queries
)

func NewSyncQueries(q *Queries) *SyncQueries {
	return &SyncQueries{queries: q}
}

func Query0x2[OP ~func(*Queries, context.Context) (T, U), T any, U any](
	op OP,
	q *SyncQueries,
	ctxt context.Context,
) (rv1 T, rv2 U) {
	q.mut.Lock()
	rv1, rv2 = op(q.queries, ctxt)
	q.mut.Unlock()
	return
}

func Query1x1[OP ~func(*Queries, context.Context, T) U, T any, U any](
	op OP,
	q *SyncQueries,
	ctxt context.Context,
	val T,
) (rv1 U) {
	q.mut.Lock()
	rv1 = op(q.queries, ctxt, val)
	q.mut.Unlock()
	return
}

func Query1x2[OP ~func(*Queries, context.Context, T) (U, V), T any, U any, V any](
	op OP,
	q *SyncQueries,
	ctxt context.Context,
	val T,
) (rv1 U, rv2 V) {
	q.mut.Lock()
	rv1, rv2 = op(q.queries, ctxt, val)
	q.mut.Unlock()
	return
}
