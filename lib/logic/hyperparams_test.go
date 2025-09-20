package logic

import (
	"context"
	"testing"

	"code.barbellmath.net/barbell-math/providentia/internal/db/migrations"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbtest "code.barbellmath.net/barbell-math/smoothbrain-test"
)

func TestHyperparams(t *testing.T) {
	t.Run("failingNoWrites", hyperparamsFailingNoWrites)
	t.Run("addGet", hyperparamsAddGet)
	t.Run("addDeleteGet", hyperparamsAddDeleteGet)
}

func hyperparamsFailingNoWrites(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)

	t.Run("invalidBarPathCalc", hyperparamsInvalidBarPathCalc(ctxt))
	t.Run("invalidBarPathTracker", hyperparamsInvalidBarPathTracker(ctxt))

	numExercises, err := ReadNumHyperparams(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, int64(len(migrations.HyperparamsSetupData)), numExercises)
}

func hyperparamsInvalidBarPathCalc(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateHyperparams(
			ctxt,
			types.BarPathCalcHyperparams{
				MinNumSamples:  1,
				TimeDeltaEps:   1,
				ApproxErr:      types.SecondOrder,
				NearZeroFilter: 1,
			},
		)
		sbtest.ContainsError(t, types.InvalidHyperparamsErr, err)
		sbtest.ContainsError(t, types.InvalidBarPathCalcErr, err)
		sbtest.ContainsError(t, types.InvalidMinNumSamplesErr, err)

		err = CreateHyperparams(
			ctxt,
			types.BarPathCalcHyperparams{
				MinNumSamples:  2,
				TimeDeltaEps:   0,
				ApproxErr:      types.SecondOrder,
				NearZeroFilter: 1,
			},
		)
		sbtest.ContainsError(t, types.InvalidHyperparamsErr, err)
		sbtest.ContainsError(t, types.InvalidBarPathCalcErr, err)
		sbtest.ContainsError(t, types.InvalidTimeDeltaEpsErr, err)

		err = CreateHyperparams(
			ctxt,
			types.BarPathCalcHyperparams{
				MinNumSamples:  2,
				TimeDeltaEps:   1,
				ApproxErr:      9999,
				NearZeroFilter: 1,
			},
		)
		sbtest.ContainsError(t, types.InvalidHyperparamsErr, err)
		sbtest.ContainsError(t, types.InvalidBarPathCalcErr, err)
		sbtest.ContainsError(t, types.ErrInvalidApproximationError, err)

		err = CreateHyperparams(
			ctxt,
			types.BarPathCalcHyperparams{
				MinNumSamples:  2,
				TimeDeltaEps:   1,
				ApproxErr:      types.SecondOrder,
				NearZeroFilter: -1,
			},
		)
		sbtest.ContainsError(t, types.InvalidHyperparamsErr, err)
		sbtest.ContainsError(t, types.InvalidBarPathCalcErr, err)
		sbtest.ContainsError(t, types.InvalidNearZeroFilterErr, err)
	}
}

func hyperparamsInvalidBarPathTracker(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		err := CreateHyperparams(
			ctxt,
			types.BarPathTrackerHyperparams{
				MinLength:   -1,
				MinFileSize: 1,
				MaxFileSize: 1,
			},
		)
		sbtest.ContainsError(t, types.InvalidHyperparamsErr, err)
		sbtest.ContainsError(t, types.InvalidBarPathTrackerErr, err)
		sbtest.ContainsError(t, types.InvalidMinLengthErr, err)

		err = CreateHyperparams(
			ctxt,
			types.BarPathTrackerHyperparams{
				MinLength:   1,
				MinFileSize: 2,
				MaxFileSize: 1,
			},
		)
		sbtest.ContainsError(t, types.InvalidHyperparamsErr, err)
		sbtest.ContainsError(t, types.InvalidBarPathTrackerErr, err)
		sbtest.ContainsError(t, types.InvalidMaxFileSizeErr, err)
	}
}

