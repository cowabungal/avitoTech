FROM golang:1.17-alpine AS builder

RUN go version

COPY . /avitoTech/
WORKDIR /avitoTech/

RUN go mod download
RUN GOOS=linux go build -o ./.bin/main ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=0 /avitoTech/.bin/main .
COPY --from=0 /avitoTech/ .

CMD ["./main"]