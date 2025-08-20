package dal

import (
	"sync"
)

type (
	SyncQueries struct {
		mut     sync.Mutex
		queries *Queries
	}
)

func NewSyncQueries(q *Queries) *SyncQueries {
	return &SyncQueries{queries: q}
}

func (s *SyncQueries) Run(op func(q *Queries)) {
	s.mut.Lock()
	op(s.queries)
	s.mut.Unlock()
}
