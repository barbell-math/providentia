package ops

import (
	"context"
	"errors"

	dal "github.com/barbell-math/providentia/internal/db/dataAccessLayer"
	sberr "github.com/barbell-math/smoothbrain-errs"
)

var (
	IncorrectNumberOfRowsErr = errors.New("Incorrect number of rows were written")
)

type (
	bulkCreateTypes interface {
		dal.BulkCreateClientsParams | dal.BulkCreateTrainingLogParams
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
