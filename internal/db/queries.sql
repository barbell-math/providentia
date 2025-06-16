-- name: BulkCreateExerciseFocusWithID :copyfrom
INSERT INTO providentia.exercise_focus (id, focus) VALUES ($1, $2);

-- name: BulkCreateExerciseKindWithID :copyfrom
INSERT INTO providentia.exercise_kind (id, kind, description) VALUES ($1, $2, $3);

-- name: BulkCreateModels :copyfrom
INSERT INTO providentia.model (id, name, description) VALUES ($1, $2, $3);

-- name: BulkCreateExerciseWithID :copyfrom
INSERT INTO providentia.exercise(
	id, name, kind_id, focus_id
) VALUES ($1, $2, $3, $4);

-- name: BulkCreateVideoDataWithID :copyfrom
INSERT INTO providentia.video_data (
	id, path, position, velocity, acceleration, force, impulse
) VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: BulkCreateClients :copyfrom
INSERT INTO providentia.client (first_name, last_name, email) VALUES ($1, $2, $3);

-- name: BulkCreateTraingLog :copyfrom
INSERT INTO providentia.training_log(
	exercise_id, exercise_kind_id, exercise_focus_id, client_id, video_id,
	date_performed, weight, sets, reps, effort,
	inter_session_cntr, inter_workout_cntr
) VALUES (
	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
);

-- name: BulkCreateModelStates :copyfrom
INSERT INTO providentia.model_state(
	client_id, training_log_id, model_id,
	v1, v2, v3, v4, v5, v6, v7, v8, v9, v10,
	time_frame, mse, pred_weight
) values (
	$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
);

-- name: GetClientIDFromEmail :one
SELECT ID FROM providentia.client WHERE email=$1;

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
	providentia.training_log.exercise_kind_id,
	providentia.training_log.exercise_focus_id,
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
