-- name: CreatePhysicsData :one
INSERT INTO providentia.physics_data(
	path, bar_path_calc_id, bar_path_track_id,
	-- data
	time, position, velocity, acceleration, jerk,
	force, impulse, work, power,
	rep_splits,
	min_vel, max_vel,
	min_acc, max_acc,
	min_force, max_force,
	min_impulse, max_impulse,
	avg_work, min_work, max_work,
	avg_power, min_power, max_power
) VALUES (
	$1,
	(
		SELECT providentia.model.id FROM providentia.hyperparams
		JOIN providentia.model
			ON providentia.model.id = providentia.hyperparams.model_id
		WHERE providentia.model.name='BarPathCalc'
			AND providentia.hyperparams.version=sqlc.arg(bar_path_calc_params_version)
	),
	(
		SELECT providentia.model.id FROM providentia.hyperparams
		JOIN providentia.model
			ON providentia.model.id = providentia.hyperparams.model_id
		WHERE providentia.model.name='BarPathTracker'
			AND providentia.hyperparams.version=sqlc.arg(bar_path_tracker_params_version)
	),
	$2,
	$3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
	$18, $19, $20, $21, $22, $23, $24, $25
) RETURNING id;

-- name: GetTotalNumPhysicsEntriesForClient :one
SELECT COUNT(*) FROM providentia.physics_data
JOIN providentia.training_log
	ON providentia.physics_data.id = providentia.training_log.physics_id
JOIN providentia.client
	ON providentia.training_log.client_id = providentia.client.id
WHERE
	providentia.client.email = $1;
-- name: BulkCreateTrainingLogs :copyfrom
INSERT INTO providentia.training_log(
	exercise_id, client_id, physics_id,
	date_performed, weight, sets, reps, effort,
	inter_session_cntr, inter_workout_cntr
) VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
);

-- name: GetTotalNumTrainingLogEntriesForClient :one
SELECT COUNT(*) FROM providentia.training_log
JOIN providentia.client
	ON providentia.training_log.client_id = providentia.client.id
WHERE providentia.client.email = $1;

-- name: GetNumWorkoutsForClient :one
SELECT COUNT(*) FROM (
	SELECT date_performed, inter_session_cntr
	FROM providentia.training_log
	JOIN providentia.client
		ON providentia.training_log.client_id = providentia.client.id
	WHERE providentia.client.email = $1
	GROUP BY date_performed, inter_session_cntr
) AS result;

-- name: GetRawWorkoutData :many
SELECT
	providentia.exercise.name,
	providentia.training_log.weight,
	providentia.training_log.sets,
	providentia.training_log.reps,
	providentia.training_log.effort,
	providentia.physics_data.path,
	providentia.physics_data.time,
	providentia.physics_data.position
FROM providentia.training_log
JOIN providentia.exercise
	ON providentia.training_log.exercise_id=providentia.exercise.id
JOIN providentia.client
	ON providentia.training_log.client_id=providentia.client.id
LEFT JOIN providentia.physics_data
	ON providentia.training_log.physics_id=providentia.physics_data.id
WHERE
	providentia.client.email = $1 AND
	providentia.training_log.inter_session_cntr = $2 AND
	providentia.training_log.date_performed = $3
ORDER BY training_log.inter_workout_cntr ASC;

-- name: GetAllWorkoutData :many
-- TODO - figure out if ordinality trick can be used here - is there a way to
-- make the join on multiple columns?
SELECT
	providentia.exercise.name,
	providentia.training_log.weight,
	providentia.training_log.sets,
	providentia.training_log.reps,
	providentia.training_log.effort,
	providentia.training_log.volume,
	providentia.training_log.exertion,
	providentia.training_log.total_reps,
	providentia.physics_data.time,
	providentia.physics_data.position,
	providentia.physics_data.velocity,
	providentia.physics_data.acceleration,
	providentia.physics_data.jerk,
	providentia.physics_data.force,
	providentia.physics_data.impulse,
	providentia.physics_data.work,
	providentia.physics_data.power,
	providentia.physics_data.rep_splits,
	providentia.physics_data.min_vel,
	providentia.physics_data.max_vel,
	providentia.physics_data.min_acc,
	providentia.physics_data.max_acc,
	providentia.physics_data.min_force,
	providentia.physics_data.max_force,
	providentia.physics_data.min_impulse,
	providentia.physics_data.max_impulse,
	providentia.physics_data.avg_work,
	providentia.physics_data.min_work,
	providentia.physics_data.max_work,
	providentia.physics_data.avg_power,
	providentia.physics_data.min_power,
	providentia.physics_data.max_power
