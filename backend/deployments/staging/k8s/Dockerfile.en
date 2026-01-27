# syntax=docker/dockerfile:1.4
# PRODUCTION-READY DOCKERFILE FOR GO BACKEND
# Architecture: Multi-stage build with security hardening
# Base image: Distroless for minimal attack surface
# Build time: ~2-3 min first build, ~5-10 sec cached rebuilds

# ============================================================================
# STAGE 1: BUILDER - Compile Go binary
# ============================================================================
FROM golang:1.25.4-alpine AS builder

# Install build dependencies
# ca-certificates: Required for HTTPS requests in build stage
# git: Required for private Go modules (if any)
# tzdata: For timezone data (will copy to runtime stage)
RUN apk add --no-cache ca-certificates git tzdata

# Set working directory
WORKDIR /build

# Enable Go module proxy for faster, reliable dependency downloads
# GOPROXY: Public proxy with direct fallback for private repos
# GOSUMDB: Checksum database for security verification
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org

# ============================================================================
# LAYER CACHING OPTIMIZATION: Copy go.mod/go.sum FIRST
# WHY: Dependencies change rarely (maybe once a week)
#      Source code changes frequently (every commit)
#      This separation means we only re-download deps when go.mod changes
# RESULT: 95% of builds use cached layers, saving 2-3 minutes per build
# ============================================================================
COPY go.mod go.sum ./

# Download dependencies with cache mount for massive speed boost
# --mount=type=cache: Persists /go/pkg/mod across builds
# RESULT: First download = 30s, subsequent = 1s (from cache)
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download && \
    go mod verify

# Copy source code (this layer invalidates on every code change)
COPY . .

# ============================================================================
# BUILD CONFIGURATION: Static binary with optimizations
# ============================================================================
# CGO_ENABLED=0: Build static binary (no dynamic linking)
#   WHY: Can run on scratch/distroless (no libc required)
#   CRITICAL for minimal images
#
# GOOS=linux GOARCH=amd64: Explicit platform (important for M1 Macs)
#   WHY: Prevents "exec format error" when building on Mac, deploying on Linux
#
# -ldflags="-w -s": Strip debug info and symbol table
#   WHY: Reduces binary size by 30-40% (e.g., 20MB -> 13MB)
#   -w: Omit DWARF debug info
#   -s: Omit symbol table
#   TRADE-OFF: Can't use debuggers (gdb/dlv) on this binary
#              For staging/prod this is acceptable for size reduction
#
# -trimpath: Remove file system paths from binary
#   WHY: Security (don't leak /home/username/project paths)
#        Reproducible builds (same code = same binary hash)
#
# -o /app/server: Output binary location
# ./cmd/server: Entry point (main package)
# ============================================================================
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    go build \
        -ldflags="-w -s" \
        -trimpath \
        -o /app/server \
        ./cmd/server

# Verify the binary is statically linked (no dependencies)
# This command ensures we can run on minimal base images
RUN ldd /app/server 2>&1 | grep -q "not a dynamic executable" || \
    (echo "ERROR: Binary is not statically linked!" && exit 1)

# ============================================================================
# STAGE 2: RUNTIME - Minimal production image
# ============================================================================
# Using distroless instead of scratch because:
# - Includes CA certificates (for HTTPS API calls)
# - Includes tzdata (for timezone conversions)
# - Includes /etc/passwd (for non-root user)
# - Still minimal: ~2MB base image, zero package manager, zero shell
#
# Alternative options:
# - scratch: Absolute minimum (0 bytes), use if no HTTPS/timezone needs
# - alpine:3.20: ~5MB, includes shell (useful for debugging), package manager
# - distroless/base: ~20MB, includes glibc (use if CGO_ENABLED=1)
FROM gcr.io/distroless/static-debian12:nonroot

# ============================================================================
# SECURITY HARDENING: Non-root user
# ============================================================================
# distroless:nonroot includes user "nonroot" with UID 65532
# WHY non-root is CRITICAL:
# 1. Container escape vulnerability: If attacker breaks out of container,
#    they get same privileges on host. Root in container = root on host = total compromise
# 2. Kubernetes PodSecurityStandards: Starting 2024+, clusters enforce
#    runAsNonRoot=true. Root containers will be REJECTED.
# 3. Compliance: SOC2, PCI-DSS, HIPAA all require least privilege principle
# 4. Defense in depth: Even if app is compromised, attacker has limited privileges
#
# This USER directive is inherited by Kubernetes (unless overridden in pod spec)
USER nonroot:nonroot

