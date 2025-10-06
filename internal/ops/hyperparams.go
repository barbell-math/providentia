package ops

import (
	"context"
	"encoding/json"
	"math"
	"runtime"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/internal/jobs"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbcsv "code.barbellmath.net/barbell-math/smoothbrain-csv"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sbjobqueue "code.barbellmath.net/barbell-math/smoothbrain-jobQueue"
	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
)

func CreateHyperparams[T types.Hyperparams](
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	params ...T,
) (opErr error) {
	modelId := getModelIdFor[T]()
	bufWriter := dal.NewBufferedWriter(
		state.Global.BatchSize,
		dal.Q.BulkCreateHyperparams,
		func() (err error) { return },
	)

	for _, iterParams := range params {
		select {
		case <-ctxt.Done():
			opErr = ctxt.Err()
			return
		default:
		}

		if opErr = validateHyperparams(iterParams); opErr != nil {
			opErr = sberr.AppendError(types.InvalidHyperparamsErr, opErr)
			return
		}

		var jsonParams []byte
		jsonParams, opErr = json.Marshal(iterParams)
		if opErr != nil {
			opErr = sberr.AppendError(
				types.InvalidHyperparamsErr,
				types.EncodingJsonHyperparamsErr,
				opErr,
			)
		}
		if opErr = bufWriter.Write(ctxt, queries, dal.BulkCreateHyperparamsParams{
			ModelID: modelId,
			Version: getVersionFrom(&iterParams),
			Params:  jsonParams,
		}); opErr != nil {
			opErr = sberr.AppendError(
				types.CouldNotAddNumHyperparamsErr, dal.FormatErr(opErr),
			)
			return
		}
	}

	if opErr = bufWriter.Flush(ctxt, queries); opErr != nil {
		opErr = sberr.AppendError(
			types.CouldNotAddNumHyperparamsErr, dal.FormatErr(opErr),
		)
		return
	}

	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"Added new hyperparams",
		"NumHyperparams", len(params),
	)
	return
}

func EnsureHyperparamsExist[T types.Hyperparams](
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	params ...T,
) (opErr error) {
	modelId := getModelIdFor[T]()
	bufWriter := dal.NewBufferedWriter(
		state.Global.BatchSize,
		dal.Q.BulkCreateHyperparams,
		func() (err error) { return },
	)
	versions := make([]int32, min(len(params), int(state.Global.BatchSize)))

	for start, end := range batchIndexes(params, int(state.Global.BatchSize)) {
		select {
		case <-ctxt.Done():
			opErr = ctxt.Err()
			return
		default:
		}

		chunk := params[start:end]
		for i := range end - start {
			versions[i] = getVersionFrom(&chunk[i])
		}

		// There is no way to express deep-JSON equality in SQL, so we have to
		// make a bunch of round trips to the database :(
		var existingParams []types.Found[T]
		existingParams, opErr = FindHyperparamsByVersionFor[T](
			ctxt, state, queries, versions[:len(chunk)]...,
		)
		if opErr != nil {
			return
		}

		for i, param := range existingParams[:len(chunk)] {
			if param.Found && param.Value == params[i+start] {
				continue
			}

			var jsonParams []byte
			jsonParams, opErr = json.Marshal(params[i+start])
			if opErr != nil {
				opErr = sberr.AppendError(
					types.InvalidHyperparamsErr,
					types.EncodingJsonHyperparamsErr,
					opErr,
				)
			}
			if opErr = bufWriter.Write(
				ctxt, queries,
				dal.BulkCreateHyperparamsParams{
					ModelID: modelId,
					Version: getVersionFrom(&params[i+start]),
					Params:  jsonParams,
				},
			); opErr != nil {
				opErr = sberr.AppendError(
					types.CouldNotAddNumHyperparamsErr, dal.FormatErr(opErr),
				)
				return
			}
		}
	}

	if opErr = bufWriter.Flush(ctxt, queries); opErr != nil {
		opErr = sberr.AppendError(
			types.CouldNotAddNumHyperparamsErr, dal.FormatErr(opErr),
		)
		return
	}

	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"Ensured hyperparams exist",
		"NumHyperparams", len(params),
	)
	return
}

