-- +goose Up
CREATE TABLE items (
    uid TEXT PRIMARY KEY,
    create_time TIMESTAMP NOT NULL,
    update_time TIMESTAMP NOT NULL,
    hash TEXT, 
    name TEXT NOT NULL,
    path TEXT NOT NULL,
    content TEXT,
    properties JSON
);

-- +goose Down
DROP TABLE items;