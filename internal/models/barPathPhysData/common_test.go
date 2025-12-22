package barpathphysdata

import (
	"math"
	"testing"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbtest "code.barbellmath.net/barbell-math/smoothbrain-test"
)

func getBasicRawData() types.PhysicsData {
	return types.PhysicsData{
		Time: []types.Second{0, 1, 2, 3, 4, 5, 6},
		Position: []types.Vec2[types.Meter, types.Meter]{
			{X: 0 * 0, Y: 0 * 0},
			{X: 1 * 1, Y: 1 * 1},
			{X: 2 * 2, Y: 2 * 2},
			{X: 3 * 3, Y: 3 * 3},
			{X: 4 * 4, Y: 4 * 4},
			{X: 5 * 5, Y: 5 * 5},
			{X: 6 * 6, Y: 6 * 6},
		},
	}
}

func TestNotEnoughSamples(t *testing.T) {
	rawData := getBasicRawData()
	params := types.BarPathCalcHyperparams{
		ApproxErr:     types.FourthOrder,
		MinNumSamples: uint64(len(rawData.Time) + 1),
		TimeDeltaEps:  1e-6,
	}
	err := Calc(&rawData, &params, 1, 1)
	sbtest.ContainsError(
		t, types.InvalidRawDataLenErr, err,
		`The minimum number of samples \(8\) was not provided, got 7 samples`,
	)
}

func TestPositionDataNotAllocated(t *testing.T) {
	rawData := getBasicRawData()
	rawData.Position = []types.Vec2[types.Meter, types.Meter]{}
	params := types.BarPathCalcHyperparams{
		ApproxErr:     types.FourthOrder,
		MinNumSamples: 5,
		TimeDeltaEps:  1e-6,
	}
	err := Calc(&rawData, &params, 1, 1)
	sbtest.ContainsError(
		t, types.InvalidRawDataLenErr, err,
		`Expected position slice of len 7, got len 0`,
	)
}

func TestInvalidExpNumReps(t *testing.T) {
	rawData := getBasicRawData()
	params := types.BarPathCalcHyperparams{
		ApproxErr:     types.FourthOrder,
		MinNumSamples: 5,
		TimeDeltaEps:  1e-6,
	}
	err := Calc(&rawData, &params, 1, 0)
	sbtest.ContainsError(
		t, types.InvalidExpNumRepsErr, err,
		`Must be >=0. Got: 0`,
	)
}

func TestTimeSeriesNotIncreasingErr(t *testing.T) {
	rawData := getBasicRawData()
	rawData.Time = []types.Second{0, 1, 2, 3, 2, 5, 6}
	params := types.BarPathCalcHyperparams{
		ApproxErr:     types.FourthOrder,
		MinNumSamples: 5,
		TimeDeltaEps:  1e-6,
	}
	err := Calc(&rawData, &params, 1, 1)
	sbtest.ContainsError(
		t, types.TimeSeriesDecreaseErr, err,
		`Time samples must be increasing`,
	)
}

func TestTimeSeriesNotMonotonicErr(t *testing.T) {
	rawData := getBasicRawData()
	rawData.Time = []types.Second{0, 1, 2, 4, 4, 5, 6}
	params := types.BarPathCalcHyperparams{
		ApproxErr:     types.FourthOrder,
		MinNumSamples: 5,
		TimeDeltaEps:  1e-6,
		NoiseFilter:   3,
	}
	err := Calc(&rawData, &params, 1, 1)
	sbtest.ContainsError(
		t, types.TimeSeriesNotMonotonicErr, err,
		`Adjacent time samples must all have the same delta \(within 0.000001 variance\)`,
	)
}

