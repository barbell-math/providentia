package barpathphysdata

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"testing"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbtest "code.barbellmath.net/barbell-math/smoothbrain-test"
)

type args struct {
	t           *testing.T
	rawDataFile string
	outFileName string
	expCenters  []types.Split
	params      types.BarPathCalcHyperparams
	numReps     int32
}

func loadAndTestCsv(a *args) {
	f, err := os.ReadFile(a.rawDataFile)
	sbtest.Nil(a.t, err)
	csvReader := csv.NewReader(bytes.NewReader(f))
	rawData, err := csvReader.ReadAll()

	samples := len(rawData) - 1
	inputData := types.PhysicsData{
		Time:     make([]types.Second, samples),
		Position: make([]types.Vec2[types.Meter, types.Meter], samples),
	}
	for i := range samples {
		rawTime, err := strconv.ParseFloat(rawData[i+1][1], 64)
		sbtest.Nil(a.t, err)
		inputData.Time[i] = types.Second(rawTime)

		rawXPos, err := strconv.ParseFloat(rawData[i+1][3], 64)
		sbtest.Nil(a.t, err)
		rawYPos, err := strconv.ParseFloat(rawData[i+1][2], 64)
		sbtest.Nil(a.t, err)
		inputData.Position[i] = types.Vec2[types.Meter, types.Meter]{
			X: types.Meter(rawXPos) / 100,
			Y: types.Meter(rawYPos) / 100,
		}
	}
	err = Calc(&inputData, &a.params, 1, a.numReps)
	sbtest.Nil(a.t, err)

	sbtest.Eq(a.t, len(a.expCenters), len(inputData.RepSplits))
	sbtest.SlicesMatch(a.t, a.expCenters, inputData.RepSplits)

	if a.outFileName == "" {
		return
	}

	timeSeriesFile, err := os.Create(a.outFileName + ".timeSeries.csv")
	sbtest.Nil(a.t, err)
	defer timeSeriesFile.Close()
	timeSeriesWriter := csv.NewWriter(timeSeriesFile)
	timeSeriesWriter.Write([]string{
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
		timeSeriesWriter.Write([]string{
			fmt.Sprintf("%f", inputData.Time[i]),
			fmt.Sprintf("%f", inputData.Position[i].X),
			fmt.Sprintf("%f", inputData.Position[i].Y),
			rawData[i+1][5],
			rawData[i+1][4],
			fmt.Sprintf("%f", inputData.Velocity[i].X),
			fmt.Sprintf("%f", inputData.Velocity[i].Y),
			rawData[i+1][7],
			rawData[i+1][6],
			fmt.Sprintf("%f", inputData.Acceleration[i].X),
			fmt.Sprintf("%f", inputData.Acceleration[i].Y),
			fmt.Sprintf("%f", inputData.Jerk[i].X),
			fmt.Sprintf("%f", inputData.Jerk[i].Y),
			fmt.Sprintf("%f", inputData.Force[i].X),
			fmt.Sprintf("%f", inputData.Force[i].Y),
			fmt.Sprintf("%f", inputData.Impulse[i].X),
			fmt.Sprintf("%f", inputData.Impulse[i].Y),
			fmt.Sprintf("%f", inputData.Work[i]),
			fmt.Sprintf("%f", inputData.Power[i]),
		})
	}
	timeSeriesWriter.Flush()

	repSeriesFile, err := os.Create(a.outFileName + ".repSeries.csv")
	sbtest.Nil(a.t, err)
	defer repSeriesFile.Close()
	repSeriesWriter := csv.NewWriter(repSeriesFile)
	repSeriesWriter.Write([]string{
		"Rep",
		"MinVelTime",
		"MinVel",
		"MaxVelTime",
		"MaxVel",
		"MinAccTime",
		"MinAcc",
		"MaxAccTime",
		"MaxAcc",
		"MinForceTime",
		"MinForce",
		"MaxForceTime",
		"MaxForce",
		"MinImpulseTime",
		"MinImpulse",
		"MaxImpulseTime",
		"MaxImpulse",
		"AvgWork",
		"MinWorkTime",
		"MinWork",
		"MaxWorkTime",
		"MaxWork",
		"AvgPower",
		"MinPowerTime",
		"MinPower",
		"MaxPowerTime",
		"MaxPower",
	})
	for i := range a.numReps {
		repSeriesWriter.Write([]string{
			fmt.Sprintf("%d", i),
			fmt.Sprintf("%f", inputData.MinVel[i].Time),
			fmt.Sprintf("%f", inputData.MinVel[i].Value),
			fmt.Sprintf("%f", inputData.MaxVel[i].Time),
			fmt.Sprintf("%f", inputData.MaxVel[i].Value),

			fmt.Sprintf("%f", inputData.MinAcc[i].Time),
			fmt.Sprintf("%f", inputData.MinAcc[i].Value),
			fmt.Sprintf("%f", inputData.MaxAcc[i].Time),
			fmt.Sprintf("%f", inputData.MaxAcc[i].Value),

			fmt.Sprintf("%f", inputData.MinForce[i].Time),
			fmt.Sprintf("%f", inputData.MinForce[i].Value),
			fmt.Sprintf("%f", inputData.MaxForce[i].Time),
			fmt.Sprintf("%f", inputData.MaxForce[i].Value),

			fmt.Sprintf("%f", inputData.MinImpulse[i].Time),
			fmt.Sprintf("%f", inputData.MinImpulse[i].Value),
			fmt.Sprintf("%f", inputData.MaxImpulse[i].Time),
			fmt.Sprintf("%f", inputData.MaxImpulse[i].Value),

			fmt.Sprintf("%f", inputData.AvgWork[i]),
			fmt.Sprintf("%f", inputData.MinWork[i].Time),
			fmt.Sprintf("%f", inputData.MinWork[i].Value),
			fmt.Sprintf("%f", inputData.MaxWork[i].Time),
			fmt.Sprintf("%f", inputData.MaxWork[i].Value),

			fmt.Sprintf("%f", inputData.AvgPower[i]),
			fmt.Sprintf("%f", inputData.MinPower[i].Time),
			fmt.Sprintf("%f", inputData.MinPower[i].Value),
			fmt.Sprintf("%f", inputData.MaxPower[i].Time),
			fmt.Sprintf("%f", inputData.MaxPower[i].Value),
		})
	}
	repSeriesWriter.Flush()
}

func TestSquatDataSecondOrder(t *testing.T) {
	loadAndTestCsv(&args{
		t:           t,
		rawDataFile: "./testData/15_08_2025_squat.csv",
		outFileName: "./testData/15_08_2025_squat.secondOrder",
		numReps:     8,
		expCenters: []types.Split{
			{StartIdx: 311, EndIdx: 379},
			{StartIdx: 437, EndIdx: 501},
			{StartIdx: 548, EndIdx: 613},
			{StartIdx: 655, EndIdx: 728},
			{StartIdx: 780, EndIdx: 850},
			{StartIdx: 911, EndIdx: 977},
			{StartIdx: 1039, EndIdx: 1106},
			{StartIdx: 1170, EndIdx: 1237},
		},
		params: types.BarPathCalcHyperparams{
			ApproxErr:       types.SecondOrder,
			NearZeroFilter:  0.1,
			SmootherWeight1: 0.5,
			SmootherWeight2: 0.5,
			SmootherWeight3: 1,
			SmootherWeight4: 0.5,
			SmootherWeight5: 0.5,
			MinNumSamples:   10,
			TimeDeltaEps:    1e-2,
			NoiseFilter:     3,
		},
	})
}

func TestSquatDataFourthOrder(t *testing.T) {
	loadAndTestCsv(&args{
		t:           t,
		rawDataFile: "./testData/15_08_2025_squat.csv",
		outFileName: "./testData/15_08_2025_squat.fourthOrder",
		numReps:     8,
		expCenters: []types.Split{
			{StartIdx: 311, EndIdx: 379},
			{StartIdx: 438, EndIdx: 501},
			{StartIdx: 548, EndIdx: 613},
			{StartIdx: 655, EndIdx: 728},
			{StartIdx: 780, EndIdx: 850},
			{StartIdx: 911, EndIdx: 977},
			{StartIdx: 1039, EndIdx: 1106},
			{StartIdx: 1170, EndIdx: 1237},
		},
		params: types.BarPathCalcHyperparams{
			ApproxErr:       types.FourthOrder,
			NearZeroFilter:  0.1,
			SmootherWeight1: 0.5,
			SmootherWeight2: 0.5,
			SmootherWeight3: 1,
			SmootherWeight4: 0.5,
			SmootherWeight5: 0.5,
			MinNumSamples:   10,
			TimeDeltaEps:    1e-2,
			NoiseFilter:     3,
		},
	})
}
