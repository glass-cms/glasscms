-- +goose Up
CREATE TABLE items (
    name TEXT PRIMARY KEY, 
    display_name TEXT NOT NULL,
    create_time TIMESTAMP NOT NULL,
    update_time TIMESTAMP NOT NULL,
    delete_time TIMESTAMP,
    hash TEXT, 
    content TEXT,
    properties JSON,
    metadata JSON
);

CREATE INDEX items_delete_time ON items(delete_time);

-- +goose Down
DROP TABLE items;