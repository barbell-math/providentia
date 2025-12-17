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
SELECT
	providentia.client.first_name,
	providentia.client.last_name,
	providentia.client.email
FROM providentia.client 
JOIN UNNEST($1::TEXT[])
WITH ORDINALITY t(email, ord)
USING (email) 
ORDER BY ord;

-- name: FindClientsByEmail :many
-- Note: this is different from GetClientsByEmail because it returns a ordinal
-- value. The ordinal value allows for checking that the requested clients
-- existed in the database.
SELECT
	providentia.client.first_name,
	providentia.client.last_name,
	providentia.client.email,
	ord::INT8
FROM providentia.client 
JOIN UNNEST($1::TEXT[])
WITH ORDINALITY t(email, ord)
USING (email) 
ORDER BY ord;

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
