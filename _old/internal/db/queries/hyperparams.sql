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
SELECT
	providentia.hyperparams.version,
	providentia.hyperparams.params
FROM providentia.hyperparams
JOIN UNNEST(@versions::INT4[])
WITH ORDINALITY t(version, ord)
USING (version)
WHERE model_id=$1
ORDER BY ord;

-- name: FindHyperparamsByVersionFor :many
-- Note: this is different from GetHyperparamsByVersionFor because it returns a
-- ordinal value. The ordinal value allows for checking that the requested
-- hyperparameters existed in the database.
SELECT
	providentia.hyperparams.version,
	providentia.hyperparams.params,
	ord::INT8
FROM providentia.hyperparams
JOIN UNNEST(@versions::INT4[])
WITH ORDINALITY t(version, ord)
USING (version)
WHERE model_id=$1
ORDER BY ord;

-- name: DeleteHyperparamsByVersionFor :one
WITH deleted_hyperparams AS (
	DELETE FROM providentia.hyperparams
	WHERE model_id=$1 AND version = ANY(@versions::INT4[])
	RETURNING id
) SELECT COUNT(*) FROM deleted_hyperparams;
