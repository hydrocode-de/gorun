-- +goose Up
CREATE TABLE users (
    id text PRIMARY KEY,
    email text NOT NULL UNIQUE,
    password_hash text NOT NULL,
    is_admin boolean NOT NULL DEFAULT FALSE,
    created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_login datetime
);

CREATE TABLE refresh_tokens (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id TEXT NOT NULL,
    token text NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at DATETIME NOT NULL,
    is_revoked BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);

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
    error_message TEXT,
    user_id TEXT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE
);

-- +goose Down
DROP TABLE runs;
DROP TABLE refresh_tokens;
DROP TABLE users;