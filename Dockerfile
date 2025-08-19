FROM golang:latest

WORKDIR /app

RUN apt-get update && apt-get install -y iputils-ping netcat-openbsd


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux go build -o main ./cmd

CMD ["./main"]


