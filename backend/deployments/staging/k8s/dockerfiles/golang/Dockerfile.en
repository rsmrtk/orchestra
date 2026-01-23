# syntax=docker/dockerfile:1.4
# Production-grade Go Dockerfile (multi-stage, static binary, non-root)

FROM golang:1.25.4-alpine AS builder
RUN apk add --no-cache ca-certificates git tzdata
WORKDIR /build

ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download && go mod verify

COPY . .

ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build -ldflags="-w -s" -trimpath -o /app/server ./cmd/server

RUN ldd /app/server 2>&1 | grep -q "not a dynamic executable" || \
    (echo "ERROR: Binary is not static" && exit 1)

FROM gcr.io/distroless/static-debian12:nonroot
USER nonroot:nonroot
WORKDIR /app

COPY --from=builder --chown=nonroot:nonroot /app/server /app/server

EXPOSE 8080
ENV PORT=8080

ENTRYPOINT ["/app/server"]
