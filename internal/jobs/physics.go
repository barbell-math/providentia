package jobs

import (
	"context"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	barpathphysdata "code.barbellmath.net/barbell-math/providentia/internal/models/barPathPhysData"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sbjobqueue "code.barbellmath.net/barbell-math/smoothbrain-jobQueue"
	"github.com/jackc/pgx/v5/pgtype"
)

type (
	Physics struct {
		BarPath []types.BarPathVariant
		Tl      *dal.BulkCreateTrainingLogsParams
		B       *sbjobqueue.Batch
		S       *types.State
		Q       *dal.SyncQueries
	}
)

func (p *Physics) JobType(_ types.PhysicsJob) {}

func (p *Physics) Batch() *sbjobqueue.Batch {
	return p.B
}

func (p *Physics) Run(ctxt context.Context) error {
	physData := dal.CreatePhysicsDataParams{
		Time:         make([][]types.Second, len(p.BarPath)),
		Position:     make([][]types.Vec2[types.Meter], len(p.BarPath)),
		Velocity:     make([][]types.Vec2[types.MeterPerSec], len(p.BarPath)),
		Acceleration: make([][]types.Vec2[types.MeterPerSec2], len(p.BarPath)),
		Jerk:         make([][]types.Vec2[types.MeterPerSec3], len(p.BarPath)),
		Force:        make([][]types.Vec2[types.Newton], len(p.BarPath)),
		Impulse:      make([][]types.Vec2[types.NewtonSec], len(p.BarPath)),
		Work:         make([][]types.Joule, len(p.BarPath)),
		Power:        make([][]types.Watt, len(p.BarPath)),
	}

	if err := p.processData(ctxt, &physData); err != nil {
		return sberr.AppendError(types.PhysicsJobQueueErr, err)
	}

	var id int64
	var err error
	p.Q.Run(func(q *dal.Queries) {
		id, err = q.CreatePhysicsData(ctxt, physData)
	})
	if err != nil {
		return sberr.AppendError(types.PhysicsJobQueueErr, err)
	}
	p.Tl.PhysicsID = pgtype.Int8{Int64: id, Valid: true}

	return nil
}

func (p *Physics) processData(
	ctxt context.Context,
	data *dal.CreatePhysicsDataParams,
) error {
	for i, set := range p.BarPath {
		// TODO - check if ctxt was canceled to stop job early
		if rawData, ok := set.TimeSeriesData(); ok {
			data.Time[i] = rawData.TimeData
			data.Position[i] = rawData.PositionData
			data.Velocity[i] = make(
				[]types.Vec2[types.MeterPerSec], len(rawData.TimeData),
			)
			data.Acceleration[i] = make(
				[]types.Vec2[types.MeterPerSec2], len(rawData.TimeData),
			)
			data.Jerk[i] = make(
				[]types.Vec2[types.MeterPerSec3], len(rawData.TimeData),
			)
			data.Impulse[i] = make(
				[]types.Vec2[types.NewtonSec], len(rawData.TimeData),
			)
			data.Force[i] = make(
				[]types.Vec2[types.Newton], len(rawData.TimeData),
			)
			data.Work[i] = make([]types.Joule, len(rawData.TimeData))
			data.Power[i] = make([]types.Watt, len(rawData.TimeData))

			if err := barpathphysdata.Calc(p.S, p.Tl.Weight, data, i); err != nil {
				return err
			}
		}
	}

	return nil
}

// TODO - dont forget to set path in the physics data!!
