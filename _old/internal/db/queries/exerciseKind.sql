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
