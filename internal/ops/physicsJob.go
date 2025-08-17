package ops

import (
	"context"
	"math"

	dal "code.barbellmath.net/barbell-math/providentia/internal/db/dataAccessLayer"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sbjobqueue "code.barbellmath.net/barbell-math/smoothbrain-jobQueue"
)

type (
	physicsJob struct {
		BarPath []types.BarPathVariant
		Tl      *dal.BulkCreateTrainingLogParams
		B       *sbjobqueue.Batch
		S       *types.State
	}
)

func (p *physicsJob) JobType(_ types.PhysicsJob) {}

func (p *physicsJob) Batch() *sbjobqueue.Batch {
	return p.B
}

func (p *physicsJob) Run(ctxt context.Context) error {
	if err := p.processTimeSeriesData(); err != nil {
		return sberr.AppendError(types.PhysicsJobQueueErr, err)
	}

	// Upload final physics data to db, returning id
	// Set video id in tl params

	return nil
}

func (p *physicsJob) processTimeSeriesData() error {
	for _, set := range p.BarPath {
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
				if math.Abs(iterDelta-delta) < p.S.PhysicsData.TimeDeltaEps {
					return sberr.Wrap(
						types.TimeSeriesNotMonotonicErr,
						"Time samples must all have the same delta (within %f variance), got deltas of %f and %f",
						p.S.PhysicsData.TimeDeltaEps, delta, iterDelta,
					)
				}
			}

			// call c algo to calculate all derivatives and such
		}
	}

	return nil
}
