.PHONY:

build:
	go build -o ./.bin/main cmd/main.go

run: build
	./.bin/main

clear:
	rm -rf ./.bin

migrate-up:
	migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5433/postgres?sslmode=disable' up

migrate-down:
	migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5433/postgres?sslmode=disable' down

run-test:
	 go test ./... -cover