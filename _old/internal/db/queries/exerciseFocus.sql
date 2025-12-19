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
