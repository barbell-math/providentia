package logic

import (
	"testing"

	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbtest "code.barbellmath.net/barbell-math/smoothbrain-test"
)

func rawWorkoutEqSavedWorkout(
	t *testing.T,
	l []types.RawWorkout,
	r []types.Workout,
) {
	sbtest.Eq(t, len(l), len(r))

	for i := range len(l) {
		sbtest.Eq(t, l[i].WorkoutID, r[i].WorkoutID)
		sbtest.Eq(t, len(l[i].Exercises), len(r[i].Exercises))

		for j := range len(l[i].Exercises) {
			sbtest.Eq(t, l[i].Exercises[j].Name, r[i].Exercises[j].Name)
			sbtest.Eq(t, l[i].Exercises[j].Weight, r[i].Exercises[j].Weight)
			sbtest.Eq(t, l[i].Exercises[j].Sets, r[i].Exercises[j].Sets)
			sbtest.Eq(t, l[i].Exercises[j].Reps, r[i].Exercises[j].Reps)
			sbtest.Eq(t, l[i].Exercises[j].Effort, r[i].Exercises[j].Effort)

			sbtest.Eq(
				t,
				types.Kilogram(l[i].Exercises[j].Sets*float64(l[i].Exercises[j].Reps))*l[i].Exercises[j].Weight,
				r[i].Exercises[j].Volume,
			)
			sbtest.Eq(
				t,
				types.RPE(l[i].Exercises[j].Sets*float64(l[i].Exercises[j].Reps))*l[i].Exercises[j].Effort,
				r[i].Exercises[j].Exertion,
			)
			sbtest.Eq(
				t,
				l[i].Exercises[j].Sets*float64(l[i].Exercises[j].Reps),
				r[i].Exercises[j].TotalReps,
			)

			for k := range len(l[i].Exercises[j].BarPath) {
				if l[i].Exercises[j].BarPath[k].Source() == types.TimeSeriesBarPathData {
					rawData, ok := l[i].Exercises[j].BarPath[k].TimeSeriesData()
					sbtest.True(t, ok)
					sbtest.SlicesMatch(t, rawData.TimeData, r[i].Exercises[j].Time[k])
					sbtest.SlicesMatch(t, rawData.PositionData, r[i].Exercises[j].Position[k])
					sbtest.Eq(t, len(rawData.TimeData), len(r[i].Exercises[j].Velocity[k]))
					sbtest.Eq(t, len(rawData.TimeData), len(r[i].Exercises[j].Acceleration[k]))
					sbtest.Eq(t, len(rawData.TimeData), len(r[i].Exercises[j].Jerk[k]))
					sbtest.Eq(t, len(rawData.TimeData), len(r[i].Exercises[j].Force[k]))
					sbtest.Eq(t, len(rawData.TimeData), len(r[i].Exercises[j].Impulse[k]))
					sbtest.Eq(t, len(rawData.TimeData), len(r[i].Exercises[j].Work[k]))
				}
				if l[i].Exercises[j].BarPath[k].Source() == types.VideoBarPathData {
					// TODO - enable when video processing is working
					// sbtest.True(t, len(r[i].Exercises[j].Time[k]) > 0)
					// sbtest.Eq(t, len(r[i].Exercises[j].Time[k]), len(r[i].Exercises[j].Position[k]))
					// sbtest.Eq(t, len(r[i].Exercises[j].Time[k]), len(r[i].Exercises[j].Velocity[k]))
					// sbtest.Eq(t, len(r[i].Exercises[j].Time[k]), len(r[i].Exercises[j].Acceleration[k]))
					// sbtest.Eq(t, len(r[i].Exercises[j].Time[k]), len(r[i].Exercises[j].Jerk[k]))
					// sbtest.Eq(t, len(r[i].Exercises[j].Time[k]), len(r[i].Exercises[j].Impulse[k]))
					// sbtest.Eq(t, len(r[i].Exercises[j].Time[k]), len(r[i].Exercises[j].Work[k]))
				}
			}
		}
	}
}
