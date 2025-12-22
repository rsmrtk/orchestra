# 🎓 Повний Розбір Kubernetes Deployment Файлів

## 📚 Зміст

1. [Що ми створили](#що-ми-створили)
2. [Архітектура системи](#архітектура-системи)
3. [Детальний розбір кожного файлу](#детальний-розбір)
4. [Як це все працює разом](#як-це-працює-разом)
5. [Практичні приклади](#практичні-приклади)
6. [Чек-лист готовності](#чек-лист-готовності)

---

## 🎯 Що ми створили

Ми створили **production-ready** Kubernetes конфігурацію з:

### Backend (14 файлів):
✅ **namespace.yaml** - Ізоляція ресурсів
✅ **secret.yaml** - Безпечне зберігання паролів
✅ **configmap.yaml** - Централізована конфігурація
✅ **postgres-pvc.yaml** - Постійне сховище для БД
✅ **postgres-deployment.yaml** - Розгортання PostgreSQL
✅ **postgres-service.yaml** - Доступ до БД
✅ **backend-deployment.yaml** - Розгортання API
✅ **backend-service.yaml** - Доступ до API
✅ **backend-hpa.yaml** - Автоматичне масштабування
✅ **backend-pdb.yaml** - Захист від downtime
✅ **backend-ingress.yaml** - Зовнішній HTTPS доступ
✅ **network-policy.yaml** - Мережева безпека

### Frontend (4 нових файли):
✅ **frontend-configmap.yaml** - Конфігурація frontend
✅ **frontend-hpa.yaml** - Автоскейлінг frontend
✅ **frontend-pdb.yaml** - Захист від downtime
✅ **frontend-network-policy.yaml** - Мережева безпека

---

## 🏗️ Архітектура системи

```
┌─────────────────────────────────────────────────────────┐
│                      INTERNET                            │
└───────────────────────┬─────────────────────────────────┘
                        │
                        ▼
            ┌───────────────────────┐
            │  Ingress Controller   │ ← SSL/TLS, Domain routing
            │  (nginx)              │
            └───────────┬───────────┘
                        │
         ┌──────────────┴──────────────┐
         │                             │
         ▼                             ▼
┌─────────────────┐          ┌─────────────────┐
│   Frontend      │          │    Backend      │
│   Pods (2-20)   │──────────▶    Pods (3-10)  │
│   Port: 3000    │  HTTP    │    Port: 8282   │
└─────────────────┘          └────────┬────────┘
         │                            │
         │ ❌ BLOCKED                 │ ✅ ALLOWED
         │                            │
         └────────────┬───────────────┘
                      ▼
              ┌───────────────┐
              │  PostgreSQL   │
              │  Pod (1)      │
              │  Port: 5432   │
              │  PVC: 10Gi    │
              └───────────────┘

NetworkPolicy Rules:
🟢 Internet → Ingress → Frontend ✅
🟢 Frontend → Backend ✅
🟢 Backend → PostgreSQL ✅
🔴 Frontend → PostgreSQL ❌ BLOCKED!
```

---

## 📖 Детальний розбір

### 1️⃣ Namespace (Фундамент)

**Файл:** `namespace.yaml`

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: orchestra-space
```

**Що це:**
- Віртуальний кластер всередині твого Kubernetes
- Ізолює ресурси від інших проектів

**Навіщо:**
- Production не перемішується з staging
- Можна встановити квоти (максимум 10 CPU, 32Gi RAM)
- RBAC - обмежити доступ тільки до production

**Аналогія:**
Namespace = окрема папка для твого проекту.

---

### 2️⃣ Secret (Безпека)

**Файл:** `secret.yaml`

```yaml
apiVersion: v1
kind: Secret
data:
  DB_PASSWORD: cnNtcnRr  # base64 encoded "rsmrtk"
```

**Що це:**
- Зберігає паролі, токени, API keys
- Base64 encoding (НЕ encryption!)

**⚠️ КРИТИЧНО:**
```bash
# НЕ КОМІТИТИ цей файл з реальними паролями!
# Створюй Secret вручну:
kubectl create secret generic orchestra-backend-secret \
  --from-literal=DB_PASSWORD=YourStrongPassword123! \
  --namespace=orchestra-space
```

**Як використовується:**
```yaml
env:
  - name: DB_PASSWORD
    valueFrom:
      secretKeyRef:
        name: orchestra-backend-secret
        key: DB_PASSWORD
```

Kubernetes автоматично декодує base64 і встановлює як env змінну.

---

### 3️⃣ ConfigMap (Конфігурація)

**Файл:** `configmap.yaml`

```yaml
apiVersion: v1
kind: ConfigMap
data:
  GIN_MODE: "release"
  LOG_LEVEL: "info"
  APP_PORT: "8282"
```

**Різниця між Secret та ConfigMap:**
| Secret | ConfigMap |
|--------|-----------|
| Паролі, токени | Порти, URLs |
| Base64 encoded | Plain text |
| Обережно! | Публічно OK |

**Чому відокремлювати config від коду:**
- Змінити налаштування без rebuild Docker image
- Різна конфігурація для staging/production
- Централізоване управління

---

### 4️⃣ PersistentVolumeClaim (Збереження даних)

**Файл:** `postgres/postgres-pvc.yaml`

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
spec:
  resources:
    requests:
      storage: 10Gi
```

**⚠️ НАЙВАЖЛИВІШИЙ ФАЙЛ ДЛЯ БД!**

**Без PVC:**
```yaml
volumes:
  - name: postgres-storage
    emptyDir: {}  # ❌ Всі дані ВТРАТЯТЬСЯ при restart!
```

**З PVC:**
```yaml
volumes:
  - name: postgres-storage
    persistentVolumeClaim:
      claimName: postgres-pvc-production  # ✅ Дані збережуться!
```

**Що відбувається:**
1. Створюєш PVC (запит на 10Gi диска)
2. Kubernetes створює PersistentVolume на Node
3. Pod монтує PVC як `/var/lib/postgresql/data`
4. PostgreSQL зберігає дані на диску
5. Pod перезапускається → Дані залишаються!

**Reclaim Policy:**
- `Retain` - Дані зберігаються після видалення PVC ✅ (безпечно)
- `Delete` - Дані видаляються автоматично ❌ (небезпечно!)

---

### 5️⃣ Deployment (Розгортання Pod'ів)

**Файл:** `backend/backend-deployment.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  replicas: 3  # Скільки Pod'ів створити
```

**Lifecycle Pod'а:**
```
1. Scheduler → Вибирає Node де достатньо ресурсів
2. Kubelet → Pull Docker image
3. Kubelet → Створює контейнер
4. Startup Probe → Чи додаток стартував? (60 сек timeout)
5. Readiness Probe → Чи готовий приймати трафік?
6. Liveness Probe → Чи живий? (постійно перевіряє)
```

**Probes - Різниця:**
| Probe | Коли | Якщо fails |
|-------|------|------------|
| Startup | При старті (одноразово) | Pod killed → restart |
| Readiness | Постійно | Pod excluded з Service |
| Liveness | Постійно | Container killed → restart |

**Resources:**
```yaml
resources:
  requests:  # Гарантований мінімум
    memory: "256Mi"
    cpu: "250m"
  limits:    # Максимум
    memory: "512Mi"
    cpu: "500m"
```

- Якщо перевищує `limits.memory` → OOMKilled
- Якщо перевищує `limits.cpu` → Throttled (повільніше працює)

---

### 6️⃣ Service (Доступ до Pod'ів)

**Файл:** `backend/backend-service.yaml`

```yaml
apiVersion: v1
kind: Service
spec:
  selector:
    app: orchestra-backend-production
  ports:
  - port: 80
    targetPort: 8282
```

**Чому потрібен Service:**
- Pod IP змінюється після restart
- Service має стабільну DNS адресу: `orchestra-backend-service`
- Load balancing між Pod'ами

**Як працює:**
```
Request → Service:80 → Load Balancer → Pod1:8282 або Pod2:8282 або Pod3:8282
```

**Types:**
- **ClusterIP** - Тільки всередині кластера (для БД)
- **NodePort** - Доступний ззовні через Node IP (для dev)
- **LoadBalancer** - Cloud load balancer (для production)
- **Ingress + ClusterIP** - Найкраще для production! ✅

---

### 7️⃣ HPA (Автоскейлінг)

**Файл:** `backend/backend-hpa.yaml`

```yaml
spec:
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        averageUtilization: 70  # Scale коли >70% CPU
```

**Як це працює:**

```
1. Поточно: 3 Pod'и, CPU 80%
2. HPA перевіряє кожні 15 сек
3. Розраховує: desiredReplicas = ceil[3 * (80/70)] = 4
4. Чекає 60 сек (stabilization)
5. Створює 4-й Pod
6. Тепер: 4 Pod'и, CPU 60% < 70% ✅
```

**Поведінка:**
- **Scale Up** - Швидко (30-60 сек stabilization)
- **Scale Down** - Повільно (5 хв stabilization)

Чому? Краще мати зайві Pod'и ніж downtime!

**Frontend vs Backend HPA:**
| | Frontend | Backend |
|--|----------|---------|
| minReplicas | 2 | 3 |
| maxReplicas | 20 | 10 |
| CPU threshold | 60% | 70% |
| Scale up stabilization | 30 сек | 60 сек |

Frontend масштабується агресивніше (users не чекають!).

---

### 8️⃣ PDB (Захист від Downtime)

**Файл:** `backend/backend-pdb.yaml`

```yaml
spec:
  minAvailable: 2  # Мінімум 2 Pod'и завжди працюють
```

**Що контролює PDB:**
✅ `kubectl drain node` (maintenance)
✅ Cluster autoscaler scale down
✅ Rolling updates

**Що НЕ контролює:**
❌ Node crash (hardware failure)
❌ OOMKilled
❌ `kubectl delete pod --force`

**Приклад (Rolling Update з 3 репліками):**
```
1. Старт: 3 Pod'i v1.0.8
2. Створюється 1 Pod v1.0.9 → Тепер 4 Pod'i
3. PDB check: 4 - 1 = 3 ≥ 2 (minAvailable) ✅
4. Видаляється 1 Pod v1.0.8 → Тепер 3 Pod'i
5. Повторюємо поки всі не оновляться
```

PDB гарантував що **завжди мінімум 2 Pod'и доступні** під час update!

---

### 9️⃣ Ingress (Зовнішній Доступ)

**Файл:** `backend/backend-ingress.yaml`

```yaml
spec:
  tls:
  - hosts:
    - api.yourdomain.com
    secretName: orchestra-backend-tls
  rules:
  - host: api.yourdomain.com
    http:
      paths:
      - path: /
        backend:
          service:
            name: orchestra-backend-service-production
```

**Request Flow:**
```
User → https://api.yourdomain.com/data
  ↓
DNS → Ingress External IP (35.123.45.67)
  ↓
Ingress Controller (nginx)
  ├─ SSL termination (HTTPS → HTTP)
  ├─ Host matching: api.yourdomain.com ✅
  └─ Path matching: / ✅
  ↓
Service: orchestra-backend-service:80
  ↓
Load Balancer → Pod:8282
  ↓
Backend додаток
```

**Переваги Ingress:**
✅ Domain name (api.yourdomain.com)
✅ SSL/TLS (HTTPS)
✅ Path-based routing (/api/v1, /api/v2)
✅ Rate limiting, auth
✅ Один IP для багатьох Services

---

### 🔟 NetworkPolicy (Безпека)

**Файл:** `network-policy.yaml`

```yaml
# Backend може приймати від:
ingress:
  - from:
    - podSelector:
        matchLabels:
          app: orchestra-frontend-production
    ports:
    - port: 8282

# Backend може викликати:
egress:
  - to:
    - podSelector:
        matchLabels:
          app: postgres
    ports:
    - port: 5432
```

**3-Tier Security Model:**
```
Internet
  ↓ ✅
Ingress Controller
  ↓ ✅
Frontend (може тільки Backend)
  ↓ ✅
Backend (може тільки PostgreSQL)
  ↓ ✅
PostgreSQL (ТІЛЬКИ від Backend!)
```

**⛔ Frontend → PostgreSQL BLOCKED!**
Якщо хакер зламає frontend, не зможе підключитися до БД.

---

## 🔄 Як це все працює разом

### Scenario 1: User відкриває сайт

```
1. User → https://yourdomain.com
2. DNS → Ingress IP
3. Ingress Controller → Frontend Service:80
4. Service → Frontend Pod:3000 (HPA: 2-20 Pod'ів)
5. Frontend serve HTML/JS
6. Response → User
```

### Scenario 2: User робить API call

```
1. User clicks button → JavaScript fetch('/api/data')
2. Frontend Pod → Backend Service:80
3. NetworkPolicy check: Frontend → Backend? ✅ ALLOWED
4. Service → Backend Pod:8282 (HPA: 3-10 Pod'ів)
5. Backend Pod → PostgreSQL Service:2828
6. NetworkPolicy check: Backend → PostgreSQL? ✅ ALLOWED
7. Service → PostgreSQL Pod:5432
8. PostgreSQL query → PVC (persistent data)
9. Response → Backend → Frontend → User
```

### Scenario 3: Traffic spike (HPA scales)

```
1. Normal: 3 backend Pod'і, CPU 40%
2. Marketing campaign → Traffic x10
3. CPU spike → 90%
4. HPA: 90% > 70% (target) → Scale up!
5. HPA розраховує: ceil[3 * (90/70)] = 4 Pod'і
6. Чекає 60 сек (stabilization)
7. Створює 4-й Pod
8. Deployment → ReplicaSet → 4-й Pod створений
9. Startup probe → Readiness probe → Pod Ready
10. Service додає Pod IP в endpoints
11. Тепер: 4 Pod'і, CPU 67% < 70% ✅
12. Traffic спадає через 10 хв
13. HPA чекає 5 хв (stabilization)
14. Scale down до 3 Pod'ів
```

### Scenario 4: Rolling Update (Zero Downtime)

```
1. Старт: 3 Pod'i v1.0.8 (Ready)
2. kubectl set image ... v1.0.9
3. Deployment створює новий ReplicaSet
4. maxSurge: 1 → Створюється 1 Pod v1.0.9
5. Тепер: 3 old + 1 new = 4 Pod'i
6. Startup + Readiness probes → новий Pod Ready
7. Service додає новий Pod в endpoints
8. PDB check: 4 - 1 = 3 ≥ 2 ✅ можна видалити old Pod
9. maxUnavailable: 1 → Видаляється 1 Pod v1.0.8
10. Тепер: 2 old + 1 new = 3 Pod'i
11. Повторюємо кроки 4-10 поки всі не оновляться
12. Фініш: 3 Pod'i v1.0.9 (Ready)
```

**Результат:** Завжди мінімум 2 Pod'і працювали!

---

## 💡 Практичні Приклади

### 1. Дізнатися чому Pod падає

```bash
# Статус Pod'а
kubectl get pod <pod-name> -n orchestra-space

# Детальна інформація
kubectl describe pod <pod-name> -n orchestra-space

# Логи
kubectl logs <pod-name> -n orchestra-space

# Логи попереднього контейнера (якщо перезапускався)
kubectl logs <pod-name> -n orchestra-space --previous

# Events (останні помилки)
kubectl get events -n orchestra-space --sort-by='.lastTimestamp' | grep <pod-name>
```

### 2. Debug мережевих проблем

```bash
# Перевірити Service endpoints
kubectl get endpoints orchestra-backend-service -n orchestra-space

# Якщо порожньо → перевір selector:
kubectl get pods -n orchestra-space --show-labels

# Тестувати connection з Pod'а
kubectl exec -it <pod-name> -n orchestra-space -- curl http://orchestra-db:2828

# Перевірити NetworkPolicy
kubectl get networkpolicy -n orchestra-space
kubectl describe networkpolicy backend-ingress-policy -n orchestra-space
```

### 3. Моніторинг ресурсів

```bash
# CPU/Memory Pod'ів
kubectl top pods -n orchestra-space

# HPA status
kubectl get hpa -n orchestra-space

# Якщо <unknown>:
kubectl get deployment metrics-server -n kube-system
```

### 4. Backup PostgreSQL

```bash
# Manual backup
kubectl exec -n orchestra-space <postgres-pod> -- \
  pg_dump -U rsmrtk orchestra-db > backup-$(date +%Y%m%d).sql

# Restore
kubectl exec -i -n orchestra-space <postgres-pod> -- \
  psql -U rsmrtk -d orchestra-db < backup-20241222.sql
```

---

## ✅ Чек-лист Готовності

### Must Have (Обов'язково)

- [x] ✅ **PVC для PostgreSQL** (не emptyDir!)
- [x] ✅ **Secret для паролів** (не plain text!)
- [x] ✅ **Resources для всіх Pod'ів** (requests та limits)
- [x] ✅ **Readiness та Liveness probes**
- [x] ✅ **HPA для backend та frontend**
- [x] ✅ **PDB для захисту від downtime**
- [x] ✅ **NetworkPolicy для ізоляції**

### Should Have (Дуже рекомендовано)

- [x] ✅ **Ingress з SSL/TLS**
- [x] ✅ **ConfigMap для конфігурації**
- [x] ✅ **Namespace для ізоляції**
- [ ] ⚠️ **Metrics Server** (перевір: `kubectl top nodes`)
- [ ] ⚠️ **Ingress Controller** (nginx, traefik)
- [ ] ⚠️ **cert-manager** для автоматичного SSL

### Nice to Have

- [ ] 📊 Prometheus + Grafana (моніторинг)
- [ ] 🚨 Alertmanager (alerts)
- [ ] 📝 Loki (centralized logging)
- [ ] 🔒 RBAC (обмежити доступ)
- [ ] 🔐 Pod Security Standards
- [ ] 🌐 Service Mesh (Istio/Linkerd)

---

## 🎓 Що далі вивчати

### 1. GitOps (Наступний крок!)
- **ArgoCD** або **Flux**
- Автоматичний deploy з Git
- Rollback одним кліком
- Git = single source of truth

### 2. Observability
- **Prometheus** - metrics
- **Grafana** - dashboards
- **Loki** - logs
- **Jaeger/Tempo** - tracing

### 3. Security
- **Sealed Secrets** - encrypted secrets в Git
- **External Secrets Operator** - secrets з Vault
- **Falco** - runtime security
- **OPA** - policy enforcement

### 4. Advanced Patterns
- **StatefulSet** - для stateful apps
- **DaemonSet** - pod на кожному Node
- **CronJob** - scheduled tasks
- **Operators** - custom controllers

### 5. Multi-Environment
- **Kustomize** - manage configurations
- **Helm** - package manager
- **Overlays** - base + environment-specific

---

## 📚 Ресурси для навчання

### Офіційна документація:
- [Kubernetes Docs](https://kubernetes.io/docs/)
- [kubectl Cheat Sheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)

### Інтерактивне навчання:
- [Katacoda](https://www.katacoda.com/courses/kubernetes)
- [Play with Kubernetes](https://labs.play-with-k8s.com/)

### Книги:
- "Kubernetes Up & Running" (Kelsey Hightower)
- "Kubernetes Patterns" (Bilgin Ibryam)

### YouTube:
- [That DevOps Guy](https://www.youtube.com/@MarcelDempers)
- [TechWorld with Nana](https://www.youtube.com/@TechWorldwithNana)

---

## 🎉 Висновок

Ти створив **production-ready** Kubernetes deployment з:
- ✅ Високою доступністю (replicas, PDB)
- ✅ Автоскейлінгом (HPA)
- ✅ Безпекою (NetworkPolicy, Secrets)
- ✅ Надійністю (PVC, Probes)
- ✅ Zero-downtime updates (RollingUpdate)

**Твій рівень:** 7/10 - Готовий до production! 🚀

**Наступні кроки:**
1. Задеплой на реальний кластер
2. Налаштуй моніторинг
3. Додай CI/CD (GitHub Actions)
4. Вивчи GitOps (ArgoCD)

**Успіхів у навчанні! 🎓**
