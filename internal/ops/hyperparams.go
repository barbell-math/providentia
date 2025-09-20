package ops

import (
	"context"
	"encoding/json"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
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
			opErr = sberr.AppendError(types.CouldNotAddNumHyperparamsErr, opErr)
			return
		}
	}

	if opErr = bufWriter.Flush(ctxt, queries); opErr != nil {
		opErr = sberr.AppendError(types.CouldNotAddNumHyperparamsErr, opErr)
		return
	}

	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"Added new hyperparams",
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

func ReadNumHyperparams(
	ctxt context.Context,
	state *types.State,
	queries *dal.SyncQueries,
) (res int64, opErr error) {
	res, opErr = dal.Query0x2(dal.Q.GetNumHyperparams, queries, ctxt)
	if opErr != nil {
		opErr = sberr.AppendError(types.CouldNotGetNumHyperparamsErr, opErr)
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
		opErr = sberr.AppendError(types.CouldNotGetNumHyperparamsErr, opErr)
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
				types.CouldNotFindRequestedHyperparamsErr, opErr,
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
		opErr = sberr.AppendError(types.CouldNotDeleteHyperparamsErr, opErr)
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
