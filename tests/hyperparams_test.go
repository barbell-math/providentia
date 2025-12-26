package tests

import (
	"context"
	"math"
	"testing"

	"code.barbellmath.net/barbell-math/providentia/internal/dal/migrations"
	"code.barbellmath.net/barbell-math/providentia/lib/logic"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbcsv "code.barbellmath.net/barbell-math/smoothbrain-csv"
	sbtest "code.barbellmath.net/barbell-math/smoothbrain-test"
)

var (
	numDefaultHyperparams = int64(len(migrations.BarPathCalcHyperparamsSetupData) +
		len(migrations.BarPathTrackerHyperparamsSetupData))
)

func TestHyperparams(t *testing.T) {
	t.Run("failingNoWrites", hyperparamsFailingNoWrites)
	t.Run("defaults", hyperparamsDefaults)
	t.Run("duplicates", hyperparamsDuplicates)
	t.Run("createRead", hyperparamsCreateRead)
	t.Run("ensureRead", hyperparamsEnsureRead)
	t.Run("createFind", hyperparamsCreateFind)
	t.Run("createDelete", hyperparamsCreateDelete)
	t.Run("createCSVRead", hyperparamsCreateCSVRead)
	t.Run("ensureCSVRead", hyperparamsEnsureCSVRead)
}

func hyperparamsFailingNoWrites(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	t.Run("invalidBarPathCalc", hyperparamsInvalidBarPathCalc(ctxt))
	t.Run("invalidBarPathTracker", hyperparamsInvalidBarPathTracker(ctxt))

	n, err := logic.ReadNumHyperparams(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, numDefaultHyperparams, n)
}

func hyperparamsInvalidBarPathCalc(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := logic.CreateHyperparams(ctxt, types.BarPathCalcHyperparams{
			MinNumSamples:  1,
			TimeDeltaEps:   1,
			ApproxErr:      types.SecondOrder,
			NoiseFilter:    1,
			NearZeroFilter: 1,
		})
		sbtest.ContainsError(
			t, types.CouldNotCreateAllHyperparamsErr, err,
			`ERROR: COPY from stdin failed: Invalid bar path calc conf`,
			`Invalid min num samples`,
			`Must be >=2. Got: 1 \(SQLSTATE 57014\)`,
		)

		err = logic.CreateHyperparams(ctxt, types.BarPathCalcHyperparams{
			MinNumSamples:  2,
			TimeDeltaEps:   -1,
			ApproxErr:      types.SecondOrder,
			NoiseFilter:    1,
			NearZeroFilter: 1,
		})
		sbtest.ContainsError(
			t, types.CouldNotCreateAllHyperparamsErr, err,
			`ERROR: COPY from stdin failed: Invalid bar path calc conf`,
			`Invalid time delta eps`,
			`Must be >=0. Got: -1.000000 \(SQLSTATE 57014\)`,
		)

		err = logic.CreateHyperparams(ctxt, types.BarPathCalcHyperparams{
			MinNumSamples:  2,
			TimeDeltaEps:   1,
			ApproxErr:      types.ApproximationError(math.MaxInt32),
			NoiseFilter:    1,
			NearZeroFilter: 1,
		})
		sbtest.ContainsError(
			t, types.CouldNotCreateAllHyperparamsErr, err,
			`ERROR: COPY from stdin failed: Invalid bar path calc conf`,
			`not a valid ApproximationError, try \[SecondOrder, FourthOrder\] \(SQLSTATE 57014\)`,
		)

		err = logic.CreateHyperparams(ctxt, types.BarPathCalcHyperparams{
			MinNumSamples:  2,
			TimeDeltaEps:   1,
			ApproxErr:      types.SecondOrder,
			NoiseFilter:    0,
			NearZeroFilter: 1,
		})
		sbtest.ContainsError(
			t, types.CouldNotCreateAllHyperparamsErr, err,
			`ERROR: COPY from stdin failed: Invalid bar path calc conf`,
			`Invalid noise filter`,
			`Must be >0. Got: 0 \(SQLSTATE 57014\)`,
		)

		err = logic.CreateHyperparams(ctxt, types.BarPathCalcHyperparams{
			MinNumSamples:  2,
			TimeDeltaEps:   1,
			ApproxErr:      types.SecondOrder,
			NoiseFilter:    1,
			NearZeroFilter: -1,
		})
		sbtest.ContainsError(
			t, types.CouldNotCreateAllHyperparamsErr, err,
			`ERROR: COPY from stdin failed: Invalid bar path calc conf`,
			`Invalid near zero filter`,
			`Must be >0. Got: -1.000000 \(SQLSTATE 57014\)`,
		)
	}
}

