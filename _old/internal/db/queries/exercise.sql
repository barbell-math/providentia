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
SELECT
	providentia.exercise.name,
	providentia.exercise.kind_id,
	providentia.exercise.focus_id
FROM providentia.exercise 
JOIN UNNEST($1::TEXT[])
WITH ORDINALITY t(name, ord)
USING (name)
ORDER BY ord;

-- name: FindExercisesByName :many
-- Note: this is different from GetExercisesByName because it returns a ordinal
-- value. The ordinal value allows for checking that the requested exercises
-- existed in the database.
SELECT
	providentia.exercise.name,
	providentia.exercise.kind_id,
	providentia.exercise.focus_id,
	ord::INT8
FROM providentia.exercise 
JOIN UNNEST($1::TEXT[])
WITH ORDINALITY t(name, ord)
USING (name) 
ORDER BY ord;

-- name: UpdateExerciseByName :exec
UPDATE providentia.exercise SET kind_id=$2, focus_id=$3
WHERE providentia.exercise.name=$1;

-- name: DeleteExercisesByName :one
WITH deleted_exercises AS (
	DELETE FROM providentia.exercise
	WHERE name = ANY($1::TEXT[])
	RETURNING id
) SELECT COUNT(*) FROM deleted_exercises;
