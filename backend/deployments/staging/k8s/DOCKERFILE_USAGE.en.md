# DOCKERFILE PRODUCTION USAGE GUIDE

## Quick Start

### Build the Image
```bash
# Navigate to project root
cd /home/r/fff/deploy/orchestra/backend

# Build with git SHA tag (recommended)
GIT_SHA=$(git rev-parse --short HEAD)
docker build -t orchestra-backend:${GIT_SHA} -f deployments/staging/k8s/Dockerfile .

# Or build with semantic version
docker build -t orchestra-backend:v1.0.0 -f deployments/staging/k8s/Dockerfile .
```

### Test Locally
```bash
# Run with environment variables
docker run --rm -p 8080:8080 \
  -e DATABASE_URL="postgres://user:pass@host:5432/dbname?sslmode=require" \
  -e GIN_MODE=release \
  orchestra-backend:${GIT_SHA}

# Test health endpoint
curl http://localhost:8080/health
```

### Security Scan (CRITICAL before deploying)
```bash
# Install trivy
# Ubuntu/Debian: apt install trivy
# Mac: brew install trivy

# Scan for vulnerabilities
trivy image orchestra-backend:${GIT_SHA}

# ACCEPTANCE CRITERIA: Zero HIGH or CRITICAL CVEs
# If found, update base image or dependencies
```

### Push to Container Registry

#### Google Container Registry (GCR)
```bash
# Authenticate
gcloud auth configure-docker

# Tag and push
docker tag orchestra-backend:${GIT_SHA} gcr.io/your-project-id/orchestra-backend:${GIT_SHA}
docker push gcr.io/your-project-id/orchestra-backend:${GIT_SHA}
```

#### Docker Hub
```bash
docker tag orchestra-backend:${GIT_SHA} yourusername/orchestra-backend:${GIT_SHA}
docker push yourusername/orchestra-backend:${GIT_SHA}
```

#### AWS ECR
```bash
aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin 123456789.dkr.ecr.us-east-1.amazonaws.com
docker tag orchestra-backend:${GIT_SHA} 123456789.dkr.ecr.us-east-1.amazonaws.com/orchestra-backend:${GIT_SHA}
docker push 123456789.dkr.ecr.us-east-1.amazonaws.com/orchestra-backend:${GIT_SHA}
```

### Deploy to Kubernetes
```bash
# Update deployment with new image
kubectl set image deployment/orchestra-backend \
  backend=gcr.io/your-project-id/orchestra-backend:${GIT_SHA} \
  -n staging

# Watch rollout status
kubectl rollout status deployment/orchestra-backend -n staging

# Verify deployment
kubectl get pods -n staging
kubectl logs -f deployment/orchestra-backend -n staging
```

---

## CI/CD Integration

### GitHub Actions Example
```yaml
name: Build and Deploy

on:
  push:
    branches: [main, staging]

env:
  PROJECT_ID: your-gcp-project
  SERVICE_NAME: orchestra-backend
  REGION: us-central1

jobs:
  build-and-deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Authenticate to GCP
        uses: google-github-actions/auth@v2
        with:
          credentials_json: ${{ secrets.GCP_SA_KEY }}

      - name: Configure Docker for GCR
        run: gcloud auth configure-docker

      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          file: ./deployments/staging/k8s/Dockerfile
          push: true
          tags: |
            gcr.io/${{ env.PROJECT_ID }}/${{ env.SERVICE_NAME }}:${{ github.sha }}
            gcr.io/${{ env.PROJECT_ID }}/${{ env.SERVICE_NAME }}:latest
          cache-from: type=registry,ref=gcr.io/${{ env.PROJECT_ID }}/${{ env.SERVICE_NAME }}:cache
          cache-to: type=registry,ref=gcr.io/${{ env.PROJECT_ID }}/${{ env.SERVICE_NAME }}:cache,mode=max

      - name: Security scan
        run: |
          docker pull gcr.io/${{ env.PROJECT_ID }}/${{ env.SERVICE_NAME }}:${{ github.sha }}
          trivy image --exit-code 1 --severity HIGH,CRITICAL \
            gcr.io/${{ env.PROJECT_ID }}/${{ env.SERVICE_NAME }}:${{ github.sha }}

      - name: Deploy to GKE
        run: |
          gcloud container clusters get-credentials staging-cluster --region ${{ env.REGION }}
          kubectl set image deployment/${{ env.SERVICE_NAME }} \
            backend=gcr.io/${{ env.PROJECT_ID }}/${{ env.SERVICE_NAME }}:${{ github.sha }} \
            -n staging
          kubectl rollout status deployment/${{ env.SERVICE_NAME }} -n staging --timeout=5m
```

