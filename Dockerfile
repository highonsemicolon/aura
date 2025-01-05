FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o aura ./cmd

FROM golang:1.23-alpine AS dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy 
RUN go install github.com/air-verse/air@latest
COPY . .
CMD ["air"]

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/aura .
EXPOSE 8080
CMD ["./aura"]