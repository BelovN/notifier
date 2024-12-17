FROM golang:1.23-bullseye AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o notifier

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/notifier /root/notifier

RUN chmod +x /root/notifier && mkdir /root/sqlite_data

CMD ["/root/notifier"]