---

## Kubernetes Integration

### Update deployment.yaml to use the image
```yaml
# deployments/staging/k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: orchestra-backend
  namespace: staging
spec:
  replicas: 3
  selector:
    matchLabels:
      app: orchestra-backend
  template:
    metadata:
      labels:
        app: orchestra-backend
    spec:
      # Security context (enforces non-root at pod level)
      securityContext:
        runAsNonRoot: true
        runAsUser: 65532  # Matches nonroot user in distroless
        fsGroup: 65532
        seccompProfile:
          type: RuntimeDefault

      containers:
      - name: backend
        image: gcr.io/your-project/orchestra-backend:REPLACE_WITH_GIT_SHA
        imagePullPolicy: Always

        # Container security context
        securityContext:
          allowPrivilegeEscalation: false
          readOnlyRootFilesystem: true
          runAsNonRoot: true
          runAsUser: 65532
          capabilities:
            drop:
              - ALL

        # Resource limits (CRITICAL for production)
        resources:
          requests:
            cpu: 100m       # 0.1 CPU core
            memory: 128Mi   # 128 MB
          limits:
            cpu: 500m       # 0.5 CPU core max
            memory: 512Mi   # 512 MB max

        # Health checks
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 30
          timeoutSeconds: 5
          failureThreshold: 3

        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 10
          timeoutSeconds: 3
          failureThreshold: 3

        # Graceful shutdown
        lifecycle:
          preStop:
            exec:
              command: ["/bin/sleep", "15"]  # Wait 15s before SIGTERM

        # Environment variables from ConfigMap/Secret
        envFrom:
        - configMapRef:
            name: orchestra-backend-config
        - secretRef:
            name: orchestra-backend-secret

        ports:
        - containerPort: 8080
          protocol: TCP
```

### Apply deployment
```bash
# Deploy
kubectl apply -f deployments/staging/k8s/deployment.yaml

# Update image
kubectl set image deployment/orchestra-backend \
  backend=gcr.io/your-project/orchestra-backend:abc123 \
  -n staging

# Rollback if something goes wrong
kubectl rollout undo deployment/orchestra-backend -n staging
```

---

## Troubleshooting

### Build Issues

#### Error: "exec format error"
```bash
# CAUSE: Built for wrong architecture (ARM on M1 Mac, deploying to AMD64 Linux)
# FIX: Explicitly set GOARCH in Dockerfile (already done in provided Dockerfile)
docker buildx build --platform linux/amd64 -f deployments/staging/k8s/Dockerfile .
```

#### Error: "go: modules not found"
```bash
# CAUSE: go.mod/go.sum not in build context
# FIX: Ensure you're building from project root
cd /home/r/fff/deploy/orchestra/backend
docker build -f deployments/staging/k8s/Dockerfile .
```

#### Slow builds
```bash
# CAUSE: Not using layer cache
# FIX 1: Use BuildKit (already enabled with # syntax=docker/dockerfile:1.4)
export DOCKER_BUILDKIT=1

# FIX 2: Use registry cache
docker buildx build --cache-from type=registry,ref=your-image:cache ...
```

### Runtime Issues

#### Error: "permission denied"
```bash
# CAUSE: Non-root user can't write to filesystem
# FIX: Use emptyDir volumes in K8s for writable directories
# In deployment.yaml:
volumes:
- name: tmp
  emptyDir: {}
volumeMounts:
- name: tmp
  mountPath: /tmp
```

#### Error: "x509: certificate signed by unknown authority"
```bash
# CAUSE: Missing CA certificates
# FIX: Already handled in Dockerfile (distroless includes CA certs)
# If using scratch, uncomment this line in Dockerfile:
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
```

#### Container crashes immediately
```bash
# DEBUG: Check logs
kubectl logs deployment/orchestra-backend -n staging --previous

# DEBUG: Describe pod
kubectl describe pod -l app=orchestra-backend -n staging
```

#### App can't connect to database
```bash
# DEBUG: Check environment variables
kubectl exec -it deployment/orchestra-backend -n staging -- env | grep DATABASE

# FIX: Ensure DATABASE_URL is set in Secret/ConfigMap
# FIX: Ensure NetworkPolicy allows connection to database
```

