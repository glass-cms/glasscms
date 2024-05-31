CREATE TABLE items (
    id INTEGER PRIMARY KEY AUTOINCREMENT, 
    name TEXT NOT NULL,
    path TEXT NOT NULL,
    content TEXT,
    hash TEXT,
    create_time TIMESTAMP NOT NULL,
    update_time TIMESTAMP NOT NULL,
    properties JSON
);