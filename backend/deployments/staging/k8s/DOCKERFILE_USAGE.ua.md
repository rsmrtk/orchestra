# ПОСІБНИК З ВИКОРИСТАННЯ PRODUCTION DOCKERFILE

## Швидкий старт

### Збірка образу
```bash
# Перехід до кореневої директорії проєкту
cd /home/r/fff/deploy/orchestra/backend

# Збірка з git SHA тегом (рекомендовано)
GIT_SHA=$(git rev-parse --short HEAD)
docker build -t orchestra-backend:${GIT_SHA} -f deployments/staging/k8s/Dockerfile .

# Або збірка з семантичною версією
docker build -t orchestra-backend:v1.0.0 -f deployments/staging/k8s/Dockerfile .
```

### Локальне тестування
```bash
# Запуск з змінними середовища
docker run --rm -p 8080:8080 \
  -e DATABASE_URL="postgres://user:pass@host:5432/dbname?sslmode=require" \
  -e GIN_MODE=release \
  orchestra-backend:${GIT_SHA}

# Тест health endpoint
curl http://localhost:8080/health
```

### Сканування безпеки (КРИТИЧНО перед розгортанням)
```bash
# Встановлення trivy
# Ubuntu/Debian: apt install trivy
# Mac: brew install trivy

# Сканування на вразливості
trivy image orchestra-backend:${GIT_SHA}

# КРИТЕРІЙ ПРИЙНЯТТЯ: Нуль HIGH або CRITICAL CVE
# Якщо знайдено, оновіть базовий образ або залежності
```

### Push в Container Registry

#### Google Container Registry (GCR)
```bash
# Автентифікація
gcloud auth configure-docker

# Тегування та push
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

### Розгортання в Kubernetes
```bash
# Оновлення deployment з новим образом
kubectl set image deployment/orchestra-backend \
  backend=gcr.io/your-project-id/orchestra-backend:${GIT_SHA} \
  -n staging

# Моніторинг статусу rollout
kubectl rollout status deployment/orchestra-backend -n staging

# Перевірка розгортання
kubectl get pods -n staging
kubectl logs -f deployment/orchestra-backend -n staging
```

---

## Інтеграція CI/CD

### Приклад GitHub Actions
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

## Інтеграція з Kubernetes

### Оновлення deployment.yaml для використання образу
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
      # Security context (примушує non-root на рівні pod)
      securityContext:
        runAsNonRoot: true
        runAsUser: 65532  # Відповідає nonroot користувачу в distroless
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

        # Resource limits (КРИТИЧНО для production)
        resources:
          requests:
            cpu: 100m       # 0.1 ядра CPU
            memory: 128Mi   # 128 MB
          limits:
            cpu: 500m       # Максимум 0.5 ядра CPU
            memory: 512Mi   # Максимум 512 MB

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

        # Плавне завершення
        lifecycle:
          preStop:
            exec:
              command: ["/bin/sleep", "15"]  # Очікування 15с перед SIGTERM

        # Змінні середовища з ConfigMap/Secret
        envFrom:
        - configMapRef:
            name: orchestra-backend-config
        - secretRef:
            name: orchestra-backend-secret

        ports:
        - containerPort: 8080
          protocol: TCP
```

### Застосування deployment
```bash
# Розгортання
kubectl apply -f deployments/staging/k8s/deployment.yaml

# Оновлення образу
kubectl set image deployment/orchestra-backend \
  backend=gcr.io/your-project/orchestra-backend:abc123 \
  -n staging

# Rollback якщо щось пішло не так
kubectl rollout undo deployment/orchestra-backend -n staging
```

---

## Вирішення проблем

### Проблеми збірки

#### Помилка: "exec format error"
```bash
# ПРИЧИНА: Зібрано для неправильної архітектури (ARM на M1 Mac, розгортання на AMD64 Linux)
# ВИПРАВЛЕННЯ: Явно встановіть GOARCH в Dockerfile (вже зроблено в наданому Dockerfile)
docker buildx build --platform linux/amd64 -f deployments/staging/k8s/Dockerfile .
```

#### Помилка: "go: modules not found"
```bash
# ПРИЧИНА: go.mod/go.sum відсутні в build context
# ВИПРАВЛЕННЯ: Переконайтеся, що збираєте з кореневої директорії проєкту
cd /home/r/fff/deploy/orchestra/backend
docker build -f deployments/staging/k8s/Dockerfile .
```

#### Повільні збірки
```bash
# ПРИЧИНА: Не використовується кеш шарів
# ВИПРАВЛЕННЯ 1: Використовуйте BuildKit (вже увімкнено з # syntax=docker/dockerfile:1.4)
export DOCKER_BUILDKIT=1

# ВИПРАВЛЕННЯ 2: Використовуйте registry кеш
docker buildx build --cache-from type=registry,ref=your-image:cache ...
```

### Проблеми Runtime

#### Помилка: "permission denied"
```bash
# ПРИЧИНА: Non-root користувач не може писати у файлову систему
# ВИПРАВЛЕННЯ: Використовуйте emptyDir volumes в K8s для директорій з можливістю запису
# В deployment.yaml:
volumes:
- name: tmp
  emptyDir: {}
volumeMounts:
- name: tmp
  mountPath: /tmp
```

#### Помилка: "x509: certificate signed by unknown authority"
```bash
# ПРИЧИНА: Відсутні CA сертифікати
# ВИПРАВЛЕННЯ: Вже оброблено в Dockerfile (distroless включає CA сертифікати)
# Якщо використовується scratch, розкоментуйте цей рядок в Dockerfile:
# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
```

