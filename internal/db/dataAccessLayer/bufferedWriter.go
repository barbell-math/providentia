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
		BulkCreateClientsParams | BulkCreateTrainingLogsParams
	}

	BufferedWriter[T bulkCreateTypes] struct {
		data     []T
		maxElems uint
		curIdx   int
		writeOp  func(ctxt context.Context, arg []T) (int64, error)
	}
)

func NewBufferedWriter[T bulkCreateTypes](
	size uint,
	writeOp func(ctxt context.Context, arg []T) (int64, error),
) BufferedWriter[T] {
	return BufferedWriter[T]{
		data:     make([]T, size),
		maxElems: size,
		curIdx:   -1,
		writeOp:  writeOp,
	}
}

func (b *BufferedWriter[T]) Write(ctxt context.Context, data ...T) error {
	for _, d := range data {
		if uint(b.curIdx+1) < b.maxElems {
			b.data[b.curIdx+1] = d
			b.curIdx++
		} else if err := b.Flush(ctxt); err != nil {
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

func (b *BufferedWriter[T]) Flush(ctxt context.Context) error {
	if b.curIdx == -1 {
		return nil
	}

	numRows, err := b.writeOp(ctxt, b.data[0:b.curIdx+1])
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
