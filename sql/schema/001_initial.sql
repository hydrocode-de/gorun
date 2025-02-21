-- +goose Up
CREATE TABLE runs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name text NOT NULL,
    title text NOT NULL,
    description text NOT NULL,
    docker_image text NOT NULL,
    mounts text NOT NULL,
    parameters text NOT NULL,
    data text NOT NULL,
    created_at DATETIME NOT NULL,
    started_at DATETIME,
    finished_at DATETIME,
    status text not null default 'pending',
    has_errored BOOLEAN NOT NULL DEFAULT FALSE,
    error_message TEXT
);

-- +goose Down
DROP TABLE runs;