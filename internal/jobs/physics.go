package jobs

import (
	"context"
	"math"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
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
		Time:         make([][]float64, len(p.BarPath)),
		Position:     make([][]float64, len(p.BarPath)),
		Velocity:     make([][]float64, len(p.BarPath)),
		Acceleration: make([][]float64, len(p.BarPath)),
		Jerk:         make([][]float64, len(p.BarPath)),
		Force:        make([][]float64, len(p.BarPath)),
		Work:         make([][]float64, len(p.BarPath)),
		Impulse:      make([][]float64, len(p.BarPath)),
	}

	if err := p.processTimeSeriesData(ctxt, &physData); err != nil {
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

func (p *Physics) processTimeSeriesData(
	ctxt context.Context,
	data *dal.CreatePhysicsDataParams,
) error {
	for i, set := range p.BarPath {
		// TODO - check if ctxt was canceled to stop job early
		if rawData, ok := set.TimeSeriesData(); ok {
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
				if math.Abs(iterDelta-delta) > p.S.PhysicsData.TimeDeltaEps {
					return sberr.Wrap(
						types.TimeSeriesNotMonotonicErr,
						"Time samples must all have the same delta (within %f variance), got deltas of %f and %f",
						p.S.PhysicsData.TimeDeltaEps, delta, iterDelta,
					)
				}
			}

			data.Time[i] = rawData.TimeData
			data.Position[i] = rawData.PositionData

			// call c algo to calculate all derivatives and such
			// should we wait to call the c algo when all time/position data is
			// available?? would make sense from cgo perspective
		}
	}

	return nil
}

// TODO - dont forget to set path in the physics data!!
