-- name: BulkCreateExerciseFocusWithID :copyfrom
-- This query is used for initilization by the migrations. The
-- UpdateExerciseFocusSerialCount query will need to be run after this to update
-- the serial counter.
INSERT INTO providentia.exercise_focus (id, focus) VALUES ($1, $2);

-- name: UpdateExerciseFocusSerialCount :exec
SELECT SETVAL(
	pg_get_serial_sequence('providentia.exercise_focus', 'id'),
	(SELECT MAX(id) FROM providentia.exercise_focus) + 1
);



-- name: BulkCreateExerciseKindWithID :copyfrom
-- This query is used for initilization by the migrations. The
-- UpdateExerciseKindSerialCount query will need to be run after this to update
-- the serial counter.
INSERT INTO providentia.exercise_kind (id, kind, description) VALUES ($1, $2, $3);

-- name: UpdateExerciseKindSerialCount :exec
SELECT SETVAL(
	pg_get_serial_sequence('providentia.exercise_kind', 'id'),
	(SELECT MAX(id) FROM providentia.exercise_kind) + 1
);



-- name: BulkCreateClients :copyfrom
INSERT INTO providentia.client (first_name, last_name, email) VALUES ($1, $2, $3);

-- name: EnsureClientsExist :exec
INSERT INTO providentia.client (first_name, last_name, email)
SELECT
	UNNEST(@first_names::TEXT[]),
	UNNEST(@last_names::TEXT[]),
	UNNEST(@emails::TEXT[])
ON CONFLICT (first_name, last_name, email) DO NOTHING;

-- name: GetNumClients :one
SELECT COUNT(*) FROM providentia.client;

-- name: GetClientIdByEmail :one
SELECT id FROM providentia.client WHERE email = $1;

-- name: GetClientsByEmail :many
SELECT first_name, last_name, email
FROM providentia.client WHERE email = ANY($1::TEXT[]);

-- name: ClientExists :one
SELECT EXISTS(SELECT 1 FROM providentia.client WHERE email = $1);

-- name: UpdateClientByEmail :exec
UPDATE providentia.client SET first_name=$1, last_name=$2
WHERE providentia.client.email=$3;

-- name: DeleteClientsByEmail :one
WITH deleted_clients AS (
    DELETE FROM providentia.client
    WHERE email = ANY($1::TEXT[])
    RETURNING id
) SELECT COUNT(*) FROM deleted_clients;



-- name: BulkCreateExerciseWithID :copyfrom
-- This query is used for initilization by the migrations. The
-- UpdateExerciseSerialCount query will need to be run after this to update
-- the serial counter.
INSERT INTO providentia.exercise(
	id, name, kind_id, focus_id
) VALUES ($1, $2, $3, $4);

-- name: UpdateExerciseSerialCount :exec
SELECT SETVAL(
	pg_get_serial_sequence('providentia.exercise', 'id'),
	(SELECT MAX(id) FROM providentia.exercise) + 1
);

-- name: BulkCreateExercises :copyfrom
INSERT INTO providentia.exercise (name, kind_id, focus_id) VALUES ($1, $2, $3);

-- name: EnsureExercisesExist :exec
INSERT INTO providentia.exercise (name, kind_id, focus_id)
SELECT
	UNNEST(@names::TEXT[]),
	UNNEST(@kinds::INT4[]),
	UNNEST(@focuses::INT4[])
ON CONFLICT (name, kind_id, focus_id) DO NOTHING;

-- name: GetNumExercises :one
SELECT COUNT(*) FROM providentia.exercise;

-- name: GetExerciseIdByName :one
SELECT id FROM providentia.exercise WHERE name = $1;

-- name: GetExercisesByName :many
SELECT name, kind_id, focus_id
FROM providentia.exercise WHERE name = ANY($1::TEXT[]);

-- name: UpdateExerciseByName :exec
UPDATE providentia.exercise SET kind_id=$2, focus_id=$3
WHERE providentia.exercise.name=$1;

-- name: DeleteExercisesByName :one
WITH deleted_exercises AS (
	DELETE FROM providentia.exercise
	WHERE name = ANY($1::TEXT[])
	RETURNING id
) SELECT COUNT(*) FROM deleted_exercises;



-- name: BulkCreateModelsWithID :copyfrom
-- This query is used for initilization by the migrations. The
-- UpdateModelSerialCount query will need to be run after this to update
-- the serial counter.
INSERT INTO providentia.model (id, name, description) VALUES ($1, $2, $3);

-- name: UpdateModelSerialCount :exec
SELECT SETVAL(
	pg_get_serial_sequence('providentia.model', 'id'),
	(SELECT MAX(id) FROM providentia.model) + 1
);



-- name: BulkCreateHyperparamsWithID :copyfrom
-- This query is used for initilization by the migrations. The
-- UpdateModelSerialCount query will need to be run after this to update
-- the serial counter.
INSERT INTO providentia.hyperparams (
	id, model_id, version, params
) VALUES ($1, $2, $3, $4); 

-- name: UpdateHyperparamsSerialCount :exec
SELECT SETVAL(
	pg_get_serial_sequence('providentia.hyperparams', 'id'),
	(SELECT MAX(id) FROM providentia.hyperparams) + 1
);

-- name: BulkCreateHyperparams :copyfrom
INSERT INTO providentia.hyperparams (
	model_id, version, params
) VALUES ($1, $2, $3);

-- name: GetNumHyperparams :one
SELECT COUNT(*) FROM providentia.hyperparams;

-- name: GetNumHyperparamsFor :one
SELECT COUNT(*) FROM providentia.hyperparams
WHERE providentia.hyperparams.model_id=$1;

-- name: GetHyperparamsByVersionFor :many
SELECT version, params FROM providentia.hyperparams
WHERE model_id=$1 AND version = ANY(@versions::INT4[]);

-- name: DeleteHyperparamsByVersionFor :one
WITH deleted_hyperparams AS (
	DELETE FROM providentia.hyperparams
	WHERE model_id=$1 AND version = ANY(@versions::INT4[])
	RETURNING id
) SELECT COUNT(*) FROM deleted_hyperparams;



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

-- name: GetAllWorkoutData :many
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



-- name: CreatePhysicsData :one
INSERT INTO providentia.physics_data(
	path,
	bar_path_calc_id, bar_path_track_id,
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
			AND providentia.hyperparams.version=$2	-- TODO - sqlc arg name
	),
	(
		SELECT providentia.model.id FROM providentia.hyperparams
		JOIN providentia.model
			ON providentia.model.id = providentia.hyperparams.model_id
		WHERE providentia.model.name='BarPathTracker'
			AND providentia.hyperparams.version=$3	-- TODO - sqlc arg name
	),
	$4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19,
	$20, $21, $22, $23, $24, $25, $26, $27
) RETURNING id;

-- name: GetTotalNumPhysicsEntriesForClient :one
SELECT COUNT(*) FROM providentia.physics_data
JOIN providentia.training_log
	ON providentia.physics_data.id = providentia.training_log.physics_id
JOIN providentia.client
	ON providentia.training_log.client_id = providentia.client.id
WHERE
	providentia.client.email = $1;
