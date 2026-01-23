# ПОСІБНИК З ВИКОРИСТАННЯ PRODUCTION DOCKERFILE (JAVASCRIPT/NODE.JS)

## Швидкий старт

### Збірка образу
```bash
cd /home/r/fff/deploy/orchestra/backend

GIT_SHA=$(git rev-parse --short HEAD)
docker build -t node-app:${GIT_SHA} -f deployments/staging/k8s/dockerfiles/javascript/Dockerfile .
```

### Локальне тестування
```bash
docker run --rm -p 3000:3000 \
  -e PORT=3000 \
  node-app:${GIT_SHA}

curl http://localhost:3000/health
```

### Сканування безпеки (перед деплоєм)
```bash
trivy image node-app:${GIT_SHA}
```

### Публікація в реєстр (приклад)
```bash
docker tag node-app:${GIT_SHA} gcr.io/your-project/node-app:${GIT_SHA}
docker push gcr.io/your-project/node-app:${GIT_SHA}
```

### Деплой в Kubernetes (приклад)
```bash
kubectl set image deployment/your-service \
  app=gcr.io/your-project/node-app:${GIT_SHA} \
  -n staging
kubectl rollout status deployment/your-service -n staging
```
