-- +goose Up
-- +goose StatementBegin
CREATE TABLE tickets (
    id int PRIMARY KEY,
    price DOUBLE PRECISION NOT NULL,
    available BOOLEAN NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tickets;
-- +goose StatementEnd
