-- +goose Up
CREATE TABLE IF NOT EXISTS example (
    id   INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS example;
