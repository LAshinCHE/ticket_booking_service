.PHONY
migrate:
	goose -dir ./migrations postgres "postgresql://ticket:ticketpass@localhost:5433/ticketdb?sslmode=disable" up

fill-data:
	psql "postgresql://ticket:ticketpass@localhost:5433/ticketdb?sslmode=disable" -f ticket_data.sql

swagger-init:
	swag init -dir ./ -generalInfo ./internal/api/http/ticket.go -o ./docs