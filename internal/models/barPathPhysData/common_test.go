package barpathphysdata

import (
	"fmt"
	"testing"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbtest "code.barbellmath.net/barbell-math/smoothbrain-test"
)

func TestStraightLine(t *testing.T) {
	rawData := dal.CreatePhysicsDataParams{
		Time: [][]types.Second{{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		Position: [][]types.Vec2[types.Meter]{{
			{X: 0 * 0, Y: 0 * 0},
			{X: 1 * 1, Y: 1 * 1},
			{X: 2 * 2, Y: 2 * 2},
			{X: 3 * 3, Y: 3 * 3},
			{X: 4 * 4, Y: 4 * 4},
			{X: 5 * 5, Y: 5 * 5},
			{X: 6 * 6, Y: 6 * 6},
			{X: 7 * 7, Y: 7 * 7},
			{X: 8 * 8, Y: 8 * 8},
			{X: 9 * 9, Y: 9 * 9},
			{X: 10 * 10, Y: 10 * 10},
		}},
		Velocity: [][]types.Vec2[types.MeterPerSec]{
			make([]types.Vec2[types.MeterPerSec], 11),
		},
		Acceleration: [][]types.Vec2[types.MeterPerSec2]{
			make([]types.Vec2[types.MeterPerSec2], 11),
		},
		Jerk: [][]types.Vec2[types.MeterPerSec3]{
			make([]types.Vec2[types.MeterPerSec3], 11),
		},
		Work: [][]types.Vec2[types.Joule]{
			make([]types.Vec2[types.Joule], 11),
		},
		Impulse: [][]types.Vec2[types.NewtonSec]{
			make([]types.Vec2[types.NewtonSec], 11),
		},
		Force: [][]types.Vec2[types.Newton]{
			make([]types.Vec2[types.Newton], 11),
		},
	}
	state := types.State{
		BarPathCalc: types.BarPathCalcConf{
			ApproxErr: types.FourthOrder,
		},
		PhysicsData: types.PhysicsDataConf{
			MinNumSamples: 10,
			TimeDeltaEps:  1e-6,
		},
	}
	err := Calc(&state, &rawData, 0)
	sbtest.Nil(t, err)
	fmt.Println("vel: ", rawData.Velocity)
	fmt.Println("acc: ", rawData.Acceleration)
	fmt.Println("jer: ", rawData.Jerk)
}
