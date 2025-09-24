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
	t.Run("createRead", hyperparamsCreateRead)
	t.Run("ensureRead", hyperparamsEnsureRead)
	t.Run("createFind", hyperparamsCreateFind)
	t.Run("addDeleteRead", hyperparamsCreateDeleteRead)
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

func hyperparamsCreateRead(t *testing.T) {
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

func hyperparamsEnsureRead(t *testing.T) {
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

	err := EnsureHyperparamsExist(ctxt, calcParams...)
	sbtest.Nil(t, err)

	err = EnsureHyperparamsExist(ctxt, trackParams...)
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

	err = EnsureHyperparamsExist(ctxt, calcParams...)
	sbtest.Nil(t, err)

	err = EnsureHyperparamsExist(ctxt, trackParams...)
	sbtest.Nil(t, err)

	res, err = ReadNumHyperparams(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 6, res)

	err = EnsureHyperparamsExist(ctxt, types.BarPathCalcHyperparams{
		Version:        1,
		MinNumSamples:  2,
		TimeDeltaEps:   3,
		ApproxErr:      types.FourthOrder,
		NearZeroFilter: 4,
	})
	sbtest.ContainsError(t, types.CouldNotAddNumHyperparamsErr, err)
}

func hyperparamsCreateFind(t *testing.T) {
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

	params, err := FindHyperparamsByVersionFor[types.BarPathCalcHyperparams](
		ctxt, 1,
	)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, len(params))
	sbtest.True(t, params[0].Found)
	sbtest.Eq(t, params[0].Value, calcParams[0])

	params, err = FindHyperparamsByVersionFor[types.BarPathCalcHyperparams](
		ctxt, 2,
	)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, len(params))
	sbtest.True(t, params[0].Found)
	sbtest.Eq(t, params[0].Value, calcParams[1])

	params2, err := FindHyperparamsByVersionFor[types.BarPathTrackerHyperparams](
		ctxt, 1,
	)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, len(params2))
	sbtest.True(t, params[0].Found)
	sbtest.Eq(t, params2[0].Value, trackParams[0])

	params2, err = FindHyperparamsByVersionFor[types.BarPathTrackerHyperparams](
		ctxt, 2,
	)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, len(params2))
	sbtest.True(t, params[0].Found)
	sbtest.Eq(t, params2[0].Value, trackParams[1])

	versions := []int32{}
	for i := range len(calcParams) {
		versions = append(versions, calcParams[i].Version, 9999)
	}
	res2, err := FindHyperparamsByVersionFor[types.BarPathCalcHyperparams](
		ctxt, versions...,
	)
	sbtest.Nil(t, err)
	for i := range len(calcParams) {
		if i%2 == 0 {
			sbtest.True(t, res2[i].Found)
			sbtest.Eq(t, calcParams[i/2], res2[i].Value)
		} else {
			sbtest.False(t, res2[i].Found)
		}
	}
}

func hyperparamsCreateDeleteRead(t *testing.T) {
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
