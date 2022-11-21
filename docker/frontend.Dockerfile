FROM golang:1.19-alpine3.16 AS build

WORKDIR /src

COPY go.mod go.sum /src/

RUN go mod download

COPY . /src/

RUN go build -o frontend ./cmd/frontend

EXPOSE 9090

CMD ["./frontend"]