func hyperparamsAddGet(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)

	calcParams := []types.BarPathCalcHyperparams{{
		Version:        1,
		MinNumSamples:  2,
		TimeDeltaEps:   1,
		ApproxErr:      types.SecondOrder,
		NearZeroFilter: 1,
	}, {
		Version:        2,
		MinNumSamples:  2,
		TimeDeltaEps:   1,
		ApproxErr:      types.SecondOrder,
		NearZeroFilter: 1,
	}}
	trackParams := []types.BarPathTrackerHyperparams{{
		Version:     1,
		MinLength:   1,
		MinFileSize: 1,
		MaxFileSize: 2,
	}, {
		Version:     2,
		MinLength:   1,
		MinFileSize: 1,
		MaxFileSize: 2,
	}}

	err := CreateHyperparams(ctxt, calcParams...)
	sbtest.Nil(t, err)

	err = CreateHyperparams(ctxt, trackParams...)
	sbtest.Nil(t, err)

	res, err := ReadNumHyperparams(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 6, res)

	res, err = ReadNumHyperparamsFor[types.BarPathCalcHyperparams](ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 3, res)
	res, err = ReadNumHyperparamsFor[types.BarPathTrackerHyperparams](ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 3, res)

	params, err := ReadHyperparamsByVersionFor[types.BarPathCalcHyperparams](
		ctxt, 1,
	)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, len(params))
	sbtest.Eq(t, params[0], calcParams[0])

	params, err = ReadHyperparamsByVersionFor[types.BarPathCalcHyperparams](
		ctxt, 2,
	)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, len(params))
	sbtest.Eq(t, params[0], calcParams[1])

	params2, err := ReadHyperparamsByVersionFor[types.BarPathTrackerHyperparams](
		ctxt, 1,
	)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, len(params2))
	sbtest.Eq(t, params2[0], trackParams[0])

	params2, err = ReadHyperparamsByVersionFor[types.BarPathTrackerHyperparams](
		ctxt, 2,
	)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, len(params2))
	sbtest.Eq(t, params2[0], trackParams[1])
}

func hyperparamsAddDeleteGet(t *testing.T) {
	ctxt, cleanup := resetApp(context.Background())
	t.Cleanup(cleanup)

	calcParams := []types.BarPathCalcHyperparams{{
		Version:        1,
		MinNumSamples:  2,
		TimeDeltaEps:   1,
		ApproxErr:      types.SecondOrder,
		NearZeroFilter: 1,
	}, {
		Version:        2,
		MinNumSamples:  2,
		TimeDeltaEps:   1,
		ApproxErr:      types.SecondOrder,
		NearZeroFilter: 1,
	}}
	trackParams := []types.BarPathTrackerHyperparams{{
		Version:     1,
		MinLength:   1,
		MinFileSize: 1,
		MaxFileSize: 2,
	}, {
		Version:     2,
		MinLength:   1,
		MinFileSize: 1,
		MaxFileSize: 2,
	}}

	err := CreateHyperparams(ctxt, calcParams...)
	sbtest.Nil(t, err)

	err = CreateHyperparams(ctxt, trackParams...)
	sbtest.Nil(t, err)

	res, err := ReadNumHyperparams(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 6, res)

	res, err = ReadNumHyperparamsFor[types.BarPathCalcHyperparams](ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 3, res)
	res, err = ReadNumHyperparamsFor[types.BarPathTrackerHyperparams](ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 3, res)

	err = DeleteHyperparams[types.BarPathCalcHyperparams](ctxt, 1)
	sbtest.Nil(t, err)

	params, err := ReadHyperparamsByVersionFor[types.BarPathCalcHyperparams](
		ctxt, 1,
	)
	sbtest.ContainsError(t, types.CouldNotFindRequestedHyperparamsErr, err)
	res, err = ReadNumHyperparamsFor[types.BarPathCalcHyperparams](ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 2, res)

	params, err = ReadHyperparamsByVersionFor[types.BarPathCalcHyperparams](
		ctxt, 2,
	)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, len(params))
	sbtest.Eq(t, params[0], calcParams[1])

	err = DeleteHyperparams[types.BarPathTrackerHyperparams](ctxt, 1)
	sbtest.Nil(t, err)

	params2, err := ReadHyperparamsByVersionFor[types.BarPathTrackerHyperparams](
		ctxt, 1,
	)
	sbtest.ContainsError(t, types.CouldNotFindRequestedHyperparamsErr, err)
	res, err = ReadNumHyperparamsFor[types.BarPathCalcHyperparams](ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 2, res)

	params2, err = ReadHyperparamsByVersionFor[types.BarPathTrackerHyperparams](
		ctxt, 2,
	)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, len(params2))
	sbtest.Eq(t, params2[0], trackParams[1])

	err = DeleteHyperparams[types.BarPathCalcHyperparams](ctxt, 1)
	sbtest.ContainsError(t, types.CouldNotDeleteHyperparamsErr, err)
	sbtest.ContainsError(t, types.CouldNotFindRequestedHyperparamsErr, err)
}
