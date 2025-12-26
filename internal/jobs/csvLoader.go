package jobs

import (
	"context"
	"io"
	"iter"

	"code.barbellmath.net/barbell-math/providentia/internal/dal"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbcsv "code.barbellmath.net/barbell-math/smoothbrain-csv"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sbjobqueue "code.barbellmath.net/barbell-math/smoothbrain-jobQueue"
	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
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
		UID       uint64
		FileChunk io.Reader
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
) (opErr error) {
	wait := false
	if opts.Batch == nil {
		wait = true
		opts.Batch, _ = sbjobqueue.BatchWithContext(ctxt)
	}

	var file string
	for file, opErr = range opts.Files {
		if opErr != nil {
			return
		}
		state.Log.Log(
			ctxt, sblog.VLevel(3),
			"JOB: UploadFromCSV: Processing data file", "File", file,
		)
		if opErr = uploadFile(ctxt, state, tx, file, opts); opErr != nil {
			return
		}
	}

	if wait {
		return opts.Batch.Wait()
	}
	return nil
}

func uploadFile[T genericCSVAvailableTypes](
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	file string,
	opts *CSVLoaderOpts[T],
) (opErr error) {
	select {
	case <-ctxt.Done():
		return ctxt.Err()
	default:
	}

	var fileChunks []*sbcsv.BasicFileChunk
	if fileChunks, opErr = sbcsv.ChunkFile(
		file, sbcsv.NewBasicFileChunk, state.ClientCSVFileChunks,
	); opErr != nil {
		return
	}
	for _, chunk := range fileChunks {
		if len(chunk.Data) == 0 {
			continue
		}
		state.CSVLoaderJobQueue.Schedule(&genericCSVLoader[T]{
			S:         state,
			Tx:        tx,
			B:         opts.Batch,
			UID:       UID_CNTR.Add(1),
			FileChunk: chunk,
			Opts:      opts.Opts,
			WriteFunc: opts.Creator,
		})
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

	params := []T{}
	if opErr = sbcsv.LoadReader(w.FileChunk, &sbcsv.LoadOpts{
		Opts:          *w.Opts,
		RequestedCols: sbcsv.ReqColsForStruct[T](),
		Op:            sbcsv.RowToStructOp(&params),
	}); opErr != nil {
		goto errReturn
	}

	if opErr = w.WriteFunc(ctxt, w.S, w.Tx, params); opErr != nil {
		goto errReturn
	}

	w.S.Log.Log(
		ctxt, sblog.VLevel(3),
		w.formatLogLine("Finished loading clients"),
		"NumRows", len(params),
	)
	return
errReturn:
	w.S.Log.Error(w.formatLogLine("Encountered error"), "Error", opErr)
	return sberr.AppendError(types.CSVLoaderJobQueueErr, opErr)
}
