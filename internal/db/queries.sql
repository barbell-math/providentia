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

-- name: GetNumClients :one
SELECT COUNT(*) FROM providentia.client;

-- name: GetClientsByEmail :many
SELECT first_name, last_name, email
FROM providentia.client WHERE email = ANY($1::text[]);

-- name: GetFullClientByEmail :one
SELECT id, first_name, last_name, email FROM providentia.client WHERE email = $1;

-- name: UpdateClientByEmail :exec
UPDATE providentia.client SET first_name=$1, last_name=$2
WHERE providentia.client.email=$3;

-- name: DeleteClientsByEmail :one
WITH deleted_clients AS (
    DELETE FROM providentia.client
    WHERE email = ANY($1::text[])
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

-- name: GetNumExercises :one
SELECT COUNT(*) FROM providentia.exercise;

-- name: GetExercisesByName :many
SELECT name, kind_id, focus_id
FROM providentia.exercise WHERE name = ANY($1::text[]);

-- name: GetFullExerciseByName :one
SELECT id, name, kind_id, focus_id FROM providentia.exercise WHERE name = $1;

-- name: UpdateExerciseByName :exec
UPDATE providentia.exercise SET kind_id=$2, focus_id=$3
WHERE providentia.exercise.name=$1;

-- name: DeleteExercisesByName :one
WITH deleted_exercises AS (
    DELETE FROM providentia.exercise
    WHERE name = ANY($1::text[])
    RETURNING id
) SELECT COUNT(*) FROM deleted_exercises;



-- name: BulkCreateModelsWithID :copyfrom
-- This query is used for initilization by the migrations. The
-- UpdateModelSerialCount query will need to be run after this to update
-- the serial counter.
INSERT INTO providentia.model (id, name, description) VALUES ($1, $2, $3);

-- name: UpdateModelSerialCount :exec
SELECT SETVAL(
	pg_get_serial_sequence('providentia.exercise', 'id'),
	(SELECT MAX(id) FROM providentia.exercise) + 1
);



-- name: BulkCreateTrainingLogs :copyfrom
INSERT INTO providentia.training_log(
	exercise_id, client_id, physics_id,
	date_performed, weight, sets, reps, effort,
	inter_session_cntr, inter_workout_cntr
) VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10
);

-- name: GetTotalNumExercisesForClient :one
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
	providentia.physics_data.work
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


-- name: CreatePhysicsData :one
INSERT INTO providentia.physics_data(
	path,
	time, position, velocity, acceleration, jerk,
	force, impulse, work
) VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING id;

-- name: GetTotalNumPhysicsEntriesForClient :one
SELECT COUNT(*) FROM providentia.physics_data
JOIN providentia.training_log
	ON providentia.physics_data.id = providentia.training_log.physics_id
JOIN providentia.client
	ON providentia.training_log.client_id = providentia.client.id
WHERE
	providentia.client.email = $1;

----- OLD ----------------------------------------------------------------------
-- name: BulkCreateModelStates :copyfrom
INSERT INTO providentia.model_state(
	client_id, training_log_id, model_id,
	v1, v2, v3, v4, v5, v6, v7, v8, v9, v10,
	time_frame, mse, pred_weight
) values (
	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
);




-- name: GetAllClientsTrainingLogData :many
SELECT
	providentia.client.email,
	providentia.exercise.name,
	providentia.training_log.date_performed,
	providentia.training_log.inter_session_cntr,
	providentia.training_log.weight,
	providentia.training_log.sets,
	providentia.training_log.reps,
	providentia.training_log.effort,
	providentia.training_log.volume,
	providentia.training_log.exertion,
	providentia.training_log.total_reps
FROM providentia.training_log
JOIN providentia.exercise
	ON providentia.training_log.exercise_id=providentia.exercise.id
JOIN providentia.client
	ON providentia.training_log.client_id=providentia.client.id
ORDER BY
	-- These cannot be labeled with providentia.training_log because you will
	-- get a `column reference "" not found` error.
	client.id DESC,
	training_log.date_performed DESC,
	training_log.id DESC
LIMIT $1;

-- name: GetClientTrainingLogData :many
SELECT
	providentia.exercise.name,
	providentia.training_log.date_performed,
	providentia.training_log.inter_session_cntr,
	providentia.training_log.weight,
	providentia.training_log.sets,
	providentia.training_log.reps,
	providentia.training_log.effort
FROM providentia.training_log
JOIN providentia.exercise
	ON providentia.training_log.exercise_id=providentia.exercise.id
JOIN providentia.client
	ON providentia.training_log.client_id=providentia.client.id
WHERE providentia.client.email=$1
ORDER BY
	-- These cannot be labeled with providentia.training_log because you will
	-- get a `column reference "" not found` error.
	training_log.date_performed DESC, training_log.id DESC
LIMIT $2;

-- name: GetExerciseIDs :one
SELECT
	providentia.exercise.id AS exercise_id,
	providentia.exercise_kind.id AS kind_id,
	providentia.exercise_focus.id AS focus_id
FROM providentia.exercise
JOIN providentia.exercise_kind
	ON providentia.exercise.kind_id=providentia.exercise_kind.id
JOIN providentia.exercise_focus
	ON providentia.exercise.focus_id=providentia.exercise_focus.id
WHERE name=$1;

-- name: ClientLastWorkoutDate :one
SELECT date_performed FROM providentia.training_log
WHERE client_id=$1
ORDER BY date_performed DESC LIMIT 1;

-- name: ClientTrainingLogDataDateRangeAscending :many
SELECT
	providentia.training_log.id,
	providentia.training_log.exercise_id,
	(sqlc.arg(date_performed)::date-providentia.training_log.date_performed) AS days_since,
	providentia.training_log.weight,
	providentia.training_log.sets,
	providentia.training_log.reps,
	providentia.training_log.effort,
	providentia.training_log.inter_session_cntr,
	providentia.training_log.inter_workout_cntr
FROM providentia.training_log
WHERE providentia.training_log.client_id=$1
	AND providentia.training_log.date_performed<sqlc.arg(date_performed)
ORDER BY
	-- These cannot be labeled with providentia.training_log because you will
	-- get a `column reference "" not found` error.
	date_performed ASC, id ASC;
