# Orchestra Production Deployment - Повний Посібник

Це production-ready Kubernetes конфігурація для Orchestra Backend + Frontend.

## 📁 Структура Файлів

```
production/k8s/
├── namespace.yaml                 # Namespace для ізоляції
├── secret.yaml                    # Секрети (паролі, API keys)
├── configmap.yaml                 # Конфігурація backend
├── network-policy.yaml            # Мережева безпека
│
├── postgres/
│   ├── postgres-pvc.yaml          # Persistent storage для БД
│   ├── postgres-deployment.yaml   # PostgreSQL Deployment
│   └── postgres-service.yaml      # Service для БД
│
└── backend/
    ├── backend-deployment.yaml    # Backend API Deployment
    ├── backend-service.yaml       # Service для backend
    ├── backend-hpa.yaml           # Автоскейлінг
    ├── backend-pdb.yaml           # Захист від downtime
    └── backend-ingress.yaml       # Зовнішній доступ

frontend/k8s/
├── frontend-configmap.yaml        # Конфігурація frontend
├── frontend-hpa.yaml              # Автоскейлінг frontend
├── frontend-pdb.yaml              # Захист від downtime
└── frontend-network-policy.yaml   # Мережева безпека
```

## 🚀 Швидкий Старт

### Крок 1: Перевірити передумови

```bash
# Kubernetes кластер працює?
kubectl cluster-info

# Metrics Server встановлений? (потрібен для HPA)
kubectl get deployment metrics-server -n kube-system

# Якщо немає Metrics Server:
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# Для minikube:
minikube addons enable metrics-server
minikube addons enable ingress
```

### Крок 2: Налаштувати Secret (ВАЖЛИВО!)

```bash
# НЕ використовуй secret.yaml з git (небезпечно!)
# Створи Secret вручну з реальними паролями:

kubectl create secret generic orchestra-backend-secret \
  --from-literal=DB_HOST=orchestra-db \
  --from-literal=DB_PORT=2828 \
  --from-literal=DB_USER=your-username \
  --from-literal=DB_PASSWORD=your-strong-password \
  --from-literal=DB_NAME=orchestra-db \
  --namespace=orchestra-space --dry-run=client -o yaml | kubectl apply -f -
```

### Крок 3: Розгорнути Backend

```bash
# 1. Створити namespace
kubectl apply -f namespace.yaml

# 2. ConfigMap (відредагуй перед apply!)
kubectl apply -f configmap.yaml

# 3. PostgreSQL
kubectl apply -f postgres/

# 4. Дочекатися поки PostgreSQL Ready
kubectl wait --for=condition=ready pod -l app=postgres -n orchestra-space --timeout=120s

# 5. Backend
kubectl apply -f backend/

# 6. NetworkPolicy
kubectl apply -f network-policy.yaml
```

### Крок 4: Розгорнути Frontend

```bash
cd ../../frontend/deployments/production/k8s/

# 1. Frontend конфігурація (відредагуй API_URL!)
kubectl apply -f frontend-configmap.yaml

# 2. Frontend deployment (використай існуючі deployment.yaml, service.yaml, ingress.yaml)
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
kubectl apply -f ingress.yaml
kubectl apply -f backend-config.yaml

# 3. HPA, PDB, NetworkPolicy
kubectl apply -f frontend-hpa.yaml
kubectl apply -f frontend-pdb.yaml
kubectl apply -f frontend-network-policy.yaml
```

### Крок 5: Перевірити

```bash
# Подивитися всі ресурси
kubectl get all -n orchestra-space

# Pod'и мають бути Running та Ready
kubectl get pods -n orchestra-space

# Services
kubectl get svc -n orchestra-space

# HPA
kubectl get hpa -n orchestra-space

# PDB
kubectl get pdb -n orchestra-space

# NetworkPolicy
kubectl get networkpolicy -n orchestra-space
```

## 🔧 Налаштування перед Deploy

### 1. Secret.yaml (Критично!)

⚠️ **НЕ КОМІТИТИ З РЕАЛЬНИМИ ПАРОЛЯМИ!**

Створи сильні паролі:
```bash
# Generate strong password
openssl rand -base64 32

# Encode в base64
echo -n "your-password" | base64
```

### 2. ConfigMap.yaml

Відредагуй:
- `GIN_MODE`: "release" для production
- `ALLOWED_ORIGINS`: твій frontend domain
- `DB_MAX_OPEN_CONNS`: залежить від навантаження

### 3. Postgres Resources

У `postgres-deployment.yaml` змін:
```yaml
resources:
  requests:
    memory: "1Gi"     # Мінімум 1-2Gi для production
    cpu: "500m"
  limits:
    memory: "2Gi"
    cpu: "1000m"
```

### 4. Postgres PVC Size

У `postgres-pvc.yaml`:
```yaml
resources:
  requests:
    storage: 10Gi     # Збільш якщо потрібно більше
```

### 5. Backend Ingress

У `backend-ingress.yaml` зміни:
```yaml
spec:
  tls:
  - hosts:
    - api.yourdomain.com    # ⚠️ ТВІЙ DOMAIN!
  rules:
  - host: api.yourdomain.com  # ⚠️ ТВІЙ DOMAIN!
```

