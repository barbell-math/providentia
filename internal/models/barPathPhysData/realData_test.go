package barpathphysdata

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"testing"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbtest "code.barbellmath.net/barbell-math/smoothbrain-test"
)

type args struct {
	t           *testing.T
	rawDataFile string
	outFile     string
	expCenters  [][]types.Split
	state       types.State
	baseData    dal.BulkCreateTrainingLogsParams
}

func loadAndTestCsv(a *args) {
	f, err := os.ReadFile(a.rawDataFile)
	sbtest.Nil(a.t, err)
	csvReader := csv.NewReader(bytes.NewReader(f))
	rawData, err := csvReader.ReadAll()

	samples := len(rawData) - 1
	inputData := dal.CreatePhysicsDataParams{
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
		rawTime, err := strconv.ParseFloat(rawData[i+1][1], 64)
		sbtest.Nil(a.t, err)
		inputData.Time[0][i] = types.Second(rawTime)

		rawXPos, err := strconv.ParseFloat(rawData[i+1][3], 64)
		sbtest.Nil(a.t, err)
		rawYPos, err := strconv.ParseFloat(rawData[i+1][2], 64)
		sbtest.Nil(a.t, err)
		inputData.Position[0][i] = types.Vec2[types.Meter, types.Meter]{
			X: types.Meter(rawXPos) / 100,
			Y: types.Meter(rawYPos) / 100,
		}
	}
	err = Calc(&a.state, &a.baseData, &inputData, 0)
	sbtest.Nil(a.t, err)

	sbtest.Eq(a.t, len(a.expCenters), len(inputData.RepSplits))
	for i := range len(a.expCenters) {
		sbtest.SlicesMatch(a.t, a.expCenters[i], inputData.RepSplits[i])
	}

	if a.outFile == "" {
		return
	}

	outF, err := os.Create(a.outFile)
	sbtest.Nil(a.t, err)
	defer outF.Close()
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
	loadAndTestCsv(&args{
		t:           t,
		rawDataFile: "./testData/15_08_2025_squat.csv",
		outFile:     "./testData/15_08_2025_squat.secondOrder.csv",
		expCenters: [][]types.Split{
			[]types.Split{
				{StartIdx: 311, EndIdx: 379},
				{StartIdx: 437, EndIdx: 501},
				{StartIdx: 548, EndIdx: 613},
				{StartIdx: 655, EndIdx: 728},
				{StartIdx: 784, EndIdx: 850},
				{StartIdx: 911, EndIdx: 977},
				{StartIdx: 1039, EndIdx: 1106},
				{StartIdx: 1170, EndIdx: 1237},
			},
		},
		state: types.State{
			BarPathCalc: types.BarPathCalcConf{
				ApproxErr:       types.SecondOrder,
				NearZeroFilter:  0.1,
				SmootherWeight1: 0.5,
				SmootherWeight2: 0.5,
				SmootherWeight3: 1,
				SmootherWeight4: 0.5,
				SmootherWeight5: 0.5,
				MinNumSamples:   10,
				TimeDeltaEps:    1e-2,
			},
		},
		baseData: dal.BulkCreateTrainingLogsParams{
			Weight: 1,
			Reps:   8,
			Sets:   1,
		},
	})
}

func TestSquatDataFourthOrder(t *testing.T) {
	loadAndTestCsv(&args{
		t:           t,
		rawDataFile: "./testData/15_08_2025_squat.csv",
		outFile:     "./testData/15_08_2025_squat.fourthOrder.csv",
		expCenters: [][]types.Split{
			[]types.Split{
				{StartIdx: 311, EndIdx: 379},
				{StartIdx: 437, EndIdx: 501},
				{StartIdx: 548, EndIdx: 613},
				{StartIdx: 655, EndIdx: 728},
				{StartIdx: 784, EndIdx: 850},
				{StartIdx: 911, EndIdx: 977},
				{StartIdx: 1039, EndIdx: 1106},
				{StartIdx: 1170, EndIdx: 1237},
			},
		},
		state: types.State{
			BarPathCalc: types.BarPathCalcConf{
				ApproxErr:       types.FourthOrder,
				NearZeroFilter:  0.1,
				SmootherWeight1: 0.5,
				SmootherWeight2: 0.5,
				SmootherWeight3: 1,
				SmootherWeight4: 0.5,
				SmootherWeight5: 0.5,
				MinNumSamples:   10,
				TimeDeltaEps:    1e-2,
			},
		},
		baseData: dal.BulkCreateTrainingLogsParams{
			Weight: 1,
			Reps:   8,
			Sets:   1,
		},
	})
}
