package dal

import (
	"context"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
	"github.com/jackc/pgx/v5"
)

type (
	createPhysicsDataReturningIdVals struct {
		*types.PhysicsData
		*int64
	}
)

const (
	createPhysicsDataReturningIdsSql = `
INSERT INTO providentia.physics_data (
	path, bar_path_calc_id, bar_path_track_id,

	time, position, velocity, acceleration, jerk,
	force, impulse, work, power,

	rep_splits,

	min_vel, max_vel,
	min_acc, max_acc,
	min_force, max_force,
	min_impulse, max_impulse,
	avg_work, min_work, max_work,
	avg_power, min_power, max_power
) VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
	$18, $19, $20, $21, $22, $23, $24, $25, $26, $27
) RETURNING id;
`
)

func CreatePhysicsDataReturningIds(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	data []createPhysicsDataReturningIdVals,
) error {
	cntr := 0
	b := pgx.Batch{}
	for _, iterPhysData := range data {
		select {
		case <-ctxt.Done():
			return ctxt.Err()
		default:
		}

		b.Queue(
			createPhysicsDataReturningIdsSql,
			iterPhysData.PhysicsData.VideoPath,
			iterPhysData.PhysicsData.BarPathCalcVersion,
			iterPhysData.PhysicsData.BarPathTrackerVersion,
			iterPhysData.PhysicsData.Time,
			iterPhysData.PhysicsData.Position,
			iterPhysData.PhysicsData.Velocity,
			iterPhysData.PhysicsData.Acceleration,
			iterPhysData.PhysicsData.Jerk,
			iterPhysData.PhysicsData.Force,
			iterPhysData.PhysicsData.Impulse,
			iterPhysData.PhysicsData.Work,
			iterPhysData.PhysicsData.Power,
			iterPhysData.PhysicsData.RepSplits,
			iterPhysData.PhysicsData.MinVel,
			iterPhysData.PhysicsData.MaxVel,
			iterPhysData.PhysicsData.MinAcc,
			iterPhysData.PhysicsData.MaxAcc,
			iterPhysData.PhysicsData.MinForce,
			iterPhysData.PhysicsData.MaxForce,
			iterPhysData.PhysicsData.MinImpulse,
			iterPhysData.PhysicsData.MaxImpulse,
			iterPhysData.PhysicsData.AvgWork,
			iterPhysData.PhysicsData.MinWork,
			iterPhysData.PhysicsData.MaxWork,
			iterPhysData.PhysicsData.AvgPower,
			iterPhysData.PhysicsData.MinPower,
			iterPhysData.PhysicsData.MaxPower,
		)

		if uint(b.Len()) >= state.Global.BatchSize {
			results := tx.SendBatch(ctxt, &b)
			for range b.Len() {
				row := results.QueryRow()
				if err := row.Scan(data[cntr]); err != nil {
					return sberr.AppendError(
						types.CouldNotCreateAllPhysicsDataErr, err,
					)
				}
				cntr++
			}

			results.Close()
			b = pgx.Batch{}
		}
	}

	state.Log.Log(
		ctxt, sblog.VLevel(3),
		"DAL: Created new physics_data entries",
		"NumRows", len(data),
	)
	return nil
}
