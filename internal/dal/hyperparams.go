package dal

import (
	"context"
	"encoding/json"
	"fmt"

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
)

const (
	hyperparamsTableName = "hyperparams"

	readHyperparamsByVersionForSql = `
SELECT
	providentia.hyperparams.version,
	providentia.hyperparams.params
FROM providentia.hyperparams
JOIN UNNEST($1::INT4[])
WITH ORDINALITY t(version, ord)
USING (version)
WHERE model_id=$2
ORDER BY ord;
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
				if len(*res) < 3 {
					*res = make([]any, 3)
				}
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

func ReadHyperparamsByVersionFor[T types.Hyperparams](
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts ReadHyperparamsByVersionForOpts[T],
) error {
	if len(*opts.Params) < len(opts.Versions) {
		*opts.Params = make([]T, len(opts.Versions))
	} else if len(*opts.Params) > len(opts.Versions) {
		*opts.Params = (*opts.Params)[:len(opts.Versions)]
	}

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
			return err
		}

		cntr := start
		for rows.Next() {
			(*opts.Params)[cntr], err = pgx.RowToStructByName[T](rows)
			if err != nil {
				rows.Close()
				return err
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
			"DAL: Read hyperparams by version",
			"NumRows", end-start,
		)
	}
	return nil
}