Налаштуй DNS:
```
api.yourdomain.com -> <ingress-external-ip>
```

### 6. Frontend ConfigMap

У `frontend-configmap.yaml`:
```yaml
data:
  API_URL: "https://api.yourdomain.com"  # ⚠️ ТВІЙ BACKEND URL!
  PUBLIC_URL: "https://yourdomain.com"   # ⚠️ ТВІЙ FRONTEND URL!
```

## 📊 Моніторинг

### Основні Команди

```bash
# Подивитися Pod'и
kubectl get pods -n orchestra-space -o wide

# Логи backend
kubectl logs -n orchestra-space -l app=orchestra-backend-production -f

# Логи PostgreSQL
kubectl logs -n orchestra-space -l app=postgres -f

# Ресурси (CPU/Memory)
kubectl top pods -n orchestra-space

# HPA status
kubectl get hpa -n orchestra-space --watch

# Events
kubectl get events -n orchestra-space --sort-by='.lastTimestamp'
```

### Health Checks

```bash
# Backend health
kubectl run -it --rm test --image=curlimages/curl --restart=Never -- \
  curl http://orchestra-backend-service-production/congratulation

# Database connection
kubectl exec -it -n orchestra-space <backend-pod> -- \
  psql -h orchestra-db -p 2828 -U rsmrtk -d orchestra-db -c "SELECT 1;"
```

## 🔄 Оновлення

### Rolling Update Backend

```bash
# Оновити image
kubectl set image deployment/orchestra-backend-production \
  orchestra-backend-production=martun5/orchestra-backend-production:v1.0.9 \
  -n orchestra-space

# Дивитися progress
kubectl rollout status deployment/orchestra-backend-production -n orchestra-space

# Якщо щось не так - rollback
kubectl rollout undo deployment/orchestra-backend-production -n orchestra-space
```

### Оновити ConfigMap

```bash
# 1. Відредагувати файл
vim configmap.yaml

# 2. Apply зміни
kubectl apply -f configmap.yaml

# 3. Restart Pod'ів (щоб підтягнути нові значення)
kubectl rollout restart deployment/orchestra-backend-production -n orchestra-space
```

## 🐛 Troubleshooting

### Pod не стартує

```bash
# Подивитися деталі
kubectl describe pod <pod-name> -n orchestra-space

# Типові проблеми:
# - ImagePullBackOff: неправильний image або немає доступу
# - CrashLoopBackOff: додаток падає (дивись логи)
# - Pending: недостатньо ресурсів
```

### Backend не може підключитися до БД

```bash
# Перевірити чи PostgreSQL Ready
kubectl get pods -n orchestra-space -l app=postgres

# Перевірити Secret
kubectl get secret orchestra-backend-secret -n orchestra-space -o yaml

# Перевірити Service endpoints
kubectl get endpoints orchestra-db -n orchestra-space

# Тестувати connection
kubectl exec -it -n orchestra-space <backend-pod> -- \
  nc -zv orchestra-db 2828
```

### HPA не працює

```bash
# Перевірити Metrics Server
kubectl top nodes
kubectl top pods -n orchestra-space

# Якщо <unknown>:
kubectl get deployment metrics-server -n kube-system
kubectl logs -n kube-system -l k8s-app=metrics-server
```

### NetworkPolicy блокує трафік

```bash
# Перевірити policies
kubectl get networkpolicy -n orchestra-space
kubectl describe networkpolicy backend-ingress-policy -n orchestra-space

# Тимчасово видалити для testing
kubectl delete networkpolicy backend-ingress-policy -n orchestra-space

# Тестувати connection
kubectl run -it --rm test --image=curlimages/curl -n orchestra-space -- \
  curl orchestra-backend-service-production/congratulation
```

## 📈 Масштабування

### Ручне Масштабування

```bash
# Збільшити кількість реплік
kubectl scale deployment/orchestra-backend-production --replicas=5 -n orchestra-space

# Перевірити
kubectl get pods -n orchestra-space -l app=orchestra-backend-production
```

### HPA Налаштування

Відредагуй `backend-hpa.yaml`:
```yaml
spec:
  minReplicas: 3        # Мінімум
  maxReplicas: 20       # Максимум
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        averageUtilization: 70  # Scale коли >70% CPU
```

## 🔒 Безпека

### Security Checklist

- [x] Паролі в Secret (не в plain text)
- [x] PVC для PostgreSQL (не emptyDir)
- [x] NetworkPolicy (ізоляція Pod'ів)
- [x] Readiness/Liveness probes
- [x] Resource limits (захист від resource exhaustion)
- [x] PDB (захист від downtime)
- [x] Ingress SSL/TLS
- [ ] RBAC (обмежити доступ до namespace)
- [ ] Pod Security Standards
- [ ] Network encryption (mTLS з service mesh)

## 📚 Додаткові Ресурси

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [HPA Walkthrough](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/)
- [NetworkPolicy Recipes](https://github.com/ahmetb/kubernetes-network-policy-recipes)
- [Production Best Practices](https://kubernetes.io/docs/setup/best-practices/)

---

**Створено для Orchestra Project** 🎵
Production-Ready Kubernetes Configuration
