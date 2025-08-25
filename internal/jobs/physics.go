package jobs

import (
	"context"
	"math"

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
		Work:         make([][]types.Vec2[types.Joule], len(p.BarPath)),
		Impulse:      make([][]types.Vec2[types.NewtonSec], len(p.BarPath)),
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
			// TODO - performance: move these checks to the C code
			// Would make error messages more opaque...
			delta := rawData.TimeData[1] - rawData.TimeData[0]
			for i := 1; i < len(rawData.TimeData); i++ {
				iterDelta := rawData.TimeData[i] - rawData.TimeData[i-1]
				if iterDelta < 0 {
					return sberr.Wrap(
						types.TimeSeriesDecreaseErr,
						"Time samples must be increasing, got a delta of %f",
						iterDelta,
					)
				}
				if math.Abs(float64(iterDelta-delta)) > float64(p.S.PhysicsData.TimeDeltaEps) {
					return sberr.Wrap(
						types.TimeSeriesNotMonotonicErr,
						"Time samples must all have the same delta (within %f variance), got deltas of %f and %f",
						p.S.PhysicsData.TimeDeltaEps, delta, iterDelta,
					)
				}
			}

			if len(rawData.TimeData) < int(p.S.PhysicsData.MinNumSamples) {
				// TODO - return err
				// Where tf should all these checks go?? It feels like they
				// don't have a "home"
			}

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
			data.Work[i] = make(
				[]types.Vec2[types.Joule], len(rawData.TimeData),
			)
			data.Impulse[i] = make(
				[]types.Vec2[types.NewtonSec], len(rawData.TimeData),
			)
			data.Force[i] = make(
				[]types.Vec2[types.Newton], len(rawData.TimeData),
			)

			barpathphysdata.Calc(p.S, data, i)
		}
	}

	return nil
}

// TODO - dont forget to set path in the physics data!!
