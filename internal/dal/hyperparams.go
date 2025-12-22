package dal

import (
	"context"
	"encoding/json"
	"fmt"

	"code.barbellmath.net/barbell-math/providentia/internal/util"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
	"github.com/jackc/pgx/v5"
)

type (
	ReadHyperparamsByVersionForOpts[T any] struct {
		Versions []int32
		Params   *[]T
	}

	FindHyperparamsByVersionForOpts[T any] struct {
		Versions []int32
		Params   *[]types.Found[T]
	}

	versionParamRes struct {
		Version int32  `db:"version"`
		Params  []byte `db:"params"`
	}
)

const (
	hyperparamsTableName = "hyperparams"

	readNumHyperparamsForSql = `
SELECT COUNT(*) FROM providentia.hyperparams WHERE model_id = $1;
`

	readHyperparamsByVersionForSql = `
SELECT version, params
FROM providentia.hyperparams
JOIN UNNEST($1::INT4[])
WITH ORDINALITY t(version, ord)
USING (version)
WHERE model_id=$2
ORDER BY ord;
`

	findHyperparamsByVersionForSql = `
SELECT ord::INT8, version, params
FROM providentia.hyperparams
JOIN UNNEST($1::INT4[])
WITH ORDINALITY t(version, ord)
USING (version)
WHERE model_id=$2
ORDER BY ord;
`

	deleteHyperparamsByVersionFor = `
DELETE FROM providentia.hyperparams WHERE model_id = $1 AND version = $2;
`
)

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

