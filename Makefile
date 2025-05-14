up:
	docker-compose up -d

down:
	docker-compose down

db-up:
	goose -dir ./booking-service/migrations postgres "postgresql://booking:bookingpass@localhost:5431/bookingdb?sslmode=disable" up
	goose -dir ./ticket-service/migrations postgres "postgresql://ticket:ticketpass@localhost:5433/ticketdb?sslmode=disable" up
	goose -dir ./payment-service/migrations postgres "postgresql://payment:paymentpass@localhost:5434/paymentdb?sslmode=disable" up

booking-db-up:
	goose -dir ./booking-service/migrations postgres "postgresql://booking:bookingpass@localhost:5431/bookingdb?sslmode=disable" up

booking-db-down:
	goose -dir ./booking-service/migrations postgres "postgresql://booking:bookingpass@localhost:5431/bookingdb?sslmode=disable" down

ticket-db-up:
	goose -dir ./ticket-service/migrations postgres "postgresql://ticket:ticketpass@localhost:5433/ticketdb?sslmode=disable" up

ticket-db-down:
	goose -dir ./ticket-service/migrations postgres "postgresql://ticket:ticketpass@localhost:5433/ticketdb?sslmode=disable" down

payment-db-up:
	goose -dir ./payment-service/migrations postgres "postgresql://payment:paymentpass@localhost:5434/paymentdb?sslmode=disable" up

payment-db-down:
	goose -dir ./payment-service/migrations postgres "postgresql://payment:paymentpass@localhost:5434/paymentdb?sslmode=disable" down

seed:
	psql "postgresql://payment:paymentpass@localhost:5434/paymentdb?sslmode=disable" -f ./generate-data/user_data.sql
	psql "postgresql://ticket:ticketpass@localhost:5433/ticketdb?sslmode=disable" -f ./generate-data/ticket_data.sql

test:
	k6 run ./generate-data/test-service/loadtest.js
