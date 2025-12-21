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

	PhysicsOpts struct {
		BarPathCalcParams    *types.BarPathCalcHyperparams
		BarTrackerCalcParams *types.BarPathTrackerHyperparams
		RawData              []types.BarPathVariant
		ExerciseData         *types.ExerciseData
	}
)

func RunPhysicsJobs(
	ctxt context.Context,
	state *types.State,
	tx pgx.Tx,
	opts PhysicsOpts,
) error {
	ceilSets := math.Ceil(opts.ExerciseData.Sets)
	floorSets := math.Floor(opts.ExerciseData.Sets)
	if opts.ExerciseData.Reps <= 0 {
		return sberr.Wrap(
			types.PhysicsJobQueueErr,
			"Supplied exercise data must have at least 1 rep",
		)
	}
	if len(opts.RawData) != int(ceilSets) {
		return sberr.Wrap(
			types.PhysicsJobQueueErr,
			"The length of the raw data (%d) must equal the ceiling of the number of sets (%f)",
			len(opts.RawData), ceilSets,
		)
	}

	batch, _ := sbjobqueue.BatchWithContext(ctxt)
	if len(opts.ExerciseData.PhysData) < len(opts.RawData) {
		opts.ExerciseData.PhysData = make(
			[]types.Optional[types.PhysicsData], len(opts.RawData),
		)
	}
	opts.ExerciseData.PhysData = opts.ExerciseData.PhysData[:len(opts.RawData)]

	var iterPhysData types.PhysicsData
	for i, exerciseSet := range opts.RawData {
		select {
		case <-ctxt.Done():
			return ctxt.Err()
		default:
		}

		if !exerciseSet.Flag.IsValid() {
			return sberr.AppendError(
				types.PhysicsJobQueueErr, types.ErrInvalidBarPathFlag,
			)
		}
		if exerciseSet.Flag == types.NoBarPathData {
			opts.ExerciseData.PhysData[i].Present = false
			continue
		}

		expReps := opts.ExerciseData.Reps
		if ceilSets > opts.ExerciseData.Sets && int(floorSets) == i {
			expReps = max(int32((opts.ExerciseData.Sets-floorSets)*float64(opts.ExerciseData.Reps)), 1)
		}
		state.PhysicsJobQueue.Schedule(&physics{
			B:                    batch,
			S:                    state,
			Tx:                   tx,
			UID:                  UID_CNTR.Add(1),
			BarPathCalcParams:    opts.BarPathCalcParams,
			BarTrackerCalcParams: opts.BarTrackerCalcParams,
			Weight:               opts.ExerciseData.Weight,
			ExpNumReps:           expReps,
			RawData:              opts.RawData[i],
			Results:              &opts.ExerciseData.PhysData[i],
		})

		opts.ExerciseData.PhysData[i] = types.Optional[types.PhysicsData]{
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

	switch p.RawData.Flag {
	case types.VideoBarPathData:
		p.Results.Value.VideoPath = p.RawData.VideoPath
		// TODO - run video model to set time and position data
	case types.TimeSeriesBarPathData:
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
