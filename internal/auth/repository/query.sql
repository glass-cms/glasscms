-- name: CreateToken :exec
INSERT INTO tokens (id, suffix, hash, create_time, expire_time) VALUES (?, ?, ?, CURRENT_TIMESTAMP, ?);

-- name: GetToken :one
SELECT id, suffix, hash, create_time, expire_time FROM tokens WHERE hash = ?;

-- name: DeleteToken :exec
DELETE FROM tokens WHERE id = ?; 