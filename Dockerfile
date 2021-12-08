FROM golang:latest

RUN go version

COPY . /avitoTech/
WORKDIR /avitoTech/

# build go app
RUN go mod download
RUN GOOS=linux go build -o app ./cmd/main.go

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN sed -i -e 's/\r$//' *.sh
RUN chmod +x wait-for-postgres.sh

CMD ["./app"]