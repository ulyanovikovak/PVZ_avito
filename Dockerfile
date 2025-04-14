FROM golang:1.23.1 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o pvz-server ./cmd/pvz-server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/pvz-server .
COPY migrations ./migrations

EXPOSE 8080

CMD ["./pvz-server"]
