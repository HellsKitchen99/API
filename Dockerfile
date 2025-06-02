FROM golang:latest AS builder

WORKDIR /apiWithDataBase

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY .. ./

RUN go build -o ./app .
CMD ["./app"]


