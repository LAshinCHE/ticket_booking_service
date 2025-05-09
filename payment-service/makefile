.PHONY: migrate

migrate:
	goose -dir ./migrations postgres "postgresql://payment:paymentpass@localhost:5434/paymentdb?sslmode=disable" up

fill-data:
	psql "postgresql://payment:paymentpass@localhost:5434/paymentdb?sslmode=disable" -f user_data.sql