# syntax=docker/dockerfile:1.4
# Production-grade Go Dockerfile
# Purpose: build a static Go binary and run it in a minimal, non-root image.
# How it works: multi-stage build -> compile in builder -> copy binary to runtime.

# ===== Stage 1: Builder =====
FROM golang:1.25.4-alpine AS builder

# Install minimal build tools. git is needed for private modules.
RUN apk add --no-cache ca-certificates git tzdata

# Use a clean build directory inside the image.
WORKDIR /build

# Configure module proxy for reliable downloads.
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org

# Copy only go.mod/go.sum first to maximize cache hits.
COPY go.mod go.sum ./

# Download deps with cache mounts for fast rebuilds.
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download && go mod verify

# Copy the rest of the source code.
COPY . .

# Build a static Linux binary.
# TARGETOS/TARGETARCH allow cross-platform builds from CI.
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 \
    GOOS=$TARGETOS \
    GOARCH=$TARGETARCH \
    go build -ldflags="-w -s" -trimpath -o /app/server ./cmd/server

# Ensure the binary is fully static (required for distroless static).
RUN ldd /app/server 2>&1 | grep -q "not a dynamic executable" || \
    (echo "ERROR: Binary is not static" && exit 1)

# ===== Stage 2: Runtime =====
FROM gcr.io/distroless/static-debian12:nonroot

# Run as non-root for security. distroless provides UID 65532.
USER nonroot:nonroot
WORKDIR /app

# Copy the compiled binary from the builder stage.
COPY --from=builder --chown=nonroot:nonroot /app/server /app/server

# Document the listening port; Kubernetes/compose maps the real port.
EXPOSE 8080
ENV PORT=8080

# Run the server as PID 1 for proper signal handling.
ENTRYPOINT ["/app/server"]
