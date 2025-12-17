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