func validateHyperparams[T types.Hyperparams](v T) (opErr error) {
	switch params := any(v).(type) {
	case types.BarPathCalcHyperparams:
		if params.MinNumSamples < 2 {
			return sberr.AppendError(
				types.InvalidBarPathCalcErr,
				sberr.Wrap(
					types.InvalidMinNumSamplesErr,
					"Must be >=2. Got: %d", params.MinNumSamples,
				),
			)
		}
		if params.TimeDeltaEps <= 0 {
			return sberr.AppendError(
				types.InvalidBarPathCalcErr,
				sberr.Wrap(
					types.InvalidTimeDeltaEpsErr,
					"Must be >=0. Got: %f", params.TimeDeltaEps,
				),
			)
		}
		if !params.ApproxErr.IsValid() {
			return sberr.AppendError(
				types.InvalidBarPathCalcErr,
				types.ErrInvalidApproximationError,
			)
		}
		if params.NearZeroFilter < 0 {
			return sberr.AppendError(
				types.InvalidBarPathCalcErr,
				sberr.Wrap(
					types.InvalidNearZeroFilterErr,
					"Must be >0. Got: %f", params.NearZeroFilter,
				),
			)
		}
	case types.BarPathTrackerHyperparams:
		if params.MinLength < 0 {
			return sberr.AppendError(
				types.InvalidBarPathTrackerErr,
				sberr.Wrap(
					types.InvalidMinLengthErr,
					"Must be >0. Got: %f", params.MinLength,
				),
			)
		}
		if params.MaxFileSize <= params.MinFileSize {
			return sberr.AppendError(
				types.InvalidBarPathTrackerErr,
				sberr.Wrap(
					types.InvalidMaxFileSizeErr,
					"Max size (%d) must be > min size (%d)",
					params.MaxFileSize,
					params.MinFileSize,
				),
			)
		}
	}
	return nil
}

func getModelIdFor[T types.Hyperparams]() types.ModelID {
	switch any((*T)(nil)).(type) {
	case *types.BarPathCalcHyperparams:
		return types.BarPathCalc
	case *types.BarPathTrackerHyperparams:
		return types.BarPathTracker
	}
	return types.UnknownModel
}

func getVersionFrom[T types.Hyperparams](v *T) int32 {
	switch params := any(v).(type) {
	case *types.BarPathCalcHyperparams:
		return params.Version
	case *types.BarPathTrackerHyperparams:
		return params.Version
	}
	return 0
}

func setVersionTo[T types.Hyperparams](v *T, version int32) {
	switch params := any(v).(type) {
	case *types.BarPathCalcHyperparams:
		params.Version = version
	case *types.BarPathTrackerHyperparams:
		params.Version = version
	}
}

func UploadHyperparamsFromCSV[T types.Hyperparams](
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	creator createFunc[T],
	opts sbcsv.Opts,
	files ...string,
) (opErr error) {
	opts.ReuseRecord = true
	batch, _ := sbjobqueue.BatchWithContext(ctxt)

	for _, file := range files {
		var fileChunks [][]byte
		if fileChunks, opErr = sbcsv.ChunkFile(
			file, sbcsv.ChunkFileOpts{
				NumRowSamples:      2,
				MinChunkRows:       1e5,
				MaxChunkRows:       math.MaxInt,
				RequestedNumChunks: runtime.NumCPU(),
			},
		); opErr != nil {
			return
		}
		for _, chunk := range fileChunks {
			state.CSVLoaderJobQueue.Schedule(&jobs.CSVLoader[T]{
				S:         state,
				Q:         queries,
				B:         batch,
				FileChunk: chunk,
				Opts:      &opts,
				WriteFunc: creator,
			})
		}
	}

	return batch.Wait()
}

func ReadNumHyperparams(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
) (res int64, opErr error) {
	res, opErr = dal.Query0x2(dal.Q.GetNumHyperparams, queries, ctxt)
	if opErr != nil {
		opErr = sberr.AppendError(
			types.CouldNotGetNumHyperparamsErr, dal.FormatErr(opErr),
		)
		return
	}
	state.Log.Log(ctxt, sblog.VLevel(3), "Read num hyperparams")
	return
}

func ReadNumHyperparamsFor[T types.Hyperparams](
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
) (res int64, opErr error) {
	modelId := getModelIdFor[T]()
	res, opErr = dal.Query1x2(dal.Q.GetNumHyperparamsFor, queries, ctxt, modelId)
	if opErr != nil {
		opErr = sberr.AppendError(
			types.CouldNotGetNumHyperparamsErr, dal.FormatErr(opErr),
		)
		return
	}
	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"Read num hyperparams",
		"ModelID", modelId,
	)
	return
}

