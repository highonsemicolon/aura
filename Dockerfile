
ARG GO_VERSION=1.25.0

FROM golang:${GO_VERSION}-bookworm AS build
WORKDIR /src

# Cache go mod first
COPY go.work* go.mod* ./
COPY apis/ ./apis/
COPY pkg/ ./pkg/
COPY services/ ./services/

# Enable module mode explicitly
ENV CGO_ENABLED=0 GO111MODULE=on

# Accept service arg
ARG SERVICE
ARG LDFLAGS=""
RUN test -n "$SERVICE" || (echo "SERVICE build-arg is required" && exit 1)

# Build
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -buildvcs=false -ldflags "$LDFLAGS" -o /out/${SERVICE} ./services/${SERVICE}

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
ARG SERVICE
COPY --from=build /out/${SERVICE} ./main
USER nonroot:nonroot
ENTRYPOINT ["/app/main"]
