FROM golang:1.23-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o notifier

FROM debian:bullseye-slim

RUN apt-get update && \
    apt-get install -y libc6 ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/notifier /app/notifier

COPY ./config /app/config

RUN chmod +x /app/notifier && mkdir /app/sqlite_data

CMD ["/app/notifier"]