func hyperparamsInvalidBarPathTracker(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := logic.CreateHyperparams(ctxt, types.BarPathTrackerHyperparams{
			MinLength:   -1,
			MinFileSize: 10,
			MaxFileSize: 10,
		})
		sbtest.ContainsError(
			t, types.CouldNotCreateAllHyperparamsErr, err,
			`ERROR: COPY from stdin failed: Invalid bar path tracker conf`,
			`Invalid min length`,
			`Must be >0. Got: -1.000000 \(SQLSTATE 57014\)`,
		)

		err = logic.CreateHyperparams(ctxt, types.BarPathTrackerHyperparams{
			MinLength:   1,
			MinFileSize: 10,
			MaxFileSize: 10,
		})
		sbtest.ContainsError(
			t, types.CouldNotCreateAllHyperparamsErr, err,
			`ERROR: COPY from stdin failed: Invalid bar path tracker conf`,
			`Invalid max file size`,
			`Max size \(10\) must be > min size \(10\) \(SQLSTATE 57014\)`,
		)
	}
}

func hyperparamsDefaults(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	res1, err := logic.ReadDefaultHyperparamsFor[types.BarPathCalcHyperparams](ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, res1, migrations.BarPathCalcHyperparamsSetupData[0])

	res2, err := logic.ReadDefaultHyperparamsFor[types.BarPathTrackerHyperparams](ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, res2, migrations.BarPathTrackerHyperparamsSetupData[0])

	n, err := logic.ReadNumHyperparams(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, numDefaultHyperparams, n)
}

func hyperparamsDuplicates(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	t.Run("modelIdVersion", hyperparamsDuplicateModelIdVersion(ctxt))
	t.Run("modelIdVersionParams", hyperparamsDuplicateModelIdVersionParams(ctxt))

	n, err := logic.ReadNumHyperparams(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, numDefaultHyperparams, n)
}

func hyperparamsDuplicateModelIdVersion(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := logic.CreateHyperparams(ctxt, types.BarPathCalcHyperparams{
			Version:        1,
			MinNumSamples:  2,
			TimeDeltaEps:   1,
			ApproxErr:      types.SecondOrder,
			NoiseFilter:    1,
			NearZeroFilter: 1,
		}, types.BarPathCalcHyperparams{
			Version:        1,
			MinNumSamples:  3,
			TimeDeltaEps:   2,
			ApproxErr:      types.FourthOrder,
			NoiseFilter:    2,
			NearZeroFilter: 2,
		})
		sbtest.ContainsError(
			t, types.CouldNotCreateAllHyperparamsErr, err,
			`ERROR: duplicate key value violates unique constraint "hyperparams_model_id_version_key" \(SQLSTATE 23505\)`,
		)
	}
}

func hyperparamsDuplicateModelIdVersionParams(
	ctxt context.Context,
) func(t *testing.T) {
	return func(t *testing.T) {
		err := logic.CreateHyperparams(ctxt, types.BarPathCalcHyperparams{
			Version:        1,
			MinNumSamples:  2,
			TimeDeltaEps:   1,
			ApproxErr:      types.SecondOrder,
			NoiseFilter:    1,
			NearZeroFilter: 1,
		}, types.BarPathCalcHyperparams{
			Version:        1,
			MinNumSamples:  2,
			TimeDeltaEps:   1,
			ApproxErr:      types.SecondOrder,
			NoiseFilter:    1,
			NearZeroFilter: 1,
		})
		sbtest.ContainsError(
			t, types.CouldNotCreateAllHyperparamsErr, err,
			`ERROR: duplicate key value violates unique constraint "hyperparams_model_id_version_key" \(SQLSTATE 23505\)`,
		)
	}
}

func hyperparamsCreateRead(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	params := []types.BarPathCalcHyperparams{
		{
			Version:        1,
			MinNumSamples:  2,
			TimeDeltaEps:   1,
			ApproxErr:      types.SecondOrder,
			NoiseFilter:    1,
			NearZeroFilter: 1,
		}, {
			Version:        2,
			MinNumSamples:  3,
			TimeDeltaEps:   2,
			ApproxErr:      types.FourthOrder,
			NoiseFilter:    2,
			NearZeroFilter: 2,
		},
	}
	err := logic.CreateHyperparams(ctxt, params...)
	sbtest.Nil(t, err)

	n, err := logic.ReadNumHyperparamsFor[types.BarPathCalcHyperparams](ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 3, n)
	n, err = logic.ReadNumHyperparamsFor[types.BarPathTrackerHyperparams](ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, n)

	readParams, err := logic.ReadHyperparamsByVersionFor[types.BarPathCalcHyperparams](
		ctxt, params[0].Version, params[1].Version,
	)
	sbtest.Nil(t, err)
	sbtest.SlicesMatch(t, params, readParams)

	readParams, err = logic.ReadHyperparamsByVersionFor[types.BarPathCalcHyperparams](ctxt, math.MaxInt32)
	sbtest.ContainsError(
		t, types.CouldNotReadAllHyperparamsErr, err,
		"Only read 0 entries out of batch of 1 requests",
	)

	err = logic.CreateHyperparams(ctxt, params...)
	sbtest.ContainsError(
		t, types.CouldNotCreateAllHyperparamsErr, err,
		`duplicate key value violates unique constraint "hyperparams_model_id_version_key" \(SQLSTATE 23505\)`,
	)

	n, err = logic.ReadNumHyperparams(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, numDefaultHyperparams+2, n)
}

