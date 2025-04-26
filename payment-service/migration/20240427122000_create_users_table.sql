-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    balance NUMERIC(12, 2) NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE users;