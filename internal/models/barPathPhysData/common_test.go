package barpathphysdata

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"testing"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbtest "code.barbellmath.net/barbell-math/smoothbrain-test"
)

func TestTimeSeriesNotIncreasingErr(t *testing.T) {
	rawData := dal.CreatePhysicsDataParams{
		Time: [][]types.Second{{0, 1, 2, 3, 2, 5, 6}},
		Position: [][]types.Vec2[types.Meter]{{
			{X: 0 * 0, Y: 0 * 0},
			{X: 1 * 1, Y: 1 * 1},
			{X: 2 * 2, Y: 2 * 2},
			{X: 3 * 3, Y: 3 * 3},
			{X: 4 * 4, Y: 4 * 4},
			{X: 5 * 5, Y: 5 * 5},
			{X: 6 * 6, Y: 6 * 6},
		}},
		Velocity: [][]types.Vec2[types.MeterPerSec]{
			make([]types.Vec2[types.MeterPerSec], 7),
		},
		Acceleration: [][]types.Vec2[types.MeterPerSec2]{
			make([]types.Vec2[types.MeterPerSec2], 7),
		},
		Jerk: [][]types.Vec2[types.MeterPerSec3]{
			make([]types.Vec2[types.MeterPerSec3], 7),
		},
		Impulse: [][]types.Vec2[types.NewtonSec]{
			make([]types.Vec2[types.NewtonSec], 7),
		},
		Force: [][]types.Vec2[types.Newton]{
			make([]types.Vec2[types.Newton], 7),
		},
		Work:  [][]types.Joule{make([]types.Joule, 7)},
		Power: [][]types.Watt{make([]types.Watt, 7)},
		RepSplits: [][]types.Split[types.Second]{
			make([]types.Split[types.Second], 7),
		},
	}
	state := types.State{
		BarPathCalc: types.BarPathCalcConf{
			ApproxErr: types.FourthOrder,
		},
		PhysicsData: types.PhysicsDataConf{
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
	rawData := dal.CreatePhysicsDataParams{
		Time: [][]types.Second{{0, 1, 2, 4, 4, 5, 6}},
		Position: [][]types.Vec2[types.Meter]{{
			{X: 0 * 0, Y: 0 * 0},
			{X: 1 * 1, Y: 1 * 1},
			{X: 2 * 2, Y: 2 * 2},
			{X: 3 * 3, Y: 3 * 3},
			{X: 4 * 4, Y: 4 * 4},
			{X: 5 * 5, Y: 5 * 5},
			{X: 6 * 6, Y: 6 * 6},
		}},
		Velocity: [][]types.Vec2[types.MeterPerSec]{
			make([]types.Vec2[types.MeterPerSec], 7),
		},
		Acceleration: [][]types.Vec2[types.MeterPerSec2]{
			make([]types.Vec2[types.MeterPerSec2], 7),
		},
		Jerk: [][]types.Vec2[types.MeterPerSec3]{
			make([]types.Vec2[types.MeterPerSec3], 7),
		},
		Impulse: [][]types.Vec2[types.NewtonSec]{
			make([]types.Vec2[types.NewtonSec], 7),
		},
		Force: [][]types.Vec2[types.Newton]{
			make([]types.Vec2[types.Newton], 7),
		},
		Work:  [][]types.Joule{make([]types.Joule, 7)},
		Power: [][]types.Watt{make([]types.Watt, 7)},
		RepSplits: [][]types.Split[types.Second]{
			make([]types.Split[types.Second], 7),
		},
	}
	state := types.State{
		BarPathCalc: types.BarPathCalcConf{
			ApproxErr: types.FourthOrder,
		},
		PhysicsData: types.PhysicsDataConf{
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
	rawData := dal.CreatePhysicsDataParams{
		Time: [][]types.Second{{0, 1, 2, 3, 4, 5, 6}},
		Position: [][]types.Vec2[types.Meter]{{
			{X: 0 * 0, Y: 0 * 0},
			{X: 1 * 1, Y: 1 * 1},
			{X: 2 * 2, Y: 2 * 2},
			{X: 3 * 3, Y: 3 * 3},
			{X: 4 * 4, Y: 4 * 4},
			{X: 5 * 5, Y: 5 * 5},
			{X: 6 * 6, Y: 6 * 6},
		}},
		Velocity: [][]types.Vec2[types.MeterPerSec]{
			make([]types.Vec2[types.MeterPerSec], 7),
		},
		Acceleration: [][]types.Vec2[types.MeterPerSec2]{
			make([]types.Vec2[types.MeterPerSec2], 7),
		},
		Jerk: [][]types.Vec2[types.MeterPerSec3]{
			make([]types.Vec2[types.MeterPerSec3], 7),
		},
		Impulse: [][]types.Vec2[types.NewtonSec]{
			make([]types.Vec2[types.NewtonSec], 7),
		},
		Force: [][]types.Vec2[types.Newton]{
			make([]types.Vec2[types.Newton], 7),
		},
		Work:  [][]types.Joule{make([]types.Joule, 7)},
		Power: [][]types.Watt{make([]types.Watt, 7)},
		RepSplits: [][]types.Split[types.Second]{
			make([]types.Split[types.Second], 7),
		},
	}
	state := types.State{
		BarPathCalc: types.BarPathCalcConf{
			ApproxErr: types.ApproximationError(500),
		},
		PhysicsData: types.PhysicsDataConf{
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

func testForAccuracy(
	t *testing.T,
	f1 func(x float64) float64,
	d1 func(x float64) float64,
	d2 func(x float64) float64,
	d3 func(x float64) float64,
	samples int,
	state *types.State,
	baseData *dal.BulkCreateTrainingLogsParams,
) {
	rawData := dal.CreatePhysicsDataParams{
		Time:         [][]types.Second{make([]types.Second, samples)},
		Position:     [][]types.Vec2[types.Meter]{make([]types.Vec2[types.Meter], samples)},
		Velocity:     [][]types.Vec2[types.MeterPerSec]{make([]types.Vec2[types.MeterPerSec], samples)},
		Acceleration: [][]types.Vec2[types.MeterPerSec2]{make([]types.Vec2[types.MeterPerSec2], samples)},
		Jerk:         [][]types.Vec2[types.MeterPerSec3]{make([]types.Vec2[types.MeterPerSec3], samples)},
		Impulse:      [][]types.Vec2[types.NewtonSec]{make([]types.Vec2[types.NewtonSec], samples)},
		Force:        [][]types.Vec2[types.Newton]{make([]types.Vec2[types.Newton], samples)},
		Work:         [][]types.Joule{make([]types.Joule, samples)},
		Power:        [][]types.Watt{make([]types.Watt, samples)},
		RepSplits:    [][]types.Split[types.Second]{make([]types.Split[types.Second], samples)},
	}
	for i := range samples {
		rawData.Time[0][i] = types.Second(i)
		rawData.Position[0][i] = types.Vec2[types.Meter]{
			X: types.Meter(f1(float64(i))),
			Y: types.Meter(f1(float64(i))),
		}
	}
	err := Calc(state, baseData, &rawData, 0)
	sbtest.Nil(t, err)

	edgeGap := 2
	if state.BarPathCalc.ApproxErr == types.FourthOrder {
		edgeGap = 3
	}
	for i := edgeGap; i < samples-edgeGap; i++ {
		d1Val := types.MeterPerSec(d1(float64(i)))
		d2Val := types.MeterPerSec2(d2(float64(i)))
		d3Val := types.MeterPerSec3(d3(float64(i)))
		sbtest.EqFloat(t, d1Val, rawData.Velocity[0][i].X, 1e-6)
		sbtest.EqFloat(t, d1Val, rawData.Velocity[0][i].Y, 1e-6)
		sbtest.EqFloat(t, d2Val, rawData.Acceleration[0][i].X, 1e-6)
		sbtest.EqFloat(t, d2Val, rawData.Acceleration[0][i].Y, 1e-6)
		sbtest.EqFloat(t, d3Val, rawData.Jerk[0][i].X, 1e-6)
		sbtest.EqFloat(t, d3Val, rawData.Jerk[0][i].Y, 1e-6)
	}
	for i := range edgeGap {
		d1Val := types.MeterPerSec(d1(float64(edgeGap)))
		d2Val := types.MeterPerSec2(d2(float64(edgeGap)))
		d3Val := types.MeterPerSec3(d3(float64(edgeGap)))
		sbtest.EqFloat(t, d1Val, rawData.Velocity[0][i].X, 1e-6)
		sbtest.EqFloat(t, d1Val, rawData.Velocity[0][i].Y, 1e-6)
		sbtest.EqFloat(t, d2Val, rawData.Acceleration[0][i].X, 1e-6)
		sbtest.EqFloat(t, d2Val, rawData.Acceleration[0][i].Y, 1e-6)
		sbtest.EqFloat(t, d3Val, rawData.Jerk[0][i].X, 1e-6)
		sbtest.EqFloat(t, d3Val, rawData.Jerk[0][i].Y, 1e-6)
	}
	for i := samples - 1; i >= samples-edgeGap; i-- {
		d1Val := types.MeterPerSec(d1(float64(samples - edgeGap - 1)))
		d2Val := types.MeterPerSec2(d2(float64(samples - edgeGap - 1)))
		d3Val := types.MeterPerSec3(d3(float64(samples - edgeGap - 1)))
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
	}
}

func TestQuadPolynomialSecondOrderAccuracy(t *testing.T) {
	// Second order error approx can exactly represent polynomials up to a power
	// of 2, so that is what we test with
	testForAccuracy(
		t,
		func(x float64) float64 { return math.Pow(x, 2) },
		func(x float64) float64 { return 2 * math.Pow(x, 1) },
		func(x float64) float64 { return 2 },
		func(x float64) float64 { return 0 },
		100,
		&types.State{
			BarPathCalc: types.BarPathCalcConf{
				ApproxErr: types.SecondOrder,
			},
			PhysicsData: types.PhysicsDataConf{
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
		func(x float64) float64 { return math.Pow(x, 4) },
		func(x float64) float64 { return 4 * math.Pow(x, 3) },
		func(x float64) float64 { return 12 * math.Pow(x, 2) },
		func(x float64) float64 { return 24 * math.Pow(x, 1) },
		100,
		&types.State{
			BarPathCalc: types.BarPathCalcConf{
				ApproxErr: types.FourthOrder,
			},
			PhysicsData: types.PhysicsDataConf{
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

func loadAndTestCsv(
	t *testing.T,
	rawDataFile string,
	outFile string,
	state *types.State,
	baseData *dal.BulkCreateTrainingLogsParams,
) {
	f, err := os.ReadFile(rawDataFile)
	sbtest.Nil(t, err)
	csvReader := csv.NewReader(bytes.NewReader(f))
	rawData, err := csvReader.ReadAll()

	samples := len(rawData) - 1
	inputData := dal.CreatePhysicsDataParams{
		Time:         [][]types.Second{make([]types.Second, samples)},
		Position:     [][]types.Vec2[types.Meter]{make([]types.Vec2[types.Meter], samples)},
		Velocity:     [][]types.Vec2[types.MeterPerSec]{make([]types.Vec2[types.MeterPerSec], samples)},
		Acceleration: [][]types.Vec2[types.MeterPerSec2]{make([]types.Vec2[types.MeterPerSec2], samples)},
		Jerk:         [][]types.Vec2[types.MeterPerSec3]{make([]types.Vec2[types.MeterPerSec3], samples)},
		Impulse:      [][]types.Vec2[types.NewtonSec]{make([]types.Vec2[types.NewtonSec], samples)},
		Force:        [][]types.Vec2[types.Newton]{make([]types.Vec2[types.Newton], samples)},
		Work:         [][]types.Joule{make([]types.Joule, samples)},
		Power:        [][]types.Watt{make([]types.Watt, samples)},
		RepSplits:    [][]types.Split[types.Second]{make([]types.Split[types.Second], samples)},
	}
	for i := range samples {
		rawTime, err := strconv.ParseFloat(rawData[i+1][1], 64)
		sbtest.Nil(t, err)
		inputData.Time[0][i] = types.Second(rawTime)

		rawXPos, err := strconv.ParseFloat(rawData[i+1][3], 64)
		sbtest.Nil(t, err)
		rawYPos, err := strconv.ParseFloat(rawData[i+1][2], 64)
		sbtest.Nil(t, err)
		inputData.Position[0][i] = types.Vec2[types.Meter]{
			X: types.Meter(rawXPos) / 100,
			Y: types.Meter(rawYPos) / 100,
		}
	}
	err = Calc(state, baseData, &inputData, 0)
	sbtest.Nil(t, err)

	outF, err := os.Create(outFile)
	sbtest.Nil(t, err)
	csvWriter := csv.NewWriter(outF)
	csvWriter.Write([]string{
		"Time",
		"PosX",
		"PosY",
		"RawVelX",
		"RawVelY",
		"CalcVelX",
		"CalcVelY",
		"RawAccX",
		"RawAccY",
		"CalcAccX",
		"CalcAccY",
		"CalcJerkX",
		"CalcJerkY",
		"CalcForceX",
		"CalcForceY",
		"CalcImpulseX",
		"CalcImpulseY",
		"CalcWork",
		"CalcPower",
	})
	for i := range samples {
		csvWriter.Write([]string{
			fmt.Sprintf("%f", inputData.Time[0][i]),
			fmt.Sprintf("%f", inputData.Position[0][i].X),
			fmt.Sprintf("%f", inputData.Position[0][i].Y),
			rawData[i+1][5],
			rawData[i+1][4],
			fmt.Sprintf("%f", inputData.Velocity[0][i].X),
			fmt.Sprintf("%f", inputData.Velocity[0][i].Y),
			rawData[i+1][7],
			rawData[i+1][6],
			fmt.Sprintf("%f", inputData.Acceleration[0][i].X),
			fmt.Sprintf("%f", inputData.Acceleration[0][i].Y),
			fmt.Sprintf("%f", inputData.Jerk[0][i].X),
			fmt.Sprintf("%f", inputData.Jerk[0][i].Y),
			fmt.Sprintf("%f", inputData.Force[0][i].X),
			fmt.Sprintf("%f", inputData.Force[0][i].Y),
			fmt.Sprintf("%f", inputData.Impulse[0][i].X),
			fmt.Sprintf("%f", inputData.Impulse[0][i].Y),
			fmt.Sprintf("%f", inputData.Work[0][i]),
			fmt.Sprintf("%f", inputData.Power[0][i]),
		})
	}
}

func TestSquatDataSecondOrder(t *testing.T) {
	state := types.State{
		BarPathCalc: types.BarPathCalcConf{
			ApproxErr:       types.SecondOrder,
			SmootherWeight1: 0.5,
			SmootherWeight2: 0.5,
			SmootherWeight3: 1,
			SmootherWeight4: 0.5,
			SmootherWeight5: 0.5,
		},
		PhysicsData: types.PhysicsDataConf{
			MinNumSamples: 10,
			TimeDeltaEps:  1e-2,
		},
	}
	baseData := dal.BulkCreateTrainingLogsParams{
		Weight: 1,
		Reps:   8,
		Sets:   1,
	}
	loadAndTestCsv(
		t,
		"./testData/15_08_2025_squat.csv",
		"./testData/15_08_2025_squat.secondOrder.csv",
		&state,
		&baseData,
	)
}

func TestSquatDataFourthOrder(t *testing.T) {
	state := types.State{
		BarPathCalc: types.BarPathCalcConf{
			ApproxErr:       types.FourthOrder,
			SmootherWeight1: 0.5,
			SmootherWeight2: 0.5,
			SmootherWeight3: 1,
			SmootherWeight4: 0.5,
			SmootherWeight5: 0.5,
		},
		PhysicsData: types.PhysicsDataConf{
			MinNumSamples: 10,
			TimeDeltaEps:  1e-2,
		},
	}
	baseData := dal.BulkCreateTrainingLogsParams{
		Weight: 1,
		Reps:   8,
		Sets:   1,
	}
	loadAndTestCsv(
		t,
		"./testData/15_08_2025_squat.csv",
		"./testData/15_08_2025_squat.fourthOrder.csv",
		&state,
		&baseData,
	)
}

// TODO - test fractional sets produce correct num of reps