---

## Security Checklist

Before deploying to production:

- [ ] Security scan passed (trivy/snyk): `trivy image orchestra-backend:${GIT_SHA}`
- [ ] No secrets in image: `docker history orchestra-backend:${GIT_SHA} | grep -i secret`
- [ ] Non-root user configured: `docker run --rm orchestra-backend:${GIT_SHA} id`
- [ ] Read-only filesystem works: Test with `securityContext.readOnlyRootFilesystem: true`
- [ ] Image signed (cosign): `cosign sign gcr.io/project/orchestra-backend:${GIT_SHA}`
- [ ] SBOM generated: `syft packages orchestra-backend:${GIT_SHA} -o spdx-json > sbom.json`

---

## Next Steps

### 1. Implement Graceful Shutdown in Go Code
```go
// internal/rest/server.go
package rest

import (
    "context"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"
)

func (s *Server) Run() error {
    srv := &http.Server{
        Addr:    ":8080",
        Handler: s.router,
    }

    // Start server in goroutine
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Failed to start server: %v", err)
        }
    }()

    // Wait for interrupt signal (SIGTERM from Kubernetes)
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
    <-quit

    log.Println("Shutting down server...")

    // Graceful shutdown with 30 second timeout
    // Kubernetes waits terminationGracePeriodSeconds (default 30s)
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Server forced to shutdown: %v", err)
        return err
    }

    log.Println("Server exited gracefully")
    return nil
}
```

### 2. Add Health Check Endpoint
```go
// internal/rest/controllers/health.go
package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func (c *Controllers) HealthCheck(ctx *gin.Context) {
    // TODO: Check database connection, dependencies
    // db.Ping(), redis.Ping(), etc.

    ctx.JSON(http.StatusOK, gin.H{
        "status": "healthy",
        "version": os.Getenv("VERSION"),
    })
}

// Register route
router.GET("/health", controllers.HealthCheck)
router.GET("/ready", controllers.ReadinessCheck)
```

### 3. Add Prometheus Metrics
```go
import "github.com/prometheus/client_golang/prometheus/promhttp"

router.GET("/metrics", gin.WrapH(promhttp.Handler()))
```

### 4. Set Up CI/CD Pipeline
- Copy GitHub Actions example above
- Configure secrets (GCP_SA_KEY, DATABASE_URL, etc.)
- Enable branch protection (require checks to pass)
- Set up staging → production promotion

---

## Resources

### Documentation
- [Docker Multi-stage Builds](https://docs.docker.com/build/building/multi-stage/)
- [Docker BuildKit](https://docs.docker.com/build/buildkit/)
- [Distroless Images](https://github.com/GoogleContainerTools/distroless)
- [Kubernetes Best Practices](https://kubernetes.io/docs/concepts/configuration/overview/)

### Security Tools
- [Trivy](https://github.com/aquasecurity/trivy) - Vulnerability scanner
- [Cosign](https://github.com/sigstore/cosign) - Image signing
- [Syft](https://github.com/anchore/syft) - SBOM generation

### Monitoring
- [Prometheus](https://prometheus.io/) - Metrics
- [Grafana](https://grafana.com/) - Dashboards
- [Loki](https://grafana.com/oss/loki/) - Log aggregation
- [Jaeger](https://www.jaegertracing.io/) - Distributed tracing

---

## Frequently Asked Questions

**Q: Why distroless instead of alpine?**
A: Distroless has zero package manager, zero shell, smaller attack surface. Alpine has shell (useful for debugging but security risk). For Go static binaries, distroless is preferred.

**Q: Can I use scratch instead of distroless?**
A: Yes, if your app doesn't make HTTPS calls (needs CA certs) and doesn't use timezones. Distroless is scratch + minimal runtime dependencies.

**Q: Why multi-stage build?**
A: Separates build dependencies (gcc, git, Go SDK) from runtime. Build image = 1GB+, runtime = 20MB. Faster deployments, lower costs, smaller attack surface.

**Q: Should I use :latest tag?**
A: NEVER in production. Use git SHA or semantic version. :latest is unpredictable and makes rollbacks impossible.

**Q: How often should I rebuild base images?**
A: At least monthly for security patches. Set up automated builds with Dependabot or Renovate to update golang:1.25.4-alpine when new versions release.
