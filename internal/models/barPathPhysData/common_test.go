package barpathphysdata

import (
	"math"
	"testing"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbtest "code.barbellmath.net/barbell-math/smoothbrain-test"
)

func getBasicRawData() dal.CreatePhysicsDataParams {
	return dal.CreatePhysicsDataParams{
		Time: [][]types.Second{
			{0, 1, 2, 3, 4, 5, 6},
			{0, 1, 2, 3, 4, 5, 6},
		},
		Position: [][]types.Vec2[types.Meter, types.Meter]{{
			{X: 0 * 0, Y: 0 * 0},
			{X: 1 * 1, Y: 1 * 1},
			{X: 2 * 2, Y: 2 * 2},
			{X: 3 * 3, Y: 3 * 3},
			{X: 4 * 4, Y: 4 * 4},
			{X: 5 * 5, Y: 5 * 5},
			{X: 6 * 6, Y: 6 * 6},
		}, {
			{X: 0 * 0, Y: 0 * 0},
			{X: 1 * 1, Y: 1 * 1},
			{X: 2 * 2, Y: 2 * 2},
			{X: 3 * 3, Y: 3 * 3},
			{X: 4 * 4, Y: 4 * 4},
			{X: 5 * 5, Y: 5 * 5},
			{X: 6 * 6, Y: 6 * 6},
		}},
		Velocity: [][]types.Vec2[types.MeterPerSec, types.MeterPerSec]{
			[]types.Vec2[types.MeterPerSec, types.MeterPerSec]{},
			[]types.Vec2[types.MeterPerSec, types.MeterPerSec]{},
		},
		Acceleration: [][]types.Vec2[types.MeterPerSec2, types.MeterPerSec2]{
			[]types.Vec2[types.MeterPerSec2, types.MeterPerSec2]{},
			[]types.Vec2[types.MeterPerSec2, types.MeterPerSec2]{},
		},
		Jerk: [][]types.Vec2[types.MeterPerSec3, types.MeterPerSec3]{
			[]types.Vec2[types.MeterPerSec3, types.MeterPerSec3]{},
			[]types.Vec2[types.MeterPerSec3, types.MeterPerSec3]{},
		},
		Impulse: [][]types.Vec2[types.NewtonSec, types.NewtonSec]{
			[]types.Vec2[types.NewtonSec, types.NewtonSec]{},
			[]types.Vec2[types.NewtonSec, types.NewtonSec]{},
		},
		Force: [][]types.Vec2[types.Newton, types.Newton]{
			[]types.Vec2[types.Newton, types.Newton]{},
			[]types.Vec2[types.Newton, types.Newton]{},
		},
		Work:      [][]types.Joule{[]types.Joule{}, []types.Joule{}},
		Power:     [][]types.Watt{[]types.Watt{}, []types.Watt{}},
		RepSplits: [][]types.Split{[]types.Split{}, []types.Split{}},
		MinVel: [][]types.PointInTime[types.Second, types.MeterPerSec]{
			[]types.PointInTime[types.Second, types.MeterPerSec]{},
			[]types.PointInTime[types.Second, types.MeterPerSec]{},
		},
		MaxVel: [][]types.PointInTime[types.Second, types.MeterPerSec]{
			[]types.PointInTime[types.Second, types.MeterPerSec]{},
			[]types.PointInTime[types.Second, types.MeterPerSec]{},
		},
		MinAcc: [][]types.PointInTime[types.Second, types.MeterPerSec2]{
			[]types.PointInTime[types.Second, types.MeterPerSec2]{},
			[]types.PointInTime[types.Second, types.MeterPerSec2]{},
		},
		MaxAcc: [][]types.PointInTime[types.Second, types.MeterPerSec2]{
			[]types.PointInTime[types.Second, types.MeterPerSec2]{},
			[]types.PointInTime[types.Second, types.MeterPerSec2]{},
		},
		MinForce: [][]types.PointInTime[types.Second, types.Newton]{
			[]types.PointInTime[types.Second, types.Newton]{},
			[]types.PointInTime[types.Second, types.Newton]{},
		},
		MaxForce: [][]types.PointInTime[types.Second, types.Newton]{
			[]types.PointInTime[types.Second, types.Newton]{},
			[]types.PointInTime[types.Second, types.Newton]{},
		},
		MinImpulse: [][]types.PointInTime[types.Second, types.NewtonSec]{
			[]types.PointInTime[types.Second, types.NewtonSec]{},
			[]types.PointInTime[types.Second, types.NewtonSec]{},
		},
		MaxImpulse: [][]types.PointInTime[types.Second, types.NewtonSec]{
			[]types.PointInTime[types.Second, types.NewtonSec]{},
			[]types.PointInTime[types.Second, types.NewtonSec]{},
		},
		AvgWork: [][]types.Joule{[]types.Joule{}, []types.Joule{}},
		MinWork: [][]types.PointInTime[types.Second, types.Joule]{
			[]types.PointInTime[types.Second, types.Joule]{},
			[]types.PointInTime[types.Second, types.Joule]{},
		},
		MaxWork: [][]types.PointInTime[types.Second, types.Joule]{
			[]types.PointInTime[types.Second, types.Joule]{},
			[]types.PointInTime[types.Second, types.Joule]{},
		},
		AvgPower: [][]types.Watt{[]types.Watt{}, []types.Watt{}},
		MinPower: [][]types.PointInTime[types.Second, types.Watt]{
			[]types.PointInTime[types.Second, types.Watt]{},
			[]types.PointInTime[types.Second, types.Watt]{},
		},
		MaxPower: [][]types.PointInTime[types.Second, types.Watt]{
			[]types.PointInTime[types.Second, types.Watt]{},
			[]types.PointInTime[types.Second, types.Watt]{},
		},
	}
}