func hyperparamsEnsureRead(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	params := []types.BarPathCalcHyperparams{
		{
			Version:        1,
			MinNumSamples:  2,
			TimeDeltaEps:   1,
			ApproxErr:      types.SecondOrder,
			NoiseFilter:    1,
			NearZeroFilter: 1,
		}, {
			Version:        2,
			MinNumSamples:  3,
			TimeDeltaEps:   2,
			ApproxErr:      types.FourthOrder,
			NoiseFilter:    2,
			NearZeroFilter: 2,
		},
	}
	err := logic.EnsureHyperparamsExist(ctxt, params...)
	sbtest.Nil(t, err)

	n, err := logic.ReadNumHyperparamsFor[types.BarPathCalcHyperparams](ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 3, n)
	n, err = logic.ReadNumHyperparamsFor[types.BarPathTrackerHyperparams](ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, n)

	readParams, err := logic.ReadHyperparamsByVersionFor[types.BarPathCalcHyperparams](
		ctxt, params[0].Version, params[1].Version,
	)
	sbtest.Nil(t, err)
	sbtest.SlicesMatch(t, params, readParams)

	readParams, err = logic.ReadHyperparamsByVersionFor[types.BarPathCalcHyperparams](ctxt, math.MaxInt32)
	sbtest.ContainsError(
		t, types.CouldNotReadAllHyperparamsErr, err,
		"Only read 0 entries out of batch of 1 requests",
	)

	err = logic.EnsureHyperparamsExist(ctxt, params...)
	sbtest.Nil(t, err)

	n, err = logic.ReadNumHyperparams(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, numDefaultHyperparams+2, n)
}

func hyperparamsCreateFind(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	params := []types.BarPathCalcHyperparams{
		{
			Version:        1,
			MinNumSamples:  2,
			TimeDeltaEps:   1,
			ApproxErr:      types.SecondOrder,
			NoiseFilter:    1,
			NearZeroFilter: 1,
		}, {
			Version:        2,
			MinNumSamples:  3,
			TimeDeltaEps:   2,
			ApproxErr:      types.FourthOrder,
			NoiseFilter:    2,
			NearZeroFilter: 2,
		},
	}
	err := logic.CreateHyperparams(ctxt, params...)
	sbtest.Nil(t, err)

	foundParams, err := logic.FindHyperparamsByVersionFor[types.BarPathCalcHyperparams](
		ctxt, params[0].Version, params[1].Version, math.MaxInt32,
	)
	sbtest.Nil(t, err)
	sbtest.SlicesMatch(t, []types.Found[types.BarPathCalcHyperparams]{{
		Found: true,
		Value: params[0],
	}, {
		Found: true,
		Value: params[1],
	}, {
		Found: false,
	}}, foundParams)

	n, err := logic.ReadNumHyperparams(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, numDefaultHyperparams+2, n)
}

func hyperparamsCreateDelete(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	params := []types.BarPathCalcHyperparams{
		{
			Version:        1,
			MinNumSamples:  2,
			TimeDeltaEps:   1,
			ApproxErr:      types.SecondOrder,
			NoiseFilter:    1,
			NearZeroFilter: 1,
		}, {
			Version:        2,
			MinNumSamples:  3,
			TimeDeltaEps:   2,
			ApproxErr:      types.FourthOrder,
			NoiseFilter:    2,
			NearZeroFilter: 2,
		}, {
			Version:        3,
			MinNumSamples:  4,
			TimeDeltaEps:   3,
			ApproxErr:      types.FourthOrder,
			NoiseFilter:    3,
			NearZeroFilter: 3,
		},
	}
	err := logic.CreateHyperparams(ctxt, params...)
	sbtest.Nil(t, err)

	n, err := logic.ReadNumHyperparams(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, numDefaultHyperparams+3, n)

	err = logic.DeleteHyperparams[types.BarPathCalcHyperparams](
		ctxt, params[0].Version,
	)
	sbtest.Nil(t, err)

	n, err = logic.ReadNumHyperparams(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, numDefaultHyperparams+2, n)

	readParams, err := logic.ReadHyperparamsByVersionFor[types.BarPathCalcHyperparams](
		ctxt, params[1].Version, params[2].Version,
	)
	sbtest.Nil(t, err)
	sbtest.SlicesMatch(t, params[1:], readParams)

	err = logic.DeleteHyperparams[types.BarPathCalcHyperparams](
		ctxt, params[0].Version,
	)
	sbtest.ContainsError(
		t, types.CouldNotDeleteAllHyperparamsErr, err,
		`Could not delete entry with version '1' \(Does it exist\?\)`,
	)

	n, err = logic.ReadNumHyperparams(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, numDefaultHyperparams+2, n)
}

