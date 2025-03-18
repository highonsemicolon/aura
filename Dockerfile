FROM golang:1.23.6 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make build

FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/tmp/main .
CMD ["./main"]