FROM providentia.training_log
JOIN providentia.exercise
	ON providentia.training_log.exercise_id=providentia.exercise.id
JOIN providentia.client
	ON providentia.training_log.client_id=providentia.client.id
LEFT JOIN providentia.physics_data
	ON providentia.training_log.physics_id=providentia.physics_data.id
WHERE
	providentia.client.email = $1 AND
	providentia.training_log.inter_session_cntr = $2 AND
	providentia.training_log.date_performed = $3
ORDER BY training_log.inter_workout_cntr ASC;

-- name: GetAllWorkoutDataBetweenDates :many
SELECT
	providentia.exercise.name,
	providentia.training_log.weight,
	providentia.training_log.sets,
	providentia.training_log.reps,
	providentia.training_log.effort,
	providentia.training_log.volume,
	providentia.training_log.exertion,
	providentia.training_log.total_reps,
	providentia.training_log.date_performed,
	providentia.training_log.inter_session_cntr,
	providentia.physics_data.time,
	providentia.physics_data.position,
	providentia.physics_data.velocity,
	providentia.physics_data.acceleration,
	providentia.physics_data.jerk,
	providentia.physics_data.force,
	providentia.physics_data.impulse,
	providentia.physics_data.work,
	providentia.physics_data.power,
	providentia.physics_data.rep_splits,
	providentia.physics_data.min_vel,
	providentia.physics_data.max_vel,
	providentia.physics_data.min_acc,
	providentia.physics_data.max_acc,
	providentia.physics_data.min_force,
	providentia.physics_data.max_force,
	providentia.physics_data.min_impulse,
	providentia.physics_data.max_impulse,
	providentia.physics_data.avg_work,
	providentia.physics_data.min_work,
	providentia.physics_data.max_work,
	providentia.physics_data.avg_power,
	providentia.physics_data.min_power,
	providentia.physics_data.max_power
FROM providentia.training_log
JOIN providentia.exercise
	ON providentia.training_log.exercise_id=providentia.exercise.id
JOIN providentia.client
	ON providentia.training_log.client_id=providentia.client.id
LEFT JOIN providentia.physics_data
	ON providentia.training_log.physics_id=providentia.physics_data.id
WHERE
	providentia.client.email = $1 AND
	providentia.training_log.date_performed BETWEEN @start::DATE AND @ending::DATE
ORDER BY 
	training_log.date_performed ASC,
	training_log.inter_session_cntr ASC,
	training_log.inter_workout_cntr ASC;

-- name: DeleteWorkout :one
WITH deleted_exercises AS (
	DELETE FROM providentia.training_log
	USING providentia.client
	WHERE
		providentia.client.id = providentia.training_log.client_id AND
		providentia.client.email = $1 AND
		providentia.training_log.inter_session_cntr = $2 AND
		providentia.training_log.date_performed = $3
	RETURNING providentia.training_log.id
) SELECT COUNT(*) FROM deleted_exercises;

-- name: DeleteWorkoutsBetweenDates :one
WITH deleted_exercises AS (
	DELETE FROM providentia.training_log
	USING providentia.client
	WHERE
		providentia.client.id = providentia.training_log.client_id AND
		providentia.client.email = $1 AND
		providentia.training_log.date_performed BETWEEN @start::DATE AND @ending::DATE
	RETURNING providentia.training_log.id
) SELECT COUNT(*) FROM deleted_exercises;
