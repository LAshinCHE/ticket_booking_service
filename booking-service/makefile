.PHONY: swagger migrate

swagger:
	swag init -dir ./ -generalInfo cmd/main.go ./internal/api/http/booking.go -o ./docs

migrate:
	goose -dir ./migrations postgres "postgresql://booking:bookingpass@localhost:5431/bookingdb?sslmode=disable" up