#### Контейнер падає негайно
```bash
# НАЛАГОДЖЕННЯ: Перевірте логи
kubectl logs deployment/orchestra-backend -n staging --previous

# НАЛАГОДЖЕННЯ: Опишіть pod
kubectl describe pod -l app=orchestra-backend -n staging
```

#### Додаток не може підключитися до бази даних
```bash
# НАЛАГОДЖЕННЯ: Перевірте змінні середовища
kubectl exec -it deployment/orchestra-backend -n staging -- env | grep DATABASE

# ВИПРАВЛЕННЯ: Переконайтеся, що DATABASE_URL встановлено в Secret/ConfigMap
# ВИПРАВЛЕННЯ: Переконайтеся, що NetworkPolicy дозволяє з'єднання з базою даних
```

---

## Чекліст безпеки

Перед розгортанням в production:

- [ ] Пройдено сканування безпеки (trivy/snyk): `trivy image orchestra-backend:${GIT_SHA}`
- [ ] Немає секретів в образі: `docker history orchestra-backend:${GIT_SHA} | grep -i secret`
- [ ] Налаштовано non-root користувача: `docker run --rm orchestra-backend:${GIT_SHA} id`
- [ ] Працює read-only файлова система: Тест з `securityContext.readOnlyRootFilesystem: true`
- [ ] Образ підписано (cosign): `cosign sign gcr.io/project/orchestra-backend:${GIT_SHA}`
- [ ] Згенеровано SBOM: `syft packages orchestra-backend:${GIT_SHA} -o spdx-json > sbom.json`

---

## Наступні кроки

### 1. Реалізація плавного завершення в Go коді
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

    // Запуск сервера в goroutine
    go func() {
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("Не вдалося запустити сервер: %v", err)
        }
    }()

    // Очікування interrupt сигналу (SIGTERM від Kubernetes)
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
    <-quit

    log.Println("Завершення роботи сервера...")

    // Плавне завершення з 30-секундним timeout
    // Kubernetes очікує terminationGracePeriodSeconds (за замовчуванням 30с)
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("Сервер примусово завершено: %v", err)
        return err
    }

    log.Println("Сервер завершено плавно")
    return nil
}
```

### 2. Додавання Health Check endpoint
```go
// internal/rest/controllers/health.go
package controllers

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func (c *Controllers) HealthCheck(ctx *gin.Context) {
    // TODO: Перевірка підключення до бази даних, залежностей
    // db.Ping(), redis.Ping(), і т.д.

    ctx.JSON(http.StatusOK, gin.H{
        "status": "healthy",
        "version": os.Getenv("VERSION"),
    })
}

// Реєстрація route
router.GET("/health", controllers.HealthCheck)
router.GET("/ready", controllers.ReadinessCheck)
```

### 3. Додавання Prometheus метрик
```go
import "github.com/prometheus/client_golang/prometheus/promhttp"

router.GET("/metrics", gin.WrapH(promhttp.Handler()))
```

### 4. Налаштування CI/CD Pipeline
- Скопіюйте приклад GitHub Actions вище
- Налаштуйте секрети (GCP_SA_KEY, DATABASE_URL, і т.д.)
- Увімкніть захист гілок (вимагайте проходження перевірок)
- Налаштуйте просування staging → production

---

## Ресурси

### Документація
- [Docker Multi-stage Builds](https://docs.docker.com/build/building/multi-stage/)
- [Docker BuildKit](https://docs.docker.com/build/buildkit/)
- [Distroless Images](https://github.com/GoogleContainerTools/distroless)
- [Kubernetes Best Practices](https://kubernetes.io/docs/concepts/configuration/overview/)

### Інструменти безпеки
- [Trivy](https://github.com/aquasecurity/trivy) - Сканер вразливостей
- [Cosign](https://github.com/sigstore/cosign) - Підписування образів
- [Syft](https://github.com/anchore/syft) - Генерація SBOM

### Моніторинг
- [Prometheus](https://prometheus.io/) - Метрики
- [Grafana](https://grafana.com/) - Дашборди
- [Loki](https://grafana.com/oss/loki/) - Агрегація логів
- [Jaeger](https://www.jaegertracing.io/) - Розподілене трейсування

---

## Часті запитання

**П: Чому distroless замість alpine?**
В: Distroless має нуль менеджера пакетів, нуль shell, меншу поверхню атаки. Alpine має shell (корисно для налагодження, але ризик безпеки). Для Go статичних бінарників distroless переважний.

**П: Чи можу я використовувати scratch замість distroless?**
В: Так, якщо ваш додаток не робить HTTPS викликів (потрібні CA сертифікати) і не використовує часові пояси. Distroless це scratch + мінімальні runtime залежності.

**П: Чому multi-stage збірка?**
В: Розділяє залежності збірки (gcc, git, Go SDK) від runtime. Образ збірки = 1GB+, runtime = 20MB. Швидші розгортання, нижчі витрати, менша поверхня атаки.

**П: Чи слід використовувати :latest тег?**
В: НІКОЛИ в production. Використовуйте git SHA або семантичну версію. :latest непередбачуваний і робить rollback неможливими.

**П: Як часто слід перебудовувати базові образи?**
В: Принаймні щомісяця для безпекових патчів. Налаштуйте автоматизовані збірки з Dependabot або Renovate для оновлення golang:1.25.4-alpine при випуску нових версій.