func validateHyperparams[T types.Hyperparams](v *T) (opErr error) {
	switch params := any(v).(type) {
	case *types.BarPathCalcHyperparams:
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
		if params.NoiseFilter <= 0 {
			return sberr.AppendError(
				types.InvalidBarPathCalcErr,
				sberr.Wrap(
					types.InvalidNoiseFilterErr,
					"Must be >0. Got: %d", params.NoiseFilter,
				),
			)
		}
	case *types.BarPathTrackerHyperparams:
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

func CreateHyperparams[T types.Hyperparams](
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	params []T,
) error {
	return genericCreate(
		ctxt, state, tx, &genericCreateOpts[T]{
			TableName: hyperparamsTableName,
			Columns:   []string{"model_id", "version", "params"},
			Data:      params,
			ValueGetter: func(v *T, res *[]any) error {
				*res = util.SliceClamp(*res, 3)
				if err := validateHyperparams(v); err != nil {
					return err
				}
				jsonParams, err := json.Marshal(v)
				(*res)[0] = getModelIdFor[T]()
				(*res)[1] = getVersionFrom(v)
				(*res)[2] = jsonParams
				return err
			},
			Err: types.CouldNotCreateAllHyperparamsErr,
		},
	)
}

func EnsureHyperparamsExist[T types.Hyperparams](
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	params []T,
) error {
	return genericEnsureExists(
		ctxt, state, tx, &genericCreateOpts[T]{
			TableName: hyperparamsTableName,
			Columns:   []string{"model_id", "version", "params"},
			Data:      params,
			ValueGetter: func(v *T, res *[]any) error {
				*res = make([]any, 3)
				if err := validateHyperparams(v); err != nil {
					return err
				}
				jsonParams, err := json.Marshal(v)
				(*res)[0] = getModelIdFor[T]()
				(*res)[1] = getVersionFrom(v)
				(*res)[2] = jsonParams
				return err
			},
			Err: types.CouldNotCreateAllHyperparamsErr,
		},
	)
}

func ReadNumHyperparams(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	num *int64,
) error {
	return genericReadTotalNum(
		ctxt, state, tx, &genericReadTotalNumOpts{
			TableName: hyperparamsTableName,
			Res:       num,
		},
	)
}

func ReadNumHyperparamsFor[T types.Hyperparams](
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	num *int64,
) error {
	modelId := getModelIdFor[T]()
	row := tx.QueryRow(ctxt, readNumHyperparamsForSql, modelId)
	state.Log.Log(
		ctxt, sblog.VLevel(3),
		fmt.Sprintf("DAL: Read total num hyperparams for %s", modelId),
	)
	return row.Scan(num)
}

func ReadHyperparamsByVersionFor[T types.Hyperparams](
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts ReadHyperparamsByVersionForOpts[T],
) error {
	*opts.Params = util.SliceClamp(*opts.Params, len(opts.Versions))
	modelId := getModelIdFor[T]()
	for start, end := range batchIndexes(opts.Versions, int(state.Global.BatchSize)) {
		select {
		case <-ctxt.Done():
			return ctxt.Err()
		default:
		}

		rows, err := tx.Query(
			ctxt, readHyperparamsByVersionForSql,
			opts.Versions[start:end], modelId,
		)
		if err != nil {
			return sberr.AppendError(types.CouldNotReadAllHyperparamsErr, err)
		}

		cntr := start
		for rows.Next() {
			iterRes, err := pgx.RowToStructByName[versionParamRes](rows)
			if err != nil {
				rows.Close()
				return sberr.AppendError(types.CouldNotReadAllHyperparamsErr, err)
			}

			setVersionTo(&(*opts.Params)[cntr], iterRes.Version)
			if err = json.Unmarshal(
				iterRes.Params, &(*opts.Params)[cntr],
			); err != nil {
				rows.Close()
				return sberr.AppendError(types.CouldNotReadAllHyperparamsErr, err)
			}
			cntr++
		}
		rows.Close()

		if cntr != end {
			return sberr.Wrap(
				types.CouldNotReadAllHyperparamsErr,
				"Only read %d entries out of batch of %d requests",
				cntr-start, end-start,
			)
		}

		state.Log.Log(
			ctxt, sblog.VLevel(3),
			fmt.Sprintf("DAL: Read hyperparams by version for %s", modelId),
			"NumRows", end-start,
		)
	}
	return nil
}

func ReadDefaultHyperparamsFor[T types.Hyperparams](
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	res *T,
) error {
	var tmp []T
	opts := ReadHyperparamsByVersionForOpts[T]{
		Versions: []int32{0},
		Params:   &tmp,
	}
	err := ReadHyperparamsByVersionFor(ctxt, state, tx, opts)
	if err != nil {
		return err
	}
	if len(tmp) != 1 {
		state.Log.Error(
			"Expected 1 result for a default hyperparameter but got more, database is not consistent with what was expected!",
			"NumResults", len(tmp),
		)
		err = sberr.Wrap(
			types.CouldNotReadAllHyperparamsErr,
			"Expected 1 result for a default hyperparameter but got %d, database is not consistent with what was expected",
			len(tmp),
		)
	}
	*res = tmp[0]

	return nil
}

func FindHyperparamsByVersionFor[T types.Hyperparams](
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts FindHyperparamsByVersionForOpts[T],
) error {
	*opts.Params = util.SliceClamp(*opts.Params, len(opts.Versions))
	modelId := getModelIdFor[T]()
	for start, end := range batchIndexes(opts.Versions, int(state.Global.BatchSize)) {
		select {
		case <-ctxt.Done():
			return ctxt.Err()
		default:
		}

		rows, err := tx.Query(
			ctxt, findHyperparamsByVersionForSql,
			opts.Versions[start:end], modelId,
		)
		if err != nil {
			return sberr.AppendError(types.CouldNotReadAllHyperparamsErr, err)
		}

		var iterVal versionParamRes
		ord := int64(0)
		found := int64(0)
		scanValues := []any{&ord, &iterVal.Version, &iterVal.Params}

		for rows.Next() {
			if err := rows.Scan(scanValues...); err != nil {
				rows.Close()
				return sberr.AppendError(types.CouldNotReadAllHyperparamsErr, err)
			}

			idx := int64(start) + ord - 1
			(*opts.Params)[idx].Found = true
			setVersionTo(&((*opts.Params)[idx].Value), iterVal.Version)
			if err = json.Unmarshal(
				iterVal.Params, &(*opts.Params)[idx].Value,
			); err != nil {
				rows.Close()
				return sberr.AppendError(types.CouldNotReadAllHyperparamsErr, err)
			}
			found++
		}
		rows.Close()

		state.Log.Log(
			ctxt, sblog.VLevel(3),
			fmt.Sprintf(
				"DAL: Found hyperparams by version for %s", modelId,
			),
			"NumFound/NumRows", fmt.Sprintf("%d/%d", found, end-start),
		)
	}
	return nil
}

func DeleteHyperparams[T types.Hyperparams](
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	versions []int32,
) error {
	// Deleting all referenced/referencing data is handled by cascade rules

	modelId := getModelIdFor[T]()
	for start, end := range batchIndexes(versions, int(state.Global.BatchSize)) {
		select {
		case <-ctxt.Done():
			return ctxt.Err()
		default:
		}

		b := pgx.Batch{}
		for i := start; i < end; i++ {
			b.Queue(deleteHyperparamsByVersionFor, modelId, versions[i])
		}
		results := tx.SendBatch(ctxt, &b)

		for i := start; i < end; i++ {
			if cmdTag, err := results.Exec(); err != nil {
				results.Close()
				return sberr.AppendError(
					types.CouldNotDeleteAllHyperparamsErr, err,
				)
			} else if cmdTag.RowsAffected() == 0 {
				results.Close()
				return sberr.Wrap(
					types.CouldNotDeleteAllHyperparamsErr,
					"Could not delete entry with version '%d' (Does it exist?)",
					versions[i],
				)
			}
		}
		results.Close()

		state.Log.Log(
			ctxt, sblog.VLevel(3),
			fmt.Sprintf("DAL: Deleted hyperparams for %s", modelId),
			"NumRows", end-start,
		)
	}
	return nil
}
