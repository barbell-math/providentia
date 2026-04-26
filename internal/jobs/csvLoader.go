package jobs

import (
	"context"
	"iter"
	"os"

	"code.barbellmath.net/barbell-math/providentia/internal/dal"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	"code.barbellmath.net/carmichaeljr/smoothbrain/sbcsv"
	"code.barbellmath.net/carmichaeljr/smoothbrain/sberrs"
	"code.barbellmath.net/carmichaeljr/smoothbrain/sbjobqueue"
	"code.barbellmath.net/carmichaeljr/smoothbrain/sblog"
	"github.com/jackc/pgx/v5"
)

type (
	genericCSVAvailableTypes interface {
		types.Client | types.Exercise | types.Hyperparams
	}

	genericCSVLoader[T genericCSVAvailableTypes] struct {
		B         *sbjobqueue.Batch
		S         *types.State
		Tx        pgx.Tx
		UID       int64
		File      string
		Opts      *sbcsv.Opts
		WriteFunc dal.CreateFunc[T]
	}

	CSVLoaderOpts[T genericCSVAvailableTypes] struct {
		*sbcsv.Opts
		Creator dal.CreateFunc[T]
		Files   iter.Seq2[string, error]
		Batch   *sbjobqueue.Batch
	}
)

func UploadFromCSV[T genericCSVAvailableTypes](
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts *CSVLoaderOpts[T],
) error {
	wait := false
	if opts.Batch == nil {
		wait = true
		opts.Batch, _ = sbjobqueue.BatchWithContext(ctxt)
	}

	for file, err := range opts.Files {
		if err != nil {
			return err
		}
		select {
		case <-ctxt.Done():
			return ctxt.Err()
		default:
		}

		state.Log.Log(
			ctxt, sblog.VLevel(3),
			formatJobLogLine("UploadFromCSV", -1, "Processing data file"),
			"File", file,
		)
		state.CSVLoaderJobQueue.Schedule(&genericCSVLoader[T]{
			S:         state,
			Tx:        tx,
			B:         opts.Batch,
			UID:       UID_CNTR.Add(1),
			File:      file,
			Opts:      opts.Opts,
			WriteFunc: opts.Creator,
		})
	}

	if wait {
		return opts.Batch.Wait()
	}
	return nil
}

func (w *genericCSVLoader[T]) JobType(_ types.CSVLoaderJob) {}

func (w *genericCSVLoader[T]) Batch() *sbjobqueue.Batch {
	return w.B
}

func (w *genericCSVLoader[T]) formatLogLine(msg string) string {
	return formatJobLogLine("genericCSVLoader", w.UID, msg)
}

func (w *genericCSVLoader[T]) Run(ctxt context.Context) (opErr error) {
	w.S.Log.Log(ctxt, sblog.VLevel(3), w.formatLogLine("Starting..."))
	var f *os.File
	params := []T{}

	f, opErr = os.Open(w.File)
	if opErr != nil {
		goto errReturn
	}

	if opErr = sbcsv.LoadReader(f, &sbcsv.LoadOpts{
		Opts:          *w.Opts,
		RequestedCols: sbcsv.ReqColsForStruct[T](),
		Op:            sbcsv.RowToStructOp(&params),
	}); opErr != nil {
		goto errReturn
	}

	// This is unfortunate... but it has to be done because a single transaction
	// is backed by a single conn which is not thread safe.
	w.B.Lock()
	if opErr = w.WriteFunc(ctxt, w.S, w.Tx, params); opErr != nil {
		goto errReturn
	}
	w.B.Unlock()

	w.S.Log.Log(
		ctxt, sblog.VLevel(3),
		w.formatLogLine("Finished loading data"),
		"NumRows", len(params),
	)
	return
errReturn:
	w.S.Log.Error(w.formatLogLine("Encountered error"), "Error", opErr)
	return sberr.AppendError(types.CSVLoaderJobQueueErr, opErr)
}
