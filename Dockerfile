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
RUN apk add --no-cache bash git curl
RUN curl -fLo /usr/local/bin/air https://github.com/air-verse/air/releases/download/v1.61.5/air_1.61.5_linux_arm64 \
    && chmod +x /usr/local/bin/air
COPY . .
CMD ["air"]

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/aura .
EXPOSE 8080
CMD ["./aura"]