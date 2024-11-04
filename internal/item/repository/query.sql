-- name: CreateItem :one
INSERT INTO
    items (
        name,
        display_name,
        create_time,
        update_time,
        delete_time,
        hash,
        content,
        properties,
        metadata
    )
VALUES
    (?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;

-- name: UpdateItem :one
UPDATE
    items
SET
    name = ?,
    display_name = ?,
    update_time = ?,
    hash = ?,
    content = ?,
    properties = ?,
    metadata = ?
WHERE
    name = ?
    AND delete_time IS NULL
RETURNING *;

-- name: DeleteItem :exec
UPDATE
    items
SET
    delete_time = ?
WHERE
    name = ?;

-- name: GetItem :one
SELECT
    *
FROM
    items
WHERE
    name = ?
    AND delete_time IS NULL;

-- name: ListItems :many
SELECT
    *
FROM
    items
WHERE
    delete_time IS NULL;