func TestTimeSeriesNotIncreasingErr(t *testing.T) {
	rawData := getBasicRawData()
	rawData.Time[0] = []types.Second{0, 1, 2, 3, 2, 5, 6}
	state := types.State{
		BarPathCalc: types.BarPathCalcConf{
			ApproxErr:     types.FourthOrder,
			MinNumSamples: 5,
			TimeDeltaEps:  1e-6,
		},
	}
	baseData := dal.BulkCreateTrainingLogsParams{
		Weight: 1,
		Reps:   1,
		Sets:   1,
	}
	err := Calc(&state, &baseData, &rawData, 0)
	sbtest.ContainsError(t, types.TimeSeriesDecreaseErr, err)
}

func TestTimeSeriesNotMonotonicErr(t *testing.T) {
	rawData := getBasicRawData()
	rawData.Time[0] = []types.Second{0, 1, 2, 4, 4, 5, 6}
	state := types.State{
		BarPathCalc: types.BarPathCalcConf{
			ApproxErr:     types.FourthOrder,
			MinNumSamples: 5,
			TimeDeltaEps:  1e-6,
		},
	}
	baseData := dal.BulkCreateTrainingLogsParams{
		Weight: 1,
		Reps:   1,
		Sets:   1,
	}
	err := Calc(&state, &baseData, &rawData, 0)
	sbtest.ContainsError(t, types.TimeSeriesNotMonotonicErr, err)
}

func TestInvalidApproxErrErr(t *testing.T) {
	rawData := getBasicRawData()
	state := types.State{
		BarPathCalc: types.BarPathCalcConf{
			ApproxErr:     types.ApproximationError(500),
			MinNumSamples: 5,
			TimeDeltaEps:  1e-6,
		},
	}
	baseData := dal.BulkCreateTrainingLogsParams{
		Weight: 1,
		Reps:   1,
		Sets:   1,
	}
	err := Calc(&state, &baseData, &rawData, 0)
	sbtest.ContainsError(t, types.ErrInvalidApproximationError, err)
}