func TestInvalidApproxErrErr(t *testing.T) {
	rawData := getBasicRawData()
	params := types.BarPathCalcHyperparams{
		ApproxErr:     types.ApproximationError(math.MaxInt32),
		MinNumSamples: 5,
		TimeDeltaEps:  1e-6,
	}
	err := Calc(&rawData, &params, 1, 1)
	sbtest.ContainsError(
		t, types.ErrInvalidApproximationError, err,
		`not a valid ApproximationError, try \[SecondOrder, FourthOrder\]`,
	)
}

type funcVals struct {
	f1 func(x float64) float64
	d1 func(x float64) float64
	d2 func(x float64) float64
	d3 func(x float64) float64
}

func testForAccuracy(
	t *testing.T,
	samples int,
	f funcVals,
	params *types.BarPathCalcHyperparams,
) {
	rawData := types.PhysicsData{
		Time:     make([]types.Second, samples),
		Position: make([]types.Vec2[types.Meter, types.Meter], samples),
	}
	for i := range samples {
		rawData.Time[i] = types.Second(i)
		rawData.Position[i] = types.Vec2[types.Meter, types.Meter]{
			X: types.Meter(f.f1(float64(i))),
			Y: types.Meter(f.f1(float64(i))),
		}
	}
	err := Calc(&rawData, params, 1, 1)
	sbtest.Nil(t, err)

	edgeGap := 2
	if params.ApproxErr == types.FourthOrder {
		edgeGap = 3
	}
	for i := edgeGap; i < samples-edgeGap; i++ {
		d1Val := types.MeterPerSec(f.d1(float64(i)))
		d2Val := types.MeterPerSec2(f.d2(float64(i)))
		d3Val := types.MeterPerSec3(f.d3(float64(i)))
		sbtest.EqFloat(t, d1Val, rawData.Velocity[i].X, 1e-6)
		sbtest.EqFloat(t, d1Val, rawData.Velocity[i].Y, 1e-6)
		sbtest.EqFloat(t, d2Val, rawData.Acceleration[i].X, 1e-6)
		sbtest.EqFloat(t, d2Val, rawData.Acceleration[i].Y, 1e-6)
		sbtest.EqFloat(t, d3Val, rawData.Jerk[i].X, 1e-6)
		sbtest.EqFloat(t, d3Val, rawData.Jerk[i].Y, 1e-6)
	}
	for i := range edgeGap {
		d1Val := types.MeterPerSec(f.d1(float64(edgeGap)))
		d2Val := types.MeterPerSec2(f.d2(float64(edgeGap)))
		d3Val := types.MeterPerSec3(f.d3(float64(edgeGap)))
		sbtest.EqFloat(t, d1Val, rawData.Velocity[i].X, 1e-6)
		sbtest.EqFloat(t, d1Val, rawData.Velocity[i].Y, 1e-6)
		sbtest.EqFloat(t, d2Val, rawData.Acceleration[i].X, 1e-6)
		sbtest.EqFloat(t, d2Val, rawData.Acceleration[i].Y, 1e-6)
		sbtest.EqFloat(t, d3Val, rawData.Jerk[i].X, 1e-6)
		sbtest.EqFloat(t, d3Val, rawData.Jerk[i].Y, 1e-6)
	}
	for i := samples - 1; i >= samples-edgeGap; i-- {
		d1Val := types.MeterPerSec(f.d1(float64(samples - edgeGap - 1)))
		d2Val := types.MeterPerSec2(f.d2(float64(samples - edgeGap - 1)))
		d3Val := types.MeterPerSec3(f.d3(float64(samples - edgeGap - 1)))
		sbtest.EqFloat(t, d1Val, rawData.Velocity[i].X, 1e-6)
		sbtest.EqFloat(t, d1Val, rawData.Velocity[i].Y, 1e-6)
		sbtest.EqFloat(t, d2Val, rawData.Acceleration[i].X, 1e-6)
		sbtest.EqFloat(t, d2Val, rawData.Acceleration[i].Y, 1e-6)
		sbtest.EqFloat(t, d3Val, rawData.Jerk[i].X, 1e-6)
		sbtest.EqFloat(t, d3Val, rawData.Jerk[i].Y, 1e-6)
	}
	for i := range samples {
		// When m=1 f=a becuase F=ma
		sbtest.EqFloat(
			t,
			float64(rawData.Force[i].X),
			float64(rawData.Acceleration[i].X),
			1e-6,
		)
		sbtest.EqFloat(
			t,
			float64(rawData.Force[i].Y),
			float64(rawData.Acceleration[i].Y),
			1e-6,
		)
		// When m=1, impulse=vel because I=\int F dt=\int ma dt = mv
		sbtest.EqFloat(
			t,
			float64(rawData.Impulse[i].X),
			float64(rawData.Velocity[i].X),
			1e-6,
		)
		sbtest.EqFloat(
			t,
			float64(rawData.Impulse[i].Y),
			float64(rawData.Velocity[i].Y),
			1e-6,
		)
		sbtest.EqFloat(
			t,
			rawData.Power[i],
			types.Watt(
				(float64(rawData.Force[i].X)*float64(rawData.Velocity[i].X))+
					(float64(rawData.Force[i].Y)*float64(rawData.Velocity[i].Y)),
			),
			1e-6,
		)
		sbtest.EqFloat(
			t,
			rawData.Work[i],
			types.Joule(
				0.5*((rawData.Velocity[i].X*rawData.Velocity[i].X)+
					(rawData.Velocity[i].Y*rawData.Velocity[i].Y)),
			),
			1e-6,
		)
	}
}

