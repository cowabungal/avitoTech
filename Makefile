.PHONY:

build:
	docker-compose build app

run:
	docker-compose up app

migrate-up:
	migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5434/postgres?sslmode=disable' up

migrate-down:
	migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5434/postgres?sslmode=disable' down

run-test:
	 go test ./... -cover