func ReadHyperparamsByVersionFor[T types.Hyperparams](
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	versions ...int32,
) (res []T, opErr error) {
	res = make([]T, len(versions))
	modelId := getModelIdFor[T]()

	for start, end := range batchIndexes(versions, int(state.Global.BatchSize)) {
		select {
		case <-ctxt.Done():
			opErr = ctxt.Err()
			return
		default:
		}

		var rawData []dal.GetHyperparamsByVersionForRow
		rawData, opErr = dal.Query1x2(
			dal.Q.GetHyperparamsByVersionFor, queries, ctxt,
			dal.GetHyperparamsByVersionForParams{
				ModelID:  modelId,
				Versions: versions[start:end],
			},
		)
		if opErr != nil {
			opErr = sberr.AppendError(
				types.CouldNotFindRequestedHyperparamsErr, dal.FormatErr(opErr),
			)
			return
		}
		if len(rawData) != end-start {
			opErr = types.CouldNotFindRequestedHyperparamsErr
			return
		}

		for i, iterRawData := range rawData {
			setVersionTo(&res[i+start], iterRawData.Version)
			if opErr = json.Unmarshal(
				iterRawData.Params, &res[i+start],
			); opErr != nil {
				opErr = sberr.AppendError(
					types.DecodingJsonHyperparamsErr, opErr,
				)
				return
			}
		}

		state.Log.Log(
			ctxt, sblog.VLevel(3),
			"Read hyperparams by version",
			"ModelID", modelId,
			"Num", len(rawData),
		)
	}

	return
}

func FindHyperparamsByVersionFor[T types.Hyperparams](
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	versions ...int32,
) (res []types.Found[T], opErr error) {
	res = make([]types.Found[T], len(versions))
	modelId := getModelIdFor[T]()

	for start, end := range batchIndexes(versions, int(state.Global.BatchSize)) {
		select {
		case <-ctxt.Done():
			opErr = ctxt.Err()
			return
		default:
		}

		var rawData []dal.FindHyperparamsByVersionForRow
		rawData, opErr = dal.Query1x2(
			dal.Q.FindHyperparamsByVersionFor, queries, ctxt,
			dal.FindHyperparamsByVersionForParams{
				ModelID:  modelId,
				Versions: versions[start:end],
			},
		)
		if opErr != nil {
			opErr = sberr.AppendError(
				types.CouldNotFindRequestedHyperparamsErr, dal.FormatErr(opErr),
			)
			return
		}

		rawDataIdx := 0
		for i := 0; i < end-start; i++ {
			res[i+start].Found = (rawDataIdx < len(rawData) && rawData[rawDataIdx].Ord-1 == int64(i))
			if res[i+start].Found {
				setVersionTo(&res[i+start].Value, rawData[rawDataIdx].Version)
				if opErr = json.Unmarshal(
					rawData[rawDataIdx].Params, &res[i+start].Value,
				); opErr != nil {
					opErr = sberr.AppendError(
						types.DecodingJsonHyperparamsErr, opErr,
					)
					return
				}
				rawDataIdx++
			}
		}

		state.Log.Log(
			ctxt, sblog.VLevel(3),
			"Found hyperparams by version",
			"ModelID", modelId,
			"Num", len(rawData),
		)
	}

	return
}

func DeleteHyperparams[T types.Hyperparams](
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
	versions ...int32,
) (opErr error) {
	// Deleting all referenced/referencing data is handled by cascade rules

	modelId := getModelIdFor[T]()

	var count int64
	count, opErr = dal.Query1x2(
		dal.Q.DeleteHyperparamsByVersionFor, queries, ctxt,
		dal.DeleteHyperparamsByVersionForParams{
			ModelID:  modelId,
			Versions: versions,
		},
	)
	if opErr != nil {
		opErr = sberr.AppendError(
			types.CouldNotDeleteHyperparamsErr, dal.FormatErr(opErr),
		)
		return
	}
	if count != int64(len(versions)) {
		opErr = sberr.AppendError(
			types.CouldNotDeleteHyperparamsErr,
			types.CouldNotFindRequestedHyperparamsErr,
		)
	}

	state.Log.Log(ctxt, sblog.VLevel(3), "Deleted hyperparams", "Num", count)
	return
}
