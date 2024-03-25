# http/Dockerfile
FROM golang:1.21

WORKDIR /app

COPY ./. .

RUN go mod download


RUN go build -o client

CMD ["./client"]

