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
		RepSplits:    make([][]types.Split, len(p.BarPath)),
	}

	for i, set := range p.BarPath {
		select {
		case <-ctxt.Done():
			return nil
		default:
		}

		if rawData, ok := set.TimeSeriesData(); ok {
			physData.Time[i] = rawData.TimeData
			physData.Position[i] = rawData.PositionData
			if err := barpathphysdata.Calc(p.S, p.Tl, &physData, i); err != nil {
				return sberr.AppendError(types.PhysicsJobQueueErr, err)
			}
		} else if _, ok := set.VideoData(); ok {
			// TODO - dont forget to set path in the physics data!!
		}
	}

	var id int64
	var err error
	id, err = dal.Query1x2(dal.Q.CreatePhysicsData, p.Q, ctxt, physData)
	if err != nil {
		return sberr.AppendError(types.PhysicsJobQueueErr, err)
	}
	p.Tl.PhysicsID = pgtype.Int8{Int64: id, Valid: true}

	return nil
}
