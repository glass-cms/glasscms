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