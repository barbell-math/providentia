package dal

import (
	"context"
	"errors"

	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
)

var (
	IncorrectNumberOfRowsErr = errors.New("Incorrect number of rows were written")
	IndexOutOfRangeErr       = errors.New("Index out of range")
)

type (
	bulkCreateTypes interface {
		BulkCreateClientsParams |
			BulkCreateTrainingLogsParams |
			BulkCreateHyperparamsParams
	}

	BufferedWriter[T bulkCreateTypes] struct {
		data     []T
		maxElems uint
		curIdx   int
		writeOp  func(q *Queries, ctxt context.Context, arg []T) (int64, error)
		preOp    func() error
	}
)

func NewBufferedWriter[T bulkCreateTypes](
	size uint,
	writeOp func(q *Queries, ctxt context.Context, arg []T) (int64, error),
	preOp func() error,
) BufferedWriter[T] {
	return BufferedWriter[T]{
		data:     make([]T, size),
		maxElems: size,
		curIdx:   -1,
		preOp:    preOp,
		writeOp:  writeOp,
	}
}

func (b *BufferedWriter[T]) Write(
	ctxt context.Context,
	queries *SyncQueries,
	data ...T,
) error {
	for _, d := range data {
		if uint(b.curIdx+1) < b.maxElems {
			b.data[b.curIdx+1] = d
			b.curIdx++
		} else if err := b.Flush(ctxt, queries); err != nil {
			return err
		}
	}
	return nil
}

func (b *BufferedWriter[T]) Pntr(idx int) (*T, error) {
	if idx < 0 || idx > b.curIdx {
		return nil, sberr.Wrap(
			IndexOutOfRangeErr,
			"Must be in range [0, %d], Got: %d", b.curIdx, idx,
		)
	}
	return &b.data[idx], nil
}

func (b *BufferedWriter[T]) Last() *T {
	if b.curIdx >= 0 {
		return &b.data[b.curIdx]
	}
	return nil
}

func (b *BufferedWriter[T]) Flush(
	ctxt context.Context,
	queries *SyncQueries,
) error {
	if b.curIdx == -1 {
		return nil
	}

	if err := b.preOp(); err != nil {
		return err
	}

	numRows, err := Query1x2(b.writeOp, queries, ctxt, b.data[0:b.curIdx+1])
	if err != nil {
		return err
	} else if numRows != int64(b.curIdx)+1 {
		return sberr.Wrap(
			IncorrectNumberOfRowsErr,
			"Expected: %d Got: %d", b.curIdx+1, numRows,
		)
	}
	b.curIdx = -1
	return nil
}
