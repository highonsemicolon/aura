FROM golang:1.23.6-alpine AS dev
WORKDIR /app

RUN apk add --no-cache git make \
    && go install github.com/air-verse/air@latest

COPY go.mod go.sum ./

COPY . .
CMD ["make", "run-hot"]