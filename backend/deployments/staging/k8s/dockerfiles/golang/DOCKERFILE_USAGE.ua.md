# ПОСІБНИК З ВИКОРИСТАННЯ PRODUCTION DOCKERFILE (GOLANG)

## Швидкий старт

### Збірка образу
```bash
cd /home/r/fff/deploy/orchestra/backend

GIT_SHA=$(git rev-parse --short HEAD)
docker build -t golang-app:${GIT_SHA} -f deployments/staging/k8s/dockerfiles/golang/Dockerfile .
```

### Локальне тестування
```bash
docker run --rm -p 8080:8080 \
  -e PORT=8080 \
  golang-app:${GIT_SHA}

curl http://localhost:8080/health
```

### Сканування безпеки (перед деплоєм)
```bash
trivy image golang-app:${GIT_SHA}
```

### Публікація в реєстр (приклад)
```bash
docker tag golang-app:${GIT_SHA} gcr.io/your-project/golang-app:${GIT_SHA}
docker push gcr.io/your-project/golang-app:${GIT_SHA}
```

### Деплой в Kubernetes (приклад)
```bash
kubectl set image deployment/your-service \
  app=gcr.io/your-project/golang-app:${GIT_SHA} \
  -n staging
kubectl rollout status deployment/your-service -n staging
```
