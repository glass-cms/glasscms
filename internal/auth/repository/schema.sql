CREATE TABLE tokens (
    id TEXT PRIMARY KEY,
    suffix TEXT NOT NULL,
    hash TEXT NOT NULL,
    create_time TIMESTAMP NOT NULL,
    expire_time TIMESTAMP NOT NULL
);

CREATE INDEX idx_tokens_hash ON tokens(hash);