# Set working directory
WORKDIR /app

# ============================================================================
# COPY ARTIFACTS FROM BUILDER STAGE
# ============================================================================
# Copy compiled binary from builder stage
# --from=builder: Reference previous stage
# --chown=nonroot:nonroot: Set ownership to non-root user
#   WHY: Prevents permission errors when app tries to read/write files
COPY --from=builder --chown=nonroot:nonroot /app/server /app/server

# Copy timezone data (if your app uses time.LoadLocation)
# Distroless includes this by default, but explicit copy if using scratch:
# COPY --from=builder --chown=nonroot:nonroot /usr/share/zoneinfo /usr/share/zoneinfo

# Copy CA certificates (if your app makes HTTPS requests)
# Distroless includes this by default, but explicit copy if using scratch:
# COPY --from=builder --chown=nonroot:nonroot /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# ============================================================================
# OPTIONAL: Copy migrations or static files if needed
# ============================================================================
# COPY --from=builder --chown=nonroot:nonroot /build/migrations /app/migrations

# ============================================================================
# RUNTIME CONFIGURATION
# ============================================================================
# Expose port (DOCUMENTATION ONLY - does not actually publish ports)
# The actual port mapping is configured in Kubernetes Service
# This documents that the app listens on 8080
EXPOSE 8080

# Environment variables (can be overridden in Kubernetes deployment)
# These are DEFAULTS only - sensitive values come from ConfigMap/Secret
ENV GIN_MODE=release \
    PORT=8080

# ============================================================================
# HEALTH CHECK (Kubernetes uses liveness/readiness probes, but this documents the endpoint)
# ============================================================================
# Kubernetes doesn't use Dockerfile HEALTHCHECK, but it's useful documentation
# Configure in K8s deployment.yaml instead:
#   livenessProbe:
#     httpGet:
#       path: /health
#       port: 8080
#     initialDelaySeconds: 10
#     periodSeconds: 30
# HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
#     CMD ["/app/server", "--healthcheck"]

# ============================================================================
# ENTRYPOINT: Run the application
# ============================================================================
# Using ENTRYPOINT (not CMD) for clarity
# ENTRYPOINT: Always executes, can't be overridden (without --entrypoint)
# CMD: Can be overridden with docker run arguments
#
# Exec form (JSON array) is preferred over shell form:
# GOOD: ["./server"]  -> Runs as PID 1, receives SIGTERM directly
# BAD:  "./server"    -> Runs through /bin/sh, SIGTERM goes to shell first
#
# WHY PID 1 matters:
# - Kubernetes sends SIGTERM to PID 1 for graceful shutdown
# - If shell is PID 1, it might not forward signal to app
# - Result: App doesn't shut down gracefully, connections get dropped
#
# GRACEFUL SHUTDOWN IN GO CODE:
# Your app should handle SIGTERM:
#   quit := make(chan os.Signal, 1)
#   signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
#   <-quit
#   ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
#   defer cancel()
#   if err := srv.Shutdown(ctx); err != nil {
#       log.Fatal("Server forced to shutdown:", err)
#   }
ENTRYPOINT ["/app/server"]

# ============================================================================
# BUILD & DEPLOYMENT INSTRUCTIONS
# ============================================================================
#
# LOCAL BUILD:
#   docker build -t orchestra-backend:local -f deployments/staging/k8s/Dockerfile .
#
# BUILD WITH GIT SHA TAG (recommended for production):
#   GIT_SHA=$(git rev-parse --short HEAD)
#   docker build -t orchestra-backend:${GIT_SHA} -f deployments/staging/k8s/Dockerfile .
#
# RUN LOCALLY (for testing):
#   docker run --rm -p 8080:8080 \
#     -e DATABASE_URL=postgres://user:pass@host:5432/db \
#     orchestra-backend:local
#
# SECURITY SCAN (run before deploying):
#   trivy image orchestra-backend:${GIT_SHA}
#   # Should show ZERO HIGH or CRITICAL vulnerabilities
#
# PUSH TO REGISTRY:
#   docker tag orchestra-backend:${GIT_SHA} gcr.io/your-project/orchestra-backend:${GIT_SHA}
#   docker push gcr.io/your-project/orchestra-backend:${GIT_SHA}
#
# DEPLOY TO KUBERNETES:
#   kubectl set image deployment/orchestra-backend \
#     backend=gcr.io/your-project/orchestra-backend:${GIT_SHA} \
#     -n staging
#
# ============================================================================
# PRODUCTION CHECKLIST
# ============================================================================
# [ ] Multi-stage build (build + runtime)
# [ ] Static binary (CGO_ENABLED=0, verified with ldd)
# [ ] Non-root user (USER nonroot)
# [ ] Minimal base image (distroless/static)
# [ ] Layer caching optimized (go.mod before source code)
# [ ] .dockerignore configured
# [ ] No secrets hardcoded (verify with: docker history image:tag)
# [ ] Security scan passed (trivy/snyk/clair)
# [ ] Image tagged with git SHA (not :latest)
# [ ] Graceful shutdown implemented (SIGTERM handler in code)
# [ ] Health checks defined in K8s deployment
# [ ] Resource limits set in K8s (requests/limits for CPU/memory)
# [ ] Tested in staging before prod
#
# ============================================================================
# TROUBLESHOOTING
# ============================================================================
#
# ERROR: "exec format error"
#   CAUSE: Binary built for wrong architecture (e.g., ARM on M1 Mac, deployed to AMD64 Linux)
#   FIX: Ensure GOOS=linux GOARCH=amd64 in build stage
#
# ERROR: "permission denied" when container starts
#   CAUSE: Binary not executable, or non-root user can't access files
#   FIX: Ensure COPY --chown=nonroot:nonroot, or RUN chmod +x /app/server
#
# ERROR: Container crashes immediately
#   DEBUG: Run with shell access (for debugging only):
#     docker run --rm -it --entrypoint /bin/sh alpine:3.20
#     # Then inspect the binary, check ldd, etc.
#
# ERROR: "x509: certificate signed by unknown authority"
#   CAUSE: Missing CA certificates (common with scratch base)
#   FIX: Use distroless/static (includes CA certs) or copy from builder:
#     COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
#
# ERROR: App can't connect to database
#   CAUSE: DATABASE_URL not set, or network policy blocking connection
#   DEBUG: Check environment variables in K8s:
#     kubectl exec -it pod-name -- env | grep DATABASE
#
# ============================================================================
# PERFORMANCE OPTIMIZATION
# ============================================================================
#
# CURRENT BUILD TIME:
# - First build (no cache): ~2-3 minutes
# - With layer cache: ~5-10 seconds
# - With go mod cache: ~1-2 seconds
#
# FURTHER OPTIMIZATIONS:
# 1. Use private GOPROXY (Athens/Artifactory) for faster, reliable downloads
# 2. Use Docker BuildKit cache backend (registry or S3):
#      docker buildx build --cache-from=type=registry,ref=image:cache ...
# 3. Parallel builds if multiple binaries:
#      RUN go build ./cmd/server & go build ./cmd/worker & wait
# 4. Use multi-platform builds for ARM + AMD64:
#      docker buildx build --platform linux/amd64,linux/arm64 ...
#
# IMAGE SIZE:
# - Builder stage: ~800MB (golang:1.25-alpine + deps)
# - Runtime stage: ~15MB (2MB distroless + 13MB binary)
# - Optimization: Binary size can be reduced further with UPX compression
#   (but adds runtime decompression overhead, rarely worth it)
#
# ============================================================================
# COST ANALYSIS (example for GKE with 10 replicas)
# ============================================================================
#
# Image size impact:
# - Large image (500MB): 10 nodes × 500MB × $0.10/GB/month = $5/month storage
#                        10 pulls/day × 500MB × $0.12/GB egress = $18/month
# - Small image (15MB): 10 nodes × 15MB × $0.10/GB/month = $0.15/month storage
#                       10 pulls/day × 15MB × $0.12/GB egress = $0.54/month
#
# SAVINGS: $22.31/month per service
# At 50 microservices: $1,115/month = $13,380/year saved
#
# This is why FAANG companies obsess over image size.
#
# ============================================================================
