-- +goose Up
-- +goose StatementBegin
CREATE TABLE bookings (
    id  SERIAL PRIMARY KEY,
    user_id int NOT NULL,
    ticket_id int NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bookings;
-- +goose StatementEnd
