-- name: BulkCreateExerciseFocus :copyfrom
INSERT INTO providentia.ExerciseFocus (ID, Focus) VALUES ($1, $2);

-- name: BulkCreateExerciseKind :copyfrom
INSERT INTO providentia.ExerciseKind (ID, Kind, Description) VALUES ($1, $2, $3);

-- name: BulkCreateExercise :copyfrom
INSERT INTO providentia.Exercise(
	ID, Name, KindID, FocusID
) VALUES ($1, $2, $3, $4);
