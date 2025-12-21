package tests

import (
	"context"
	"testing"

	"code.barbellmath.net/barbell-math/providentia/internal/dal/migrations"
	"code.barbellmath.net/barbell-math/providentia/lib/logic"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbtest "code.barbellmath.net/barbell-math/smoothbrain-test"
)

var (
	testingCalcHyperparams = types.BarPathCalcHyperparams{
		MinNumSamples:   5,
		TimeDeltaEps:    1,
		ApproxErr:       types.SecondOrder,
		NoiseFilter:     1,
		NearZeroFilter:  2,
		SmootherWeight1: 1,
		SmootherWeight2: 1,
		SmootherWeight3: 1,
		SmootherWeight4: 1,
		SmootherWeight5: 1,
	}
)

func TestPhysicsData(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	t.Run("errorCases", physicsDataErrorCases(ctxt))
	t.Run("constantNumReps", physicsDataConstantNumReps(ctxt))
	t.Run("lastSetLessReps", physicsDataLastSetLessReps(ctxt))
	t.Run("sparseRawData", physicsDataSparseRawData(ctxt))
}

func physicsDataErrorCases(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("noReps", physicsDataNoReps(ctxt))
		t.Run("notEnoughRawData", physicsDataNotEnoughRawData(ctxt))
	}
}

func physicsDataNoReps(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		exerciseData := types.ExerciseData{
			Weight: 1,
			Sets:   1,
			Reps:   0,
		}

		err := logic.CalcPhysicsData(
			ctxt,
			&migrations.BarPathCalcHyperparamsSetupData[0],
			&migrations.BarPathTrackerHyperparamsSetupData[0],
			&exerciseData,
			logic.BarPathTimeSeriesData(types.RawTimeSeriesData{}),
		)
		sbtest.ContainsError(
			t, types.PhysicsJobQueueErr, err,
			`Supplied exercise data must have at least 1 rep`,
		)
	}
}

func physicsDataNotEnoughRawData(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		exerciseData := types.ExerciseData{
			Weight: 1,
			Sets:   2,
			Reps:   1,
		}

		err := logic.CalcPhysicsData(
			ctxt,
			&migrations.BarPathCalcHyperparamsSetupData[0],
			&migrations.BarPathTrackerHyperparamsSetupData[0],
			&exerciseData,
			logic.BarPathTimeSeriesData(types.RawTimeSeriesData{}),
		)
		sbtest.ContainsError(
			t, types.PhysicsJobQueueErr, err,
			`The length of the raw data \(1\) must equal the ceiling of the number of sets \(2.000000\)`,
		)
	}
}

func physicsDataConstantNumReps(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		exerciseData := types.ExerciseData{
			Weight: 1,
			Sets:   2,
			Reps:   2,
		}
		rawData := logic.BarPathTimeSeriesData(types.RawTimeSeriesData{
			TimeData: []types.Second{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
			PositionData: []types.Vec2[types.Meter, types.Meter]{
				{X: 0, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 2},
				{X: 3, Y: 3},
				{X: 2, Y: 2}, {X: 1, Y: 1}, {X: 0, Y: 0},
				{X: 1, Y: 1}, {X: 2, Y: 2},
				{X: 3, Y: 3},
				{X: 2, Y: 2}, {X: 1, Y: 1}, {X: 0, Y: 0},
			},
		})

		err := logic.CalcPhysicsData(
			ctxt,
			&testingCalcHyperparams,
			&migrations.BarPathTrackerHyperparamsSetupData[0],
			&exerciseData,
			rawData, rawData,
		)
		sbtest.Nil(t, err)
		sbtest.True(t, exerciseData.PhysData[0].Present)
		sbtest.True(t, exerciseData.PhysData[1].Present)
		sbtest.SlicesMatch(
			t, exerciseData.PhysData[0].Value.RepSplits, []types.Split{
				{StartIdx: 0, EndIdx: 6},
				{StartIdx: 5, EndIdx: 13},
			},
		)
		sbtest.SlicesMatch(
			t, exerciseData.PhysData[1].Value.RepSplits, []types.Split{
				{StartIdx: 0, EndIdx: 6},
				{StartIdx: 5, EndIdx: 13},
			},
		)
	}
}

func physicsDataLastSetLessReps(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		exerciseData := types.ExerciseData{
			Weight: 1,
			Sets:   1.5,
			Reps:   2,
		}
		rawData := logic.BarPathTimeSeriesData(types.RawTimeSeriesData{
			TimeData: []types.Second{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
			PositionData: []types.Vec2[types.Meter, types.Meter]{
				{X: 0, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 2},
				{X: 3, Y: 3},
				{X: 2, Y: 2}, {X: 1, Y: 1}, {X: 0, Y: 0},
				{X: 1, Y: 1}, {X: 2, Y: 2},
				{X: 3, Y: 3},
				{X: 2, Y: 2}, {X: 1, Y: 1}, {X: 0, Y: 0},
			},
		})

		err := logic.CalcPhysicsData(
			ctxt,
			&testingCalcHyperparams,
			&migrations.BarPathTrackerHyperparamsSetupData[0],
			&exerciseData,
			rawData, rawData,
		)
		sbtest.Nil(t, err)
		sbtest.True(t, exerciseData.PhysData[0].Present)
		sbtest.True(t, exerciseData.PhysData[1].Present)
		sbtest.SlicesMatch(
			t, exerciseData.PhysData[0].Value.RepSplits, []types.Split{
				{StartIdx: 0, EndIdx: 6},
				{StartIdx: 5, EndIdx: 13},
			},
		)
		sbtest.SlicesMatch(
			t, exerciseData.PhysData[1].Value.RepSplits, []types.Split{
				{StartIdx: 0, EndIdx: 6},
			},
		)
	}
}

func physicsDataSparseRawData(ctxt context.Context) func(t *testing.T) {
	return func(t *testing.T) {
		exerciseData := types.ExerciseData{
			Weight: 1,
			Sets:   3,
			Reps:   2,
		}
		rawData := logic.BarPathTimeSeriesData(types.RawTimeSeriesData{
			TimeData: []types.Second{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
			PositionData: []types.Vec2[types.Meter, types.Meter]{
				{X: 0, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 2},
				{X: 3, Y: 3},
				{X: 2, Y: 2}, {X: 1, Y: 1}, {X: 0, Y: 0},
				{X: 1, Y: 1}, {X: 2, Y: 2},
				{X: 3, Y: 3},
				{X: 2, Y: 2}, {X: 1, Y: 1}, {X: 0, Y: 0},
			},
		})

		err := logic.CalcPhysicsData(
			ctxt,
			&testingCalcHyperparams,
			&migrations.BarPathTrackerHyperparamsSetupData[0],
			&exerciseData,
			rawData, types.BarPathVariant{}, rawData,
		)
		sbtest.Nil(t, err)
		sbtest.True(t, exerciseData.PhysData[0].Present)
		sbtest.False(t, exerciseData.PhysData[1].Present)
		sbtest.True(t, exerciseData.PhysData[2].Present)
		sbtest.SlicesMatch(
			t, exerciseData.PhysData[0].Value.RepSplits, []types.Split{
				{StartIdx: 0, EndIdx: 6},
				{StartIdx: 5, EndIdx: 13},
			},
		)
		sbtest.SlicesMatch(
			t, exerciseData.PhysData[1].Value.RepSplits, []types.Split{},
		)
		sbtest.SlicesMatch(
			t, exerciseData.PhysData[2].Value.RepSplits, []types.Split{
				{StartIdx: 0, EndIdx: 6},
				{StartIdx: 5, EndIdx: 13},
			},
		)
	}
}
