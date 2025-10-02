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
		BarPath              []types.BarPathVariant
		Tl                   *dal.BulkCreateTrainingLogsParams
		B                    *sbjobqueue.Batch
		Q                    *dal.SyncQueries
		BarPathCalcParams    *types.BarPathCalcHyperparams
		BarTrackerCalcParams *types.BarPathTrackerHyperparams
	}
)

func (p *Physics) JobType(_ types.PhysicsJob) {}

func (p *Physics) Batch() *sbjobqueue.Batch {
	return p.B
}

func (p *Physics) Run(ctxt context.Context) error {
	var physData dal.CreatePhysicsDataParams
	barpathphysdata.InitBarPathCalcPhysicsData(
		&physData, p.BarPathCalcParams, len(p.BarPath),
	)
	// TODO
	// barpathtracker.InitBarPathTrackerData(&physData, p.BarTrackerCalcParams)

	for i, set := range p.BarPath {
		select {
		case <-ctxt.Done():
			return nil
		default:
		}

		if rawData, ok := set.TimeSeriesData(); ok {
			physData.Time[i] = rawData.TimeData
			physData.Position[i] = rawData.PositionData
			if err := barpathphysdata.Calc(
				p.Tl, &physData, p.BarPathCalcParams, i,
			); err != nil {
				return sberr.AppendError(types.PhysicsJobQueueErr, err)
			}
		} else if rawData, ok := set.VideoData(); ok {
			physData.Path[i] = rawData
			// TODO - run video model
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