func hyperparamsCreateCSVRead(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	err := logic.CreateHyperparamsFromCSV[types.BarPathCalcHyperparams](
		ctxt, &sbcsv.Opts{}, "./testData/hyperparamData/1.barPathCalc.csv",
	)
	sbtest.Nil(t, err)

	n, err := logic.ReadNumHyperparams(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, numDefaultHyperparams+2, n)

	client, err := logic.ReadHyperparamsByVersionFor[types.BarPathCalcHyperparams](
		ctxt, 1, 2,
	)
	sbtest.Nil(t, err)
	sbtest.SlicesMatch(t, client, []types.BarPathCalcHyperparams{
		{
			Version:         1,
			MinNumSamples:   2,
			TimeDeltaEps:    1,
			ApproxErr:       types.SecondOrder,
			NoiseFilter:     6,
			NearZeroFilter:  1,
			SmootherWeight1: 0.1,
			SmootherWeight2: 0.2,
			SmootherWeight3: 0.3,
			SmootherWeight4: 0.4,
			SmootherWeight5: 0.5,
		}, {
			Version:         2,
			MinNumSamples:   2,
			TimeDeltaEps:    1,
			ApproxErr:       types.FourthOrder,
			NoiseFilter:     6,
			NearZeroFilter:  1,
			SmootherWeight1: 0.1,
			SmootherWeight2: 0.2,
			SmootherWeight3: 0.3,
			SmootherWeight4: 0.4,
			SmootherWeight5: 0.5,
		},
	})

	_, err = logic.ReadHyperparamsByVersionFor[types.BarPathCalcHyperparams](ctxt, 3)
	sbtest.ContainsError(t, types.CouldNotReadAllHyperparamsErr, err)

	err = logic.CreateHyperparamsFromCSV[types.BarPathCalcHyperparams](
		ctxt, &sbcsv.Opts{}, "./testData/hyperparamData/1.barPathCalc.csv",
	)
	sbtest.ContainsError(t, types.CSVLoaderJobQueueErr, err)
	sbtest.ContainsError(
		t, types.CouldNotCreateAllHyperparamsErr, err,
		`duplicate key value violates unique constraint "hyperparams_model_id_version_key" \(SQLSTATE 23505\)`,
	)

	n, err = logic.ReadNumHyperparams(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, numDefaultHyperparams+2, n)
}

func hyperparamsEnsureCSVRead(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	err := logic.EnsureHyperparamsExistFromCSV[types.BarPathCalcHyperparams](
		ctxt, &sbcsv.Opts{}, "./testData/hyperparamData/1.barPathCalc.csv",
	)
	sbtest.Nil(t, err)

	n, err := logic.ReadNumHyperparams(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, numDefaultHyperparams+2, n)

	client, err := logic.ReadHyperparamsByVersionFor[types.BarPathCalcHyperparams](
		ctxt, 1, 2,
	)
	sbtest.Nil(t, err)
	sbtest.SlicesMatch(t, client, []types.BarPathCalcHyperparams{
		{
			Version:         1,
			MinNumSamples:   2,
			TimeDeltaEps:    1,
			ApproxErr:       types.SecondOrder,
			NoiseFilter:     6,
			NearZeroFilter:  1,
			SmootherWeight1: 0.1,
			SmootherWeight2: 0.2,
			SmootherWeight3: 0.3,
			SmootherWeight4: 0.4,
			SmootherWeight5: 0.5,
		}, {
			Version:         2,
			MinNumSamples:   2,
			TimeDeltaEps:    1,
			ApproxErr:       types.FourthOrder,
			NoiseFilter:     6,
			NearZeroFilter:  1,
			SmootherWeight1: 0.1,
			SmootherWeight2: 0.2,
			SmootherWeight3: 0.3,
			SmootherWeight4: 0.4,
			SmootherWeight5: 0.5,
		},
	})

	_, err = logic.ReadHyperparamsByVersionFor[types.BarPathCalcHyperparams](ctxt, 3)
	sbtest.ContainsError(t, types.CouldNotReadAllHyperparamsErr, err)

	err = logic.EnsureHyperparamsExistFromCSV[types.BarPathCalcHyperparams](
		ctxt, &sbcsv.Opts{}, "./testData/hyperparamData/1.barPathCalc.csv",
	)
	sbtest.Nil(t, err)

	n, err = logic.ReadNumHyperparams(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, numDefaultHyperparams+2, n)
}
