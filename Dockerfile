ARG GO_VERSION=1.25.0

FROM golang:${GO_VERSION}-bookworm AS build
WORKDIR /src

# Install buf
RUN go install github.com/bufbuild/buf/cmd/buf@latest

# Copy proto config + proto files first (better caching)
COPY buf.* ./
COPY apis/ ./apis/

# Generate protobuf code
RUN /go/bin/buf generate

# Copy go modules and source
COPY go.work* go.mod* ./
COPY pkg/ ./pkg/
COPY services/ ./services/

# Enable module mode
ENV CGO_ENABLED=0 GO111MODULE=on

# Accept service arg
ARG SERVICE
ARG LDFLAGS=""
RUN test -n "$SERVICE" || (echo "SERVICE build-arg is required" && exit 1)

# Build binaries
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -buildvcs=false -ldflags "$LDFLAGS" -o /out/${SERVICE} ./services/${SERVICE}

# Final image
FROM gcr.io/distroless/static:nonroot
WORKDIR /app
ARG SERVICE
COPY --from=build /out/${SERVICE} ./main
USER nonroot:nonroot
ENTRYPOINT ["/app/main"]
