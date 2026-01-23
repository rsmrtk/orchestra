# ПОСІБНИК З ВИКОРИСТАННЯ PRODUCTION DOCKERFILE (PYTHON)

## Швидкий старт

### Збірка образу
```bash
cd /home/r/fff/deploy/orchestra/backend

GIT_SHA=$(git rev-parse --short HEAD)
docker build -t python-app:${GIT_SHA} -f deployments/staging/k8s/dockerfiles/python/Dockerfile .
```

### Локальне тестування
```bash
docker run --rm -p 8000:8000 \
  -e PORT=8000 \
  python-app:${GIT_SHA}

curl http://localhost:8000/health
```

### Сканування безпеки (перед деплоєм)
```bash
trivy image python-app:${GIT_SHA}
```

### Публікація в реєстр (приклад)
```bash
docker tag python-app:${GIT_SHA} gcr.io/your-project/python-app:${GIT_SHA}
docker push gcr.io/your-project/python-app:${GIT_SHA}
```

### Деплой в Kubernetes (приклад)
```bash
kubectl set image deployment/your-service \
  app=gcr.io/your-project/python-app:${GIT_SHA} \
  -n staging
kubectl rollout status deployment/your-service -n staging
```
