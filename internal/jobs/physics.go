package jobs

import (
	"context"
	"math"

	barpathphysdata "code.barbellmath.net/barbell-math/providentia/internal/models/barPathPhysData"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sberr "code.barbellmath.net/barbell-math/smoothbrain-errs"
	sbjobqueue "code.barbellmath.net/barbell-math/smoothbrain-jobQueue"
	sblog "code.barbellmath.net/barbell-math/smoothbrain-logging"
	"github.com/jackc/pgx/v5"
)

type (
	physics struct {
		B   *sbjobqueue.Batch
		S   *types.State
		Tx  pgx.Tx
		UID uint64

		BarPathCalcParams    *types.BarPathCalcHyperparams
		BarTrackerCalcParams *types.BarPathTrackerHyperparams
		Weight               types.Kilogram
		ExpNumReps           int32
		RawData              types.BarPathVariant
		Results              *types.Optional[types.PhysicsData]
	}

	// TODO - make work with exercise data somehow? - will need to check len of
	// slices
	PhysicsOpts struct {
		Weight               types.Kilogram
		Sets                 float64
		Reps                 int32
		BarPathCalcParams    *types.BarPathCalcHyperparams
		BarTrackerCalcParams *types.BarPathTrackerHyperparams
		RawData              []types.BarPathVariant
		Results              *[]types.Optional[types.PhysicsData]
	}
)

func RunPhysicsJobs(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts PhysicsOpts,
) error {
	ceilSets := math.Ceil(opts.Sets)
	floorSets := math.Floor(opts.Sets)
	if opts.Reps <= 0 {
		return sberr.Wrap(
			types.PhysicsJobQueueErr,
			"Supplied data must have at least 1 rep",
		)
	}
	if len(opts.RawData) != int(ceilSets) {
		return sberr.Wrap(
			types.PhysicsJobQueueErr,
			"The length of the supplied data (%d) must equal the ceiling of the number of sets (%d)",
			len(opts.RawData), ceilSets,
		)
	}

	batch, _ := sbjobqueue.BatchWithContext(ctxt)
	if len(*opts.Results) < len(opts.RawData) {
		*opts.Results = make(
			[]types.Optional[types.PhysicsData], len(opts.RawData),
		)
	}
	*opts.Results = (*opts.Results)[:len(opts.RawData)]

	var iterPhysData types.PhysicsData
	for i, exerciseSet := range opts.RawData {
		select {
		case <-ctxt.Done():
			return nil
		default:
		}

		if exerciseSet.Flag.IsValid() {
			expReps := opts.Reps
			if ceilSets > opts.Sets && int(floorSets) == i {
				expReps = max(int32((opts.Sets-floorSets)*float64(opts.Reps)), 1)
			}

			state.PhysicsJobQueue.Schedule(&physics{
				B:                    batch,
				S:                    state,
				Tx:                   tx,
				UID:                  UID_CNTR.Add(1),
				BarPathCalcParams:    opts.BarPathCalcParams,
				BarTrackerCalcParams: opts.BarTrackerCalcParams,
				Weight:               opts.Weight,
				ExpNumReps:           expReps,
				RawData:              opts.RawData[i],
				Results:              &(*opts.Results)[i],
			})
		} else {
			(*opts.Results)[i].Present = false
			continue
		}

		(*opts.Results)[i] = types.Optional[types.PhysicsData]{
			Present: true,
			Value:   iterPhysData,
		}
		iterPhysData.Time = iterPhysData.Time[:0]
		iterPhysData.Position = iterPhysData.Position[:0]
	}

	return batch.Wait()
}

func (p *physics) JobType(_ types.PhysicsJob) {}

func (p *physics) Batch() *sbjobqueue.Batch {
	return p.B
}

func (p *physics) formatLogLine(msg string) string {
	return formatJobLogLine("physics", p.UID, msg)
}

func (p *physics) Run(ctxt context.Context) (opErr error) {
	p.S.Log.Log(ctxt, sblog.VLevel(3), p.formatLogLine("Starting..."))

	if p.RawData.Flag == types.VideoBarPathData {
		p.Results.Value.VideoPath = p.RawData.VideoPath
		// TODO - run video model to set time and position data
	} else if p.RawData.Flag == types.TimeSeriesBarPathData {
		p.Results.Value.Time = p.RawData.TimeSeries.TimeData
		p.Results.Value.Position = p.RawData.TimeSeries.PositionData
	}

	if opErr = barpathphysdata.Calc(
		&p.Results.Value, p.BarPathCalcParams, p.Weight, p.ExpNumReps,
	); opErr != nil {
		goto errReturn
	}

	p.Results.Present = true
	p.S.Log.Log(
		ctxt, sblog.VLevel(3),
		p.formatLogLine("Finished processing physics data"),
	)
	return nil
errReturn:
	p.S.Log.Error(p.formatLogLine("Encountered error"), "Error", opErr)
	return sberr.AppendError(types.PhysicsJobQueueErr, opErr)
}
