CREATE TABLE items (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    create_time TIMESTAMP NOT NULL,
    update_time TIMESTAMP NOT NULL,
    hash TEXT, 
    name TEXT NOT NULL,
    path TEXT NOT NULL,
    content TEXT,
    properties JSON
);