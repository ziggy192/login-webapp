FROM golang:1.19-alpine3.16 AS build

WORKDIR /src

COPY go.mod go.sum /src/

RUN go mod download

COPY . /src/

RUN go build -o api ./cmd/api

EXPOSE 8080

CMD ["./api"]
