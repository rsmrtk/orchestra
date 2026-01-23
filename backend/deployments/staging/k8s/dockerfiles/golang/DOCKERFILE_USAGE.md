# DOCKERFILE PRODUCTION USAGE GUIDE (GOLANG)

## Quick Start

### Build the Image
```bash
cd /home/r/fff/deploy/orchestra/backend

GIT_SHA=$(git rev-parse --short HEAD)
docker build -t golang-app:${GIT_SHA} -f deployments/staging/k8s/dockerfiles/golang/Dockerfile .
```

### Test Locally
```bash
docker run --rm -p 8080:8080 \
  -e PORT=8080 \
  golang-app:${GIT_SHA}

curl http://localhost:8080/health
```

### Security Scan (Before Deploy)
```bash
trivy image golang-app:${GIT_SHA}
```

### Push to Registry (Example)
```bash
docker tag golang-app:${GIT_SHA} gcr.io/your-project/golang-app:${GIT_SHA}
docker push gcr.io/your-project/golang-app:${GIT_SHA}
```

### Deploy to Kubernetes (Example)
```bash
kubectl set image deployment/your-service \
  app=gcr.io/your-project/golang-app:${GIT_SHA} \
  -n staging
kubectl rollout status deployment/your-service -n staging
```