func TestFractionalSets(t *testing.T) {
	rawData := getBasicRawData()
	state := types.State{
		BarPathCalc: types.BarPathCalcConf{
			ApproxErr:     types.SecondOrder,
			MinNumSamples: 5,
			TimeDeltaEps:  1e-6,
		},
	}
	baseData := dal.BulkCreateTrainingLogsParams{
		Weight: 1,
		Reps:   2,
		Sets:   1.5,
	}
	err := Calc(&state, &baseData, &rawData, 0)
	sbtest.Nil(t, err)
	err = Calc(&state, &baseData, &rawData, 1)
	sbtest.Nil(t, err)
	sbtest.Eq(t, len(rawData.RepSplits), 2)
	sbtest.Eq(t, len(rawData.RepSplits[0]), 2)
	sbtest.Eq(t, len(rawData.RepSplits[1]), 1)

	baseData = dal.BulkCreateTrainingLogsParams{
		Weight: 1,
		Reps:   1,
		Sets:   1.5,
	}
	err = Calc(&state, &baseData, &rawData, 0)
	sbtest.Nil(t, err)
	err = Calc(&state, &baseData, &rawData, 1)
	sbtest.Nil(t, err)
	sbtest.Eq(t, len(rawData.RepSplits), 2)
	sbtest.Eq(t, len(rawData.RepSplits[0]), 1)
	sbtest.Eq(t, len(rawData.RepSplits[1]), 1)
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
	state *types.State,
	baseData *dal.BulkCreateTrainingLogsParams,
) {
	rawData := dal.CreatePhysicsDataParams{
		Time: [][]types.Second{make([]types.Second, samples)},
		Position: [][]types.Vec2[types.Meter, types.Meter]{
			make([]types.Vec2[types.Meter, types.Meter], samples),
		},
		Velocity: [][]types.Vec2[types.MeterPerSec, types.MeterPerSec]{
			[]types.Vec2[types.MeterPerSec, types.MeterPerSec]{},
		},
		Acceleration: [][]types.Vec2[types.MeterPerSec2, types.MeterPerSec2]{
			[]types.Vec2[types.MeterPerSec2, types.MeterPerSec2]{},
		},
		Jerk: [][]types.Vec2[types.MeterPerSec3, types.MeterPerSec3]{
			[]types.Vec2[types.MeterPerSec3, types.MeterPerSec3]{},
		},
		Impulse: [][]types.Vec2[types.NewtonSec, types.NewtonSec]{
			[]types.Vec2[types.NewtonSec, types.NewtonSec]{},
		},
		Force: [][]types.Vec2[types.Newton, types.Newton]{
			[]types.Vec2[types.Newton, types.Newton]{},
		},
		Work:      [][]types.Joule{[]types.Joule{}},
		Power:     [][]types.Watt{[]types.Watt{}},
		RepSplits: [][]types.Split{[]types.Split{}},
		MinVel: [][]types.PointInTime[types.Second, types.MeterPerSec]{
			[]types.PointInTime[types.Second, types.MeterPerSec]{},
		},
		MaxVel: [][]types.PointInTime[types.Second, types.MeterPerSec]{
			[]types.PointInTime[types.Second, types.MeterPerSec]{},
		},
		MinAcc: [][]types.PointInTime[types.Second, types.MeterPerSec2]{
			[]types.PointInTime[types.Second, types.MeterPerSec2]{},
		},
		MaxAcc: [][]types.PointInTime[types.Second, types.MeterPerSec2]{
			[]types.PointInTime[types.Second, types.MeterPerSec2]{},
		},
		MinForce: [][]types.PointInTime[types.Second, types.Newton]{
			[]types.PointInTime[types.Second, types.Newton]{},
		},
		MaxForce: [][]types.PointInTime[types.Second, types.Newton]{
			[]types.PointInTime[types.Second, types.Newton]{},
		},
		MinImpulse: [][]types.PointInTime[types.Second, types.NewtonSec]{
			[]types.PointInTime[types.Second, types.NewtonSec]{},
		},
		MaxImpulse: [][]types.PointInTime[types.Second, types.NewtonSec]{
			[]types.PointInTime[types.Second, types.NewtonSec]{},
		},
		AvgWork: [][]types.Joule{[]types.Joule{}},
		MinWork: [][]types.PointInTime[types.Second, types.Joule]{
			[]types.PointInTime[types.Second, types.Joule]{},
		},
		MaxWork: [][]types.PointInTime[types.Second, types.Joule]{
			[]types.PointInTime[types.Second, types.Joule]{},
		},
		AvgPower: [][]types.Watt{[]types.Watt{}},
		MinPower: [][]types.PointInTime[types.Second, types.Watt]{
			[]types.PointInTime[types.Second, types.Watt]{},
		},
		MaxPower: [][]types.PointInTime[types.Second, types.Watt]{
			[]types.PointInTime[types.Second, types.Watt]{},
		},
	}
	for i := range samples {
		rawData.Time[0][i] = types.Second(i)
		rawData.Position[0][i] = types.Vec2[types.Meter, types.Meter]{
			X: types.Meter(f.f1(float64(i))),
			Y: types.Meter(f.f1(float64(i))),
		}
	}
	err := Calc(state, baseData, &rawData, 0)
	sbtest.Nil(t, err)

	edgeGap := 2
	if state.BarPathCalc.ApproxErr == types.FourthOrder {
		edgeGap = 3
	}
	for i := edgeGap; i < samples-edgeGap; i++ {
		d1Val := types.MeterPerSec(f.d1(float64(i)))
		d2Val := types.MeterPerSec2(f.d2(float64(i)))
		d3Val := types.MeterPerSec3(f.d3(float64(i)))
		sbtest.EqFloat(t, d1Val, rawData.Velocity[0][i].X, 1e-6)
		sbtest.EqFloat(t, d1Val, rawData.Velocity[0][i].Y, 1e-6)
		sbtest.EqFloat(t, d2Val, rawData.Acceleration[0][i].X, 1e-6)
		sbtest.EqFloat(t, d2Val, rawData.Acceleration[0][i].Y, 1e-6)
		sbtest.EqFloat(t, d3Val, rawData.Jerk[0][i].X, 1e-6)
		sbtest.EqFloat(t, d3Val, rawData.Jerk[0][i].Y, 1e-6)
	}
	for i := range edgeGap {
		d1Val := types.MeterPerSec(f.d1(float64(edgeGap)))
		d2Val := types.MeterPerSec2(f.d2(float64(edgeGap)))
		d3Val := types.MeterPerSec3(f.d3(float64(edgeGap)))
		sbtest.EqFloat(t, d1Val, rawData.Velocity[0][i].X, 1e-6)
		sbtest.EqFloat(t, d1Val, rawData.Velocity[0][i].Y, 1e-6)
		sbtest.EqFloat(t, d2Val, rawData.Acceleration[0][i].X, 1e-6)
		sbtest.EqFloat(t, d2Val, rawData.Acceleration[0][i].Y, 1e-6)
		sbtest.EqFloat(t, d3Val, rawData.Jerk[0][i].X, 1e-6)
		sbtest.EqFloat(t, d3Val, rawData.Jerk[0][i].Y, 1e-6)
	}
	for i := samples - 1; i >= samples-edgeGap; i-- {
		d1Val := types.MeterPerSec(f.d1(float64(samples - edgeGap - 1)))
		d2Val := types.MeterPerSec2(f.d2(float64(samples - edgeGap - 1)))
		d3Val := types.MeterPerSec3(f.d3(float64(samples - edgeGap - 1)))
		sbtest.EqFloat(t, d1Val, rawData.Velocity[0][i].X, 1e-6)
		sbtest.EqFloat(t, d1Val, rawData.Velocity[0][i].Y, 1e-6)
		sbtest.EqFloat(t, d2Val, rawData.Acceleration[0][i].X, 1e-6)
		sbtest.EqFloat(t, d2Val, rawData.Acceleration[0][i].Y, 1e-6)
		sbtest.EqFloat(t, d3Val, rawData.Jerk[0][i].X, 1e-6)
		sbtest.EqFloat(t, d3Val, rawData.Jerk[0][i].Y, 1e-6)
	}
	for i := range samples {
		// When m=1 f=a becuase F=ma
		sbtest.EqFloat(
			t,
			float64(rawData.Force[0][i].X),
			float64(rawData.Acceleration[0][i].X),
			1e-6,
		)
		sbtest.EqFloat(
			t,
			float64(rawData.Force[0][i].Y),
			float64(rawData.Acceleration[0][i].Y),
			1e-6,
		)
		// When m=1, impulse=vel because I=\int F dt=\int ma dt = mv
		sbtest.EqFloat(
			t,
			float64(rawData.Impulse[0][i].X),
			float64(rawData.Velocity[0][i].X),
			1e-6,
		)
		sbtest.EqFloat(
			t,
			float64(rawData.Impulse[0][i].Y),
			float64(rawData.Velocity[0][i].Y),
			1e-6,
		)
		sbtest.EqFloat(
			t,
			rawData.Power[0][i],
			types.Watt(
				(float64(rawData.Force[0][i].X)*float64(rawData.Velocity[0][i].X))+
					(float64(rawData.Force[0][i].Y)*float64(rawData.Velocity[0][i].Y)),
			),
			1e-6,
		)
		sbtest.EqFloat(
			t,
			rawData.Work[0][i],
			types.Joule(
				0.5*((rawData.Velocity[0][i].X*rawData.Velocity[0][i].X)+
					(rawData.Velocity[0][i].Y*rawData.Velocity[0][i].Y)),
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
		&types.State{
			BarPathCalc: types.BarPathCalcConf{
				ApproxErr:     types.SecondOrder,
				MinNumSamples: 10,
				TimeDeltaEps:  1e-6,
			},
		},
		&dal.BulkCreateTrainingLogsParams{
			Weight: 1,
			Reps:   1,
			Sets:   1,
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
		&types.State{
			BarPathCalc: types.BarPathCalcConf{
				ApproxErr:     types.FourthOrder,
				MinNumSamples: 10,
				TimeDeltaEps:  1e-6,
			},
		},
		&dal.BulkCreateTrainingLogsParams{
			Weight: 1,
			Reps:   1,
			Sets:   1,
		},
	)
}
