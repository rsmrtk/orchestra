# DOCKERFILE PRODUCTION USAGE GUIDE (PYTHON)

## Quick Start

### Build the Image
```bash
cd /home/r/fff/deploy/orchestra/backend

GIT_SHA=$(git rev-parse --short HEAD)
docker build -t python-app:${GIT_SHA} -f deployments/staging/k8s/dockerfiles/python/Dockerfile .
```

### Test Locally
```bash
docker run --rm -p 8000:8000 \
  -e PORT=8000 \
  python-app:${GIT_SHA}

curl http://localhost:8000/health
```

### Security Scan (Before Deploy)
```bash
trivy image python-app:${GIT_SHA}
```

### Push to Registry (Example)
```bash
docker tag python-app:${GIT_SHA} gcr.io/your-project/python-app:${GIT_SHA}
docker push gcr.io/your-project/python-app:${GIT_SHA}
```

### Deploy to Kubernetes (Example)
```bash
kubectl set image deployment/your-service \
  app=gcr.io/your-project/python-app:${GIT_SHA} \
  -n staging
kubectl rollout status deployment/your-service -n staging
```
