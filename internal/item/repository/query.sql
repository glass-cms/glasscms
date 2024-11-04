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

-- name: UpsertItem :one
INSERT INTO items (
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
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(name) DO UPDATE SET
    display_name = excluded.display_name,
    create_time = excluded.create_time,
    update_time = excluded.update_time,
    delete_time = excluded.delete_time,
    hash = excluded.hash,
    content = excluded.content,
    properties = excluded.properties,
    metadata = excluded.metadata
RETURNING name, display_name, create_time, update_time, delete_time, hash, content, properties, metadata;