func TestQuadPolynomialSecondOrderAccuracy(t *testing.T) {
	// Second order error approx can exactly represent polynomials up to a power
	// of 2, so that is what we test with
	testForAccuracy(
		t,
		100,
		funcVals{
			f1: func(x float64) float64 { return math.Pow(x, 2) },
			d1: func(x float64) float64 { return 2 * math.Pow(x, 1) },
			d2: func(x float64) float64 { return 2 },
			d3: func(x float64) float64 { return 0 },
		},
		&types.BarPathCalcHyperparams{
			ApproxErr:     types.SecondOrder,
			MinNumSamples: 10,
			TimeDeltaEps:  1e-6,
			NoiseFilter:   3,
		},
	)
}

func TestQuadPolynomialFourthOrderAccuracy(t *testing.T) {
	// Fourth order error approx can exactly represent polynomials up to a power
	// of 4, so that is what we test with
	testForAccuracy(
		t,
		100,
		funcVals{
			f1: func(x float64) float64 { return math.Pow(x, 4) },
			d1: func(x float64) float64 { return 4 * math.Pow(x, 3) },
			d2: func(x float64) float64 { return 12 * math.Pow(x, 2) },
			d3: func(x float64) float64 { return 24 * math.Pow(x, 1) },
		},
		&types.BarPathCalcHyperparams{
			ApproxErr:     types.FourthOrder,
			MinNumSamples: 10,
			TimeDeltaEps:  1e-6,
			NoiseFilter:   3,
		},
	)
}

func TestRepsExtendToStartAndEnd(t *testing.T) {
	rawData := types.PhysicsData{
		Time: []types.Second{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
		Position: []types.Vec2[types.Meter, types.Meter]{
			{X: 0, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 2},
			{X: 3, Y: 3},
			{X: 2, Y: 2}, {X: 1, Y: 1}, {X: 0, Y: 0},
			{X: 1, Y: 1}, {X: 2, Y: 2},
			{X: 3, Y: 3},
			{X: 2, Y: 2}, {X: 1, Y: 1}, {X: 0, Y: 0},
		},
	}
	params := types.BarPathCalcHyperparams{
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
	err := Calc(&rawData, &params, 1, 2)
	sbtest.Nil(t, err)
	sbtest.SlicesMatch(t, rawData.RepSplits, []types.Split{
		{StartIdx: 0, EndIdx: 6},
		{StartIdx: 5, EndIdx: 13},
	})
}
