.PHONY:

build:
	go build -o ./.bin/main cmd/main.go

run: build
	./.bin/main

clear:
	rm -rf ./.bin