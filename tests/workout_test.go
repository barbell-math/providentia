package tests

import (
	"context"
	"testing"
	"time"

	"code.barbellmath.net/barbell-math/providentia/internal/util"
	"code.barbellmath.net/barbell-math/providentia/lib/logic"
	"code.barbellmath.net/barbell-math/providentia/lib/types"
	sbtest "code.barbellmath.net/barbell-math/smoothbrain-test"
)

var (
	testPhysicsData1 = types.PhysicsData{
		VideoPath: "",
		Time:      []types.Second{0, 1, 2, 3},
		Position: []types.Vec2[types.Meter, types.Meter]{
			{X: 0, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 2}, {X: 3, Y: 3},
		},
		Velocity:     []types.Vec2[types.MeterPerSec, types.MeterPerSec]{},
		Acceleration: []types.Vec2[types.MeterPerSec2, types.MeterPerSec2]{},
		Jerk:         []types.Vec2[types.MeterPerSec3, types.MeterPerSec3]{},
		Force:        []types.Vec2[types.Newton, types.Newton]{},
		Impulse:      []types.Vec2[types.NewtonSec, types.NewtonSec]{},
		Work:         []types.Joule{},
		Power:        []types.Watt{},

		RepSplits: []types.Split{},

		MinVel: []types.PointInTime[types.Second, types.MeterPerSec]{},
		MaxVel: []types.PointInTime[types.Second, types.MeterPerSec]{},

		MinAcc: []types.PointInTime[types.Second, types.MeterPerSec2]{},
		MaxAcc: []types.PointInTime[types.Second, types.MeterPerSec2]{},

		MinForce: []types.PointInTime[types.Second, types.Newton]{},
		MaxForce: []types.PointInTime[types.Second, types.Newton]{},

		MinImpulse: []types.PointInTime[types.Second, types.NewtonSec]{},
		MaxImpulse: []types.PointInTime[types.Second, types.NewtonSec]{},

		AvgWork: []types.Joule{},
		MinWork: []types.PointInTime[types.Second, types.Joule]{},
		MaxWork: []types.PointInTime[types.Second, types.Joule]{},

		AvgPower: []types.Watt{},
		MinPower: []types.PointInTime[types.Second, types.Watt]{},
		MaxPower: []types.PointInTime[types.Second, types.Watt]{},
	}

	testPhysicsData2 = types.PhysicsData{
		VideoPath: "",
		Time:      []types.Second{0, 1, 2, 3, 4},
		Position: []types.Vec2[types.Meter, types.Meter]{
			{X: 0, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 2}, {X: 3, Y: 3},
			{X: 4, Y: 4},
		},
		Velocity:     []types.Vec2[types.MeterPerSec, types.MeterPerSec]{},
		Acceleration: []types.Vec2[types.MeterPerSec2, types.MeterPerSec2]{},
		Jerk:         []types.Vec2[types.MeterPerSec3, types.MeterPerSec3]{},
		Force:        []types.Vec2[types.Newton, types.Newton]{},
		Impulse:      []types.Vec2[types.NewtonSec, types.NewtonSec]{},
		Work:         []types.Joule{},
		Power:        []types.Watt{},

		RepSplits: []types.Split{},

		MinVel: []types.PointInTime[types.Second, types.MeterPerSec]{},
		MaxVel: []types.PointInTime[types.Second, types.MeterPerSec]{},

		MinAcc: []types.PointInTime[types.Second, types.MeterPerSec2]{},
		MaxAcc: []types.PointInTime[types.Second, types.MeterPerSec2]{},

		MinForce: []types.PointInTime[types.Second, types.Newton]{},
		MaxForce: []types.PointInTime[types.Second, types.Newton]{},

		MinImpulse: []types.PointInTime[types.Second, types.NewtonSec]{},
		MaxImpulse: []types.PointInTime[types.Second, types.NewtonSec]{},

		AvgWork: []types.Joule{},
		MinWork: []types.PointInTime[types.Second, types.Joule]{},
		MaxWork: []types.PointInTime[types.Second, types.Joule]{},

		AvgPower: []types.Watt{},
		MinPower: []types.PointInTime[types.Second, types.Watt]{},
		MaxPower: []types.PointInTime[types.Second, types.Watt]{},
	}
)

func optionalWorkoutsEqual(
	t *testing.T,
	l []types.Optional[types.Workout],
	r []types.Optional[types.Workout],
) {
	sbtest.Eq(t, len(l), len(r))
	for i := range len(l) {
		sbtest.Eq(t, l[i].Present, r[i].Present)
		workoutsEqual(
			t, []types.Workout{l[i].Value}, []types.Workout{r[i].Value},
		)
	}
}

func workoutsEqual(
	t *testing.T,
	l []types.Workout,
	r []types.Workout,
) {
	sbtest.Eq(t, len(l), len(r))

	for i := range len(l) {
		sbtest.Eq(t, l[i].WorkoutId.ClientEmail, r[i].WorkoutId.ClientEmail)
		sbtest.Eq(t, l[i].WorkoutId.Session, r[i].WorkoutId.Session)
		sbtest.True(t, util.DateEqual(
			l[i].WorkoutId.DatePerformed, r[i].WorkoutId.DatePerformed,
		))
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
				r[i].Exercises[j].AbstractData.Value.Volume,
			)
			sbtest.Eq(
				t,
				types.RPE(l[i].Exercises[j].Sets*float64(l[i].Exercises[j].Reps))*l[i].Exercises[j].Effort,
				r[i].Exercises[j].AbstractData.Value.Exertion,
			)
			sbtest.Eq(
				t,
				l[i].Exercises[j].Sets*float64(l[i].Exercises[j].Reps),
				r[i].Exercises[j].AbstractData.Value.TotalReps,
			)

			for k := range len(l[i].Exercises[j].PhysData) {
				sbtest.Eq(
					t,
					l[i].Exercises[j].PhysData[k].Present,
					r[i].Exercises[j].PhysData[k].Present,
				)
				sbtest.SlicesMatch(
					t,
					l[i].Exercises[j].PhysData[k].Value.Time,
					r[i].Exercises[j].PhysData[k].Value.Time,
				)
				sbtest.SlicesMatch(
					t,
					l[i].Exercises[j].PhysData[k].Value.Position,
					r[i].Exercises[j].PhysData[k].Value.Position,
				)
				sbtest.SlicesMatch(
					t,
					l[i].Exercises[j].PhysData[k].Value.Velocity,
					r[i].Exercises[j].PhysData[k].Value.Velocity,
				)
				sbtest.SlicesMatch(
					t,
					l[i].Exercises[j].PhysData[k].Value.Acceleration,
					r[i].Exercises[j].PhysData[k].Value.Acceleration,
				)
				sbtest.SlicesMatch(
					t,
					l[i].Exercises[j].PhysData[k].Value.Jerk,
					r[i].Exercises[j].PhysData[k].Value.Jerk,
				)
				sbtest.SlicesMatch(
					t,
					l[i].Exercises[j].PhysData[k].Value.Force,
					r[i].Exercises[j].PhysData[k].Value.Force,
				)
				sbtest.SlicesMatch(
					t,
					l[i].Exercises[j].PhysData[k].Value.Impulse,
					r[i].Exercises[j].PhysData[k].Value.Impulse,
				)
				sbtest.SlicesMatch(
					t,
					l[i].Exercises[j].PhysData[k].Value.Work,
					r[i].Exercises[j].PhysData[k].Value.Work,
				)
				sbtest.SlicesMatch(
					t,
					l[i].Exercises[j].PhysData[k].Value.Power,
					r[i].Exercises[j].PhysData[k].Value.Power,
				)
			}
		}
	}
}

func TestWorkout(t *testing.T) {
	t.Run("createReadNoPhysData", workoutCreateReadNoPhysData)
	t.Run("createReadPhysData", workoutCreateReadPhysData)
	t.Run("createFindNoPhysData", workoutCreateFindNoPhysData)
	t.Run("createFindPhysData", workoutCreateFindPhysData)
	t.Run("createFindBetweenDates", workoutCreateFindBetweenDates)
	t.Run("createDeletePhysData", workoutCreateDeletePhysData)
	t.Run("createDeleteBetweenDates", workoutCreateDeleteBetweenDates)
}

func workoutCreateReadNoPhysData(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	err := logic.CreateClients(ctxt, types.Client{
		FirstName: "FName",
		LastName:  "LName",
		Email:     "email@email.com",
	})
	sbtest.Nil(t, err)

	workouts := []types.Workout{
		{
			WorkoutId: types.WorkoutId{
				ClientEmail:   "email@email.com",
				Session:       1,
				DatePerformed: time.Now(),
			},
			Exercises: []types.ExerciseData{
				{
					Name:   "Squat",
					Weight: 365,
					Sets:   5,
					Reps:   5,
					Effort: 10,
				},
			},
		},
	}
	err = logic.CreateWorkouts(ctxt, workouts...)
	sbtest.Nil(t, err)

	res, err := logic.ReadWorkoutsById(ctxt, workouts[0].WorkoutId)
	sbtest.Nil(t, err)
	workoutsEqual(t, workouts, res)

	n, err := logic.ReadNumWorkoutsForClient(ctxt, "email@email.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, n)

	err = logic.CreateWorkouts(ctxt, workouts...)
	sbtest.ContainsError(t, types.CouldNotCreateAllWorkoutsErr, err)
	sbtest.ContainsError(
		t, types.CouldNotCreateAllTrainingLogsErr, err,
		`duplicate key value violates unique constraint "training_log_client_id_date_performed_inter_session_cntr_in_key" \(SQLSTATE 23505\)`,
	)

	res, err = logic.ReadWorkoutsById(ctxt, types.WorkoutId{
		ClientEmail:   "asdf",
		Session:       1,
		DatePerformed: time.Now(),
	})
	sbtest.ContainsError(
		t, types.CouldNotReadAllWorkoutsErr, err,
		`Could not read entry with id '{ClientEmail:asdf Session:1 DatePerformed:.*}' \(Does id exist\?\)`,
	)

	n, err = logic.ReadNumWorkoutsForClient(ctxt, "email@email.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, n)
}

func workoutCreateReadPhysData(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	err := logic.CreateClients(ctxt, types.Client{
		FirstName: "FName",
		LastName:  "LName",
		Email:     "email@email.com",
	})
	sbtest.Nil(t, err)

	workouts := []types.Workout{
		{
			WorkoutId: types.WorkoutId{
				ClientEmail:   "email@email.com",
				Session:       1,
				DatePerformed: time.Now(),
			},
			Exercises: []types.ExerciseData{
				{
					Name:   "Squat",
					Weight: 365,
					Sets:   5,
					Reps:   5,
					Effort: 10,
					PhysData: []types.Optional[types.PhysicsData]{
						{Present: true, Value: testPhysicsData1},
						{Present: true, Value: testPhysicsData2},
					},
				},
			},
		},
	}
	err = logic.CreateWorkouts(ctxt, workouts...)
	sbtest.Nil(t, err)

	res, err := logic.ReadWorkoutsById(ctxt, workouts[0].WorkoutId)
	sbtest.Nil(t, err)
	workoutsEqual(t, workouts, res)

	n, err := logic.ReadNumWorkoutsForClient(ctxt, "email@email.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, n)

	err = logic.CreateWorkouts(ctxt, workouts...)
	sbtest.ContainsError(t, types.CouldNotCreateAllWorkoutsErr, err)
	sbtest.ContainsError(
		t, types.CouldNotCreateAllTrainingLogsErr, err,
		`duplicate key value violates unique constraint "training_log_client_id_date_performed_inter_session_cntr_in_key" \(SQLSTATE 23505\)`,
	)

	res, err = logic.ReadWorkoutsById(ctxt, types.WorkoutId{
		ClientEmail:   "asdf",
		Session:       1,
		DatePerformed: time.Now(),
	})
	sbtest.ContainsError(
		t, types.CouldNotReadAllWorkoutsErr, err,
		`Could not read entry with id '{ClientEmail:asdf Session:1 DatePerformed:.*}' \(Does id exist\?\)`,
	)

	n, err = logic.ReadNumWorkoutsForClient(ctxt, "email@email.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, n)
}

func workoutCreateFindNoPhysData(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	err := logic.CreateClients(ctxt, types.Client{
		FirstName: "FName",
		LastName:  "LName",
		Email:     "email@email.com",
	})
	sbtest.Nil(t, err)

	workouts := []types.Optional[types.Workout]{
		{
			Present: true,
			Value: types.Workout{
				WorkoutId: types.WorkoutId{
					ClientEmail:   "email@email.com",
					Session:       1,
					DatePerformed: time.Now(),
				},
				Exercises: []types.ExerciseData{
					{
						Name:   "Squat",
						Weight: 365,
						Sets:   5,
						Reps:   5,
						Effort: 10,
					},
				},
			},
		},
	}
	err = logic.CreateWorkouts(ctxt, workouts[0].Value)
	sbtest.Nil(t, err)

	res, err := logic.FindWorkoutsById(ctxt, workouts[0].Value.WorkoutId)
	sbtest.Nil(t, err)
	optionalWorkoutsEqual(t, workouts, res)

	n, err := logic.ReadNumWorkoutsForClient(ctxt, "email@email.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, n)

	res, err = logic.FindWorkoutsById(ctxt, types.WorkoutId{
		ClientEmail:   "asdf",
		Session:       1,
		DatePerformed: time.Now(),
	})
	sbtest.Nil(t, err)
	sbtest.Eq(t, len(res), 1)
	sbtest.False(t, res[0].Present)

	n, err = logic.ReadNumWorkoutsForClient(ctxt, "email@email.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, n)
}

func workoutCreateFindPhysData(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	err := logic.CreateClients(ctxt, types.Client{
		FirstName: "FName",
		LastName:  "LName",
		Email:     "email@email.com",
	})
	sbtest.Nil(t, err)

	workouts := []types.Optional[types.Workout]{
		{
			Present: true,
			Value: types.Workout{
				WorkoutId: types.WorkoutId{
					ClientEmail:   "email@email.com",
					Session:       1,
					DatePerformed: time.Now(),
				},
				Exercises: []types.ExerciseData{
					{
						Name:   "Squat",
						Weight: 365,
						Sets:   5,
						Reps:   5,
						Effort: 10,
						PhysData: []types.Optional[types.PhysicsData]{
							{Present: true, Value: testPhysicsData1},
							{Present: true, Value: testPhysicsData2},
						},
					},
				},
			},
		},
	}
	err = logic.CreateWorkouts(ctxt, workouts[0].Value)
	sbtest.Nil(t, err)

	res, err := logic.FindWorkoutsById(ctxt, workouts[0].Value.WorkoutId)
	sbtest.Nil(t, err)
	optionalWorkoutsEqual(t, workouts, res)

	n, err := logic.ReadNumWorkoutsForClient(ctxt, "email@email.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, n)

	res, err = logic.FindWorkoutsById(ctxt, types.WorkoutId{
		ClientEmail:   "asdf",
		Session:       1,
		DatePerformed: time.Now(),
	})
	sbtest.Nil(t, err)
	sbtest.Eq(t, len(res), 1)
	sbtest.False(t, res[0].Present)

	n, err = logic.ReadNumWorkoutsForClient(ctxt, "email@email.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, n)
}

func workoutCreateFindBetweenDates(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	err := logic.CreateClients(ctxt, types.Client{
		FirstName: "FName",
		LastName:  "LName",
		Email:     "email@email.com",
	})
	sbtest.Nil(t, err)

	startTime := time.Now()

	workouts := []types.Workout{
		{
			WorkoutId: types.WorkoutId{
				ClientEmail:   "email@email.com",
				Session:       1,
				DatePerformed: startTime,
			},
			Exercises: []types.ExerciseData{
				{
					Name:   "Squat",
					Weight: 365,
					Sets:   2,
					Reps:   2,
					Effort: 10,
					PhysData: []types.Optional[types.PhysicsData]{
						{Present: true, Value: testPhysicsData1},
						{Present: true, Value: testPhysicsData2},
					},
				},
			},
		},
		{
			WorkoutId: types.WorkoutId{
				ClientEmail:   "email@email.com",
				Session:       1,
				DatePerformed: startTime.Add(24 * time.Hour),
			},
			Exercises: []types.ExerciseData{
				{
					Name:   "Bench",
					Weight: 225,
					Sets:   1,
					Reps:   1,
					Effort: 10,
					PhysData: []types.Optional[types.PhysicsData]{
						{Present: true, Value: testPhysicsData1},
					},
				},
			},
		},
		{
			WorkoutId: types.WorkoutId{
				ClientEmail:   "email@email.com",
				Session:       1,
				DatePerformed: startTime.Add(48 * time.Hour),
			},
			Exercises: []types.ExerciseData{
				{
					Name:   "Deadlift",
					Weight: 405,
					Sets:   1,
					Reps:   1,
					Effort: 10,
				},
			},
		},
	}
	err = logic.CreateWorkouts(ctxt, workouts...)
	sbtest.Nil(t, err)

	res, err := logic.FindWorkoutsInDateRange(
		ctxt, workouts[0].WorkoutId.ClientEmail,
		time.Now().Add(-1*time.Hour), startTime.Add(72*time.Hour),
	)
	sbtest.Nil(t, err)
	workoutsEqual(t, workouts, res)

	res, err = logic.FindWorkoutsInDateRange(
		ctxt, workouts[0].WorkoutId.ClientEmail,
		time.Now().Add(-1*time.Hour), time.Now().Add(48*time.Hour),
	)
	sbtest.Nil(t, err)
	workoutsEqual(t, workouts[:2], res)

	res, err = logic.FindWorkoutsInDateRange(
		ctxt, workouts[0].WorkoutId.ClientEmail,
		time.Now().Add(23*time.Hour), time.Now().Add(48*time.Hour),
	)
	sbtest.Nil(t, err)
	workoutsEqual(t, workouts[1:2], res)

	res, err = logic.FindWorkoutsInDateRange(
		ctxt, workouts[0].WorkoutId.ClientEmail,
		time.Now().Add(23*time.Hour), time.Now().Add(72*time.Hour),
	)
	sbtest.Nil(t, err)
	workoutsEqual(t, workouts[1:], res)

	_, err = logic.FindWorkoutsInDateRange(
		ctxt, "asdf", time.Now().Add(1*time.Hour), time.Now(),
	)
	sbtest.ContainsError(
		t, types.CouldNotReadAllWorkoutsErr, err,
		`Start date \(.*\) must be before end date \(.*\)`,
	)

	n, err := logic.ReadNumClients(ctxt)
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, n)

	n, err = logic.ReadNumWorkoutsForClient(ctxt, "email@email.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, 3, n)
}

func workoutCreateDeletePhysData(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	err := logic.CreateClients(ctxt, types.Client{
		FirstName: "FName",
		LastName:  "LName",
		Email:     "email@email.com",
	})
	sbtest.Nil(t, err)

	workouts := []types.Workout{
		{
			WorkoutId: types.WorkoutId{
				ClientEmail:   "email@email.com",
				Session:       1,
				DatePerformed: time.Now(),
			},
			Exercises: []types.ExerciseData{
				{
					Name:   "Squat",
					Weight: 365,
					Sets:   5,
					Reps:   5,
					Effort: 10,
					PhysData: []types.Optional[types.PhysicsData]{
						{Present: true, Value: testPhysicsData1},
						{Present: true, Value: testPhysicsData2},
					},
				},
			},
		},
	}
	err = logic.CreateWorkouts(ctxt, workouts...)
	sbtest.Nil(t, err)

	n, err := logic.ReadNumWorkoutsForClient(ctxt, "email@email.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, n)

	err = logic.DeleteWorkouts(ctxt, workouts[0].WorkoutId)
	sbtest.Nil(t, err)

	_, err = logic.ReadWorkoutsById(ctxt, workouts[0].WorkoutId)
	sbtest.ContainsError(
		t, types.CouldNotReadAllWorkoutsErr, err,
		`Could not read entry with id '{ClientEmail:email@email.com Session:1 DatePerformed:.*}' \(Does id exist\?\)`,
	)

	err = logic.DeleteWorkouts(ctxt, types.WorkoutId{
		ClientEmail:   "asdf",
		Session:       1,
		DatePerformed: time.Now(),
	})
	sbtest.ContainsError(t, types.CouldNotDeleteAllWorkoutsErr, err)
	sbtest.ContainsError(
		t, types.CouldNotDeleteAllTrainingLogsErr, err,
		`Could not delete entry with id '{ClientEmail:asdf Session:1 DatePerformed:.*}' \(Does id exist\?\)`,
	)

	n, err = logic.ReadNumWorkoutsForClient(ctxt, "email@email.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, 0, n)
}

func workoutCreateDeleteBetweenDates(t *testing.T) {
	ctxt, cleanup := resetApp(t, context.Background())
	t.Cleanup(cleanup)

	err := logic.CreateClients(ctxt, types.Client{
		FirstName: "FName",
		LastName:  "LName",
		Email:     "email@email.com",
	})
	sbtest.Nil(t, err)

	workouts := []types.Workout{
		{
			WorkoutId: types.WorkoutId{
				ClientEmail:   "email@email.com",
				Session:       1,
				DatePerformed: time.Now(),
			},
			Exercises: []types.ExerciseData{
				{
					Name:   "Squat",
					Weight: 365,
					Sets:   5,
					Reps:   5,
					Effort: 10,
					PhysData: []types.Optional[types.PhysicsData]{
						{Present: true, Value: testPhysicsData1},
						{Present: true, Value: testPhysicsData2},
					},
				},
			},
		},
	}
	err = logic.CreateWorkouts(ctxt, workouts...)
	sbtest.Nil(t, err)

	n, err := logic.ReadNumWorkoutsForClient(ctxt, "email@email.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, 1, n)

	_, err = logic.DeleteWorkoutsInDateRange(
		ctxt, "email@email.com",
		time.Now().Add(24*time.Hour), time.Now().Add(48*time.Hour),
	)
	sbtest.ContainsError(t, types.CouldNotDeleteAllWorkoutsErr, err)
	sbtest.ContainsError(
		t, types.CouldNotDeleteAllTrainingLogsErr, err,
		`no rows in result set`,
	)

	res, err := logic.DeleteWorkoutsInDateRange(
		ctxt, "email@email.com",
		time.Now().Add(-1*time.Hour), time.Now().Add(24*time.Hour),
	)
	sbtest.Nil(t, err)
	sbtest.Eq(t, res, 1)

	_, err = logic.ReadWorkoutsById(ctxt, workouts[0].WorkoutId)
	sbtest.ContainsError(
		t, types.CouldNotReadAllWorkoutsErr, err,
		`Could not read entry with id '{ClientEmail:email@email.com Session:1 DatePerformed:.*}' \(Does id exist\?\)`,
	)

	_, err = logic.DeleteWorkoutsInDateRange(
		ctxt, "asdf", time.Now().Add(1*time.Hour), time.Now(),
	)
	sbtest.ContainsError(
		t, types.CouldNotDeleteAllWorkoutsErr, err,
		`Start date \(.*\) must be before end date \(.*\)`,
	)

	n, err = logic.ReadNumWorkoutsForClient(ctxt, "email@email.com")
	sbtest.Nil(t, err)
	sbtest.Eq(t, 0, n)
}
