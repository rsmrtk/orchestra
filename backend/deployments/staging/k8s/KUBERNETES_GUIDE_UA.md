# 🎓 Повний гайд по Kubernetes для початківців (Українською)

## 📖 Зміст

1. [Що таке Kubernetes і навіщо він потрібен](#що-таке-kubernetes)
2. [Основні концепції та архітектура](#основні-концепції)
3. [Як працюють YAML файли](#як-працюють-yaml-файли)
4. [Детальний розбір кожного ресурсу](#детальний-розбір)
5. [Як все це працює разом](#як-все-працює-разом)
6. [Практичні сценарії](#практичні-сценарії)
7. [Що очікувати і як тестувати](#що-очікувати)

---

## 🤔 Що таке Kubernetes?

### Проста аналогія

Уяви що твій додаток - це група робітників:
- **Docker** - це будинок для одного робітника (контейнер для одного процесу)
- **Kubernetes** - це весь завод з менеджерами, які керують робітниками:
  - Наймають нових робітників (створюють контейнери)
  - Звільняють тих що не працюють (перезапускають контейнери)
  - Розподіляють роботу (load balancing)
  - Слідкують щоб завжди була потрібна кількість (auto-healing)
  - Збільшують/зменшують команду під навантаження (auto-scaling)

### Реальний приклад

**БЕЗ Kubernetes:**
```
Твій сервер:
- 1 Docker контейнер з backend
- Якщо впав - треба вручну перезапустити
- Якщо багато трафіку - треба вручну запустити більше
- Якщо сервер впав - все зламалося
```

**З Kubernetes:**
```
Кластер (група серверів):
- Kubernetes автоматично запускає 3 копії backend (high availability)
- Якщо одна впала - автоматично перезапускає
- Якщо багато трафіку - автоматично запускає більше копій
- Якщо сервер впав - переносить контейнери на інший сервер
- Один IP для доступу, навіть якщо контейнери на різних серверах
```

---

## 🏗 Основні концепції

### 1. Кластер = Група серверів

```
Kubernetes Cluster
├── Master Node (Control Plane) - мозок кластера
│   ├── API Server - приймає команди (kubectl)
│   ├── Scheduler - вирішує де запустити Pod'и
│   ├── Controller Manager - слідкує щоб все працювало
│   └── etcd - база даних кластера
│
└── Worker Nodes - робочі коні, де працюють додатки
    ├── Node 1
    │   ├── Pod 1 (твій backend контейнер)
    │   ├── Pod 2 (ще один backend контейнер)
    │   └── Pod 3 (інший сервіс)
    ├── Node 2
    │   ├── Pod 4
    │   └── Pod 5
    └── Node 3
        └── Pod 6
```

### 2. Основні об'єкти (від найменшого до найбільшого)

```
Container (Docker) - один процес в ізоляції
    ↓
Pod - одна або більше контейнерів разом (найменша одиниця в K8s)
    ↓
ReplicaSet - гарантує що N Pod'ів працює
    ↓
Deployment - управляє ReplicaSet (оновлення, rollback)
    ↓
Service - стабільна точка доступу до Pod'ів
    ↓
Ingress - HTTP маршрутизація ззовні кластера
```

### 3. Архітектура твого додатку в K8s

```
Інтернет
    ↓
[Ingress Controller] - nginx який приймає HTTPS
    ↓ (маршрутизує за domain/path)
[Service] - Load Balancer між Pod'ами
    ↓ (розподіляє трафік)
[Pod 1]  [Pod 2]  [Pod 3] - твої backend контейнери
    ↓ (читають конфігурацію)
[ConfigMap]  [Secret] - конфігурація та секрети
```

---

## 📝 Як працюють YAML файли

### Структура будь-якого Kubernetes YAML

Кожен YAML файл має 4 основні секції:

```yaml
apiVersion: v1              # 1. ВЕРСІЯ API - яку версію Kubernetes API використати
kind: Pod                   # 2. ТИП РЕСУРСУ - що ми створюємо
metadata:                   # 3. МЕТАДАНІ - інформація ПРО ресурс
  name: my-pod              #    - ім'я (унікальне в namespace)
  namespace: default        #    - в якому namespace
  labels:                   #    - мітки для організації
    app: backend
    env: staging
spec:                       # 4. СПЕЦИФІКАЦІЯ - ЯК має виглядати ресурс
  # Тут вся конфігурація (різна для кожного типу)
```

### Як Kubernetes читає YAML

```
1. Ти пишеш: kubectl apply -f deployment.yaml
           ↓
2. kubectl відправляє YAML до API Server
           ↓
3. API Server зберігає в etcd (база даних)
           ↓
4. Controller Manager бачить новий Deployment
           ↓
5. Створює ReplicaSet
           ↓
6. ReplicaSet створює Pod'и
           ↓
7. Scheduler вибирає на якій Node запустити кожен Pod
           ↓
8. Kubelet на Node'і завантажує Docker образ
           ↓
9. Docker запускає контейнер
           ↓
10. Pod працює! 🎉
```

### Декларативний підхід

**Імперативний** (як в Docker Compose):
```bash
docker run -d -p 8080:8080 my-app  # Ти кажеш ЯК запустити
docker stop my-app                  # Ти кажеш ЯК зупинити
```

**Декларативний** (в Kubernetes):
```yaml
# deployment.yaml - ти описуєш ЩО ти хочеш
spec:
  replicas: 3  # Я хочу 3 Pod'и
```

Kubernetes сам:
- Створює 3 Pod'и
- Якщо один впав - створить новий
- Якщо їх 5 - видалить зайві
- Якщо 0 - створить 3

**Ти не кажеш ЯК це зробити, ти кажеш ЩО ти хочеш отримати!**

---

## 🔍 Детальний розбір кожного ресурсу

### 1. Namespace - Віртуальний кластер

**Що це простими словами:**
Уяви що у тебе одна квартира (кластер), але ти хочеш відокремити кухню (staging) від спальні (production). Namespace - це віртуальні кімнати.

**Навіщо:**
```
orchestra-staging/     - тут все для тестування
  ├── backend Pod'и
  ├── database
  └── redis

orchestra-production/  - тут реальні користувачі
  ├── backend Pod'и
  ├── database
  └── redis

Вони НЕ бачать один одного!
```

**YAML структура:**
```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: orchestra-staging    # Ім'я namespace
  labels:                    # Мітки для організації
    environment: staging
```

**Що очікувати:**
- Після `kubectl apply -f namespace.yaml` створюється пустий namespace
- Всі інші ресурси створюються всередині цього namespace
- Можна видалити весь namespace і все всередині видалиться

---

### 2. ConfigMap - Конфігурація додатку

**Що це простими словами:**
ConfigMap - це файл з налаштуваннями (як `.env` файл), тільки в Kubernetes.

**Навіщо:**
```
БЕЗ ConfigMap:
- Налаштування в Docker образі
- Щоб змінити - треба пересобрати образ
- Один образ = одне середовище

З ConfigMap:
- Налаштування окремо від образу
- Змінити можна без пересборки
- Один образ для staging і production (різні ConfigMap)
```

**YAML структура:**
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-config
  namespace: my-namespace
data:                        # Секція з даними
  APP_PORT: "8383"          # key: value (все в лапках!)
  GIN_MODE: "debug"
  DATABASE_HOST: "postgres-service"

  # Можна цілий файл:
  app.properties: |
    server.port=8383
    server.name=backend
```

**Як використовується в Pod'і:**
```yaml
# В deployment.yaml
env:
- name: APP_PORT              # Ім'я env змінної в контейнері
  valueFrom:
    configMapKeyRef:
      name: my-config         # Ім'я ConfigMap
      key: APP_PORT           # Ключ з ConfigMap
```

**Що очікувати:**
- Pod читає `APP_PORT` як env змінну
- Якщо змінити ConfigMap - треба перезапустити Pod'и
- Можна монтувати як файли замість env змінних

---

### 3. Secret - Секретна інформація

**Що це простими словами:**
Secret - те саме що ConfigMap, але для паролів, токенів, ключів.

**ВАЖЛИВО:**
```
Secret != безпека!
- Дані в base64 (це НЕ шифрування!)
- echo "password" | base64  →  cGFzc3dvcmQK
- echo "cGFzc3dvcmQK" | base64 -d  →  password

Kubernetes МОЖЕ шифрувати Secrets at rest (в etcd), але це треба налаштувати.
```

**YAML структура:**
```yaml
apiVersion: v1
kind: Secret
metadata:
  name: my-secret
  namespace: my-namespace
type: Opaque                 # Тип: Opaque, kubernetes.io/tls, kubernetes.io/dockerconfigjson
data:                        # Дані в base64!
  API_KEY: bXktc2VjcmV0      # "my-secret" в base64
  DB_PASSWORD: cGFzc3dvcmQ=  # "password" в base64
```

**Як закодувати:**
```bash
# Закодувати
echo -n "my-secret-key" | base64
# Результат: bXktc2VjcmV0LWtleQ==

# Декодувати
echo "bXktc2VjcmV0LWtleQ==" | base64 -d
# Результат: my-secret-key
```

**Краще створювати через kubectl:**
```bash
kubectl create secret generic my-secret \
  --from-literal=API_KEY=my-secret-key \
  --from-literal=DB_PASSWORD=password123 \
  -n my-namespace
```

**Що очікувати:**
- Kubernetes автоматично декодує base64 при використанні
- НЕ комітити Secret з реальними даними в Git!
- Використовуй Sealed Secrets, External Secrets, Vault

---

### 4. Deployment - Управління Pod'ами

**Що це простими словами:**
Deployment - це менеджер який слідкує щоб твої контейнери працювали.

**Як це працює:**
```
Ти кажеш: "Я хочу 3 Pod'и з моїм backend"
           ↓
Deployment створює ReplicaSet
           ↓
ReplicaSet створює 3 Pod'и
           ↓
Якщо Pod впав:
  ReplicaSet бачить що тепер 2 Pod'и (не 3)
           ↓
  ReplicaSet створює новий Pod
           ↓
  Знову 3 Pod'и!
```

**Ключові частини YAML:**

#### 4.1. Скільки Pod'ів запустити
```yaml
spec:
  replicas: 3                # Кількість копій твого додатку
```

#### 4.2. Як оновлювати
```yaml
spec:
  strategy:
    type: RollingUpdate      # Поступово міняти старі Pod'и на нові
    rollingUpdate:
      maxUnavailable: 1      # Максимум 1 Pod може бути недоступний
      maxSurge: 1            # Максимум 1 додатковий Pod під час оновлення
```

**Приклад Rolling Update:**
```
Поточно: 3 Pod'и (version 1.0)
[Pod-v1] [Pod-v1] [Pod-v1]

Крок 1: Створити 1 новий Pod (maxSurge: 1)
[Pod-v1] [Pod-v1] [Pod-v1] [Pod-v2] ← новий

Крок 2: Дочекатися поки новий Pod готовий
[Pod-v1] [Pod-v1] [Pod-v1] [Pod-v2✓]

Крок 3: Видалити 1 старий Pod (maxUnavailable: 1)
[Pod-v1] [Pod-v1] [Pod-v2✓]

Крок 4: Створити ще 1 новий Pod
[Pod-v1] [Pod-v1] [Pod-v2✓] [Pod-v2]

...продовжується поки всі не оновляться...

Результат: 3 Pod'и (version 2.0), без downtime!
[Pod-v2] [Pod-v2] [Pod-v2]
```

#### 4.3. Який Docker образ використати
```yaml
spec:
  template:
    spec:
      containers:
      - name: backend
        image: ghcr.io/user/backend:v1.0.0  # Registry/Image:Tag
        imagePullPolicy: Always              # Завжди тягнути новий образ
```

#### 4.4. Ресурси (CPU, Memory)
```yaml
resources:
  requests:              # МІНІМУМ - для scheduler (де розмістити Pod)
    memory: "128Mi"      # 128 мебібайт RAM
    cpu: "100m"          # 100 міллікор = 0.1 CPU core
  limits:                # МАКСИМУМ - для захисту (щоб не з'їв всі ресурси)
    memory: "256Mi"      # Якщо перевищить - Pod буде killed (OOMKilled)
    cpu: "200m"          # Якщо перевищить - CPU буде throttled (повільніше)
```

**Приклад:**
```
requests: 100m CPU, 128Mi RAM
limits: 200m CPU, 256Mi RAM

Що це означає:
- Kubernetes гарантує 0.1 CPU core і 128MB RAM
- Pod може використати до 0.2 CPU core і 256MB RAM
- Якщо використає > 256MB - Kubernetes його killed
- Якщо використає > 0.2 CPU - Kubernetes обмежить швидкість
```

#### 4.5. Health Checks (Перевірки здоров'я)

**Liveness Probe** - чи живий контейнер (чи не завис):
```yaml
livenessProbe:
  httpGet:
    path: /health          # HTTP GET запит до цього endpoint
    port: 8383
  initialDelaySeconds: 30  # Почекати 30 сек після старту (дати час запуститися)
  periodSeconds: 10        # Перевіряти кожні 10 секунд
  timeoutSeconds: 5        # Таймаут запиту 5 секунд
  failureThreshold: 3      # 3 невдачі підряд = Pod dead, перезапустити
```

**Що відбувається:**
```
Pod запускається
   ↓
30 секунд чекання (initialDelaySeconds)
   ↓
Кожні 10 секунд: GET http://pod-ip:8383/health
   ↓
Якщо 3 рази підряд failed (timeout або HTTP 500+):
   ↓
Kubernetes: "Pod завис! Перезапускаю..."
   ↓
Pod перезапускається
```

**Readiness Probe** - чи готовий приймати трафік:
```yaml
readinessProbe:
  httpGet:
    path: /health
    port: 8383
  initialDelaySeconds: 5   # Швидше ніж liveness
  periodSeconds: 5
  failureThreshold: 3
```

**Різниця liveness vs readiness:**
```
Liveness failed → Kubernetes перезапускає Pod
Readiness failed → Kubernetes НЕ перезапускає, просто не шле трафік на цей Pod

Приклад:
- 3 Pod'и: [Pod1✓] [Pod2✓] [Pod3✓]
- Pod2 перевантажений, readiness failed
- Тепер: [Pod1✓] [Pod2✗] [Pod3✓]
- Трафік йде тільки на Pod1 і Pod3
- Pod2 не перезапускається, просто "відпочиває"
- Коли Pod2 відновився - знову приймає трафік
```

---

### 5. Service - Мережевий доступ

**Що це простими словами:**
Service - це як телефонна книга. Замість того щоб запам'ятовувати IP кожного Pod'а (вони змінюються!), ти дзвониш на один номер (Service), а він перенаправляє на вільного "оператора" (Pod).

**Проблема без Service:**
```
Pod'и мають динамічні IP:
- Pod1: 10.244.1.5
- Pod2: 10.244.2.8
- Pod3: 10.244.3.12

Якщо Pod1 перезапустився:
- Pod1: 10.244.1.99  ← нова IP!

Як frontend знає куди звертатися? 🤔
```

**Рішення з Service:**
```
Service має:
- Стабільне ім'я: backend-service
- Стабільну IP: 10.96.0.50 (Cluster IP)
- DNS: backend-service.namespace.svc.cluster.local

Frontend завжди звертається: http://backend-service:80
Service сам знаходить Pod'и і розподіляє трафік!
```

**Як Service знаходить Pod'и:**
```yaml
# Service YAML
spec:
  selector:                  # Шукати Pod'и з цими labels
    app: backend

# Deployment YAML
spec:
  template:
    metadata:
      labels:                # Pod'и мають ці labels
        app: backend

Селектори МАЮТЬ СПІВПАДАТИ!
```

**Типи Service:**

#### ClusterIP (default) - тільки всередині кластера
```yaml
spec:
  type: ClusterIP
  ports:
  - port: 80                 # Порт Service (на якому слухає)
    targetPort: 8383         # Порт контейнера (куди перенаправляє)
```

```
Доступ:
- Всередині кластера: http://backend-service:80 ✓
- Ззовні кластера: ✗ (недоступно)

Використання: backend ↔ backend, backend ↔ database
```

#### NodePort - доступний на кожній Node
```yaml
spec:
  type: NodePort
  ports:
  - port: 80
    targetPort: 8383
    nodePort: 30080          # Порт на кожній Node (30000-32767)
```

```
Доступ:
- Всередині: http://backend-service:80 ✓
- Ззовні: http://<будь-яка-node-ip>:30080 ✓

Використання: dev/test середовища, коли нема LoadBalancer
```

#### LoadBalancer - зовнішній Load Balancer
```yaml
spec:
  type: LoadBalancer
  ports:
  - port: 80
    targetPort: 8383
```

```
Cloud provider (AWS/GCP/Azure) створює зовнішній LB
External IP: 34.123.45.67

Доступ:
- Ззовні: http://34.123.45.67:80 ✓

Використання: production, публічні API
```

**DNS імена:**
```
<service-name>.<namespace>.svc.cluster.local

Приклади:
backend-service                                    (той самий namespace)
backend-service.orchestra-staging                  (з іншого namespace)
backend-service.orchestra-staging.svc.cluster.local (повне FQDN)
```

---

### 6. Ingress - HTTP/HTTPS маршрутизація

**Що це простими словами:**
Ingress - це охоронець на вході в кластер. Він дивиться на URL і вирішує куди направити запит.

**Проблема без Ingress:**
```
Для кожного сервісу потрібен LoadBalancer:
- backend-lb: 34.123.45.67  (коштує $$)
- frontend-lb: 34.123.45.68 (коштує $$)
- api-lb: 34.123.45.69      (коштує $$)

Дорого і незручно!
```

**Рішення з Ingress:**
```
Один LoadBalancer: 34.123.45.67

Ingress маршрутизує:
api.example.com/users    → backend-service
api.example.com/posts    → blog-service
app.example.com          → frontend-service
admin.example.com        → admin-service

Один IP, багато сервісів!
```

**Як працює:**
```
1. Встановлюєш Ingress Controller (nginx, traefik)
           ↓
2. Ingress Controller створює LoadBalancer
           ↓
3. Створюєш Ingress ресурси з правилами
           ↓
4. Ingress Controller читає правила і налаштовує nginx
           ↓
5. Трафік: Інтернет → LB → Ingress Controller → Service → Pod
```

**YAML структура:**
```yaml
spec:
  rules:
  - host: api.example.com            # Домен
    http:
      paths:
      - path: /api                   # Шлях
        pathType: Prefix             # Prefix = /api, /api/users, /api/posts
        backend:
          service:
            name: backend-service    # Куди перенаправити
            port:
              number: 80
```

**Приклади маршрутизації:**
```yaml
# Маршрутизація по host (domain)
rules:
- host: api.example.com              # → backend-service
- host: app.example.com              # → frontend-service
- host: admin.example.com            # → admin-service

# Маршрутизація по path
rules:
- host: example.com
  paths:
  - path: /api                       # → backend-service
  - path: /app                       # → frontend-service
  - path: /admin                     # → admin-service
```

**SSL/TLS (HTTPS):**
```yaml
spec:
  tls:
  - hosts:
    - api.example.com
    secretName: api-tls              # Secret з SSL сертифікатом
```

**Автоматичний SSL з cert-manager:**
```yaml
metadata:
  annotations:
    cert-manager.io/cluster-issuer: "letsencrypt-prod"

# cert-manager автоматично отримає SSL від Let's Encrypt!
```

---

### 7. HPA - Автоматичне масштабування

**Що це простими словами:**
HPA - це автоматичний "найомщик" Pod'ів. Коли багато роботи - наймає більше, коли мало - звільняє зайвих.

**Як працює:**
```
1. Metrics Server збирає метрики (CPU, Memory) кожні 15 секунд
           ↓
2. HPA Controller перевіряє метрики
           ↓
3. Порівнює з target (наприклад 70% CPU)
           ↓
4. Якщо вище - збільшує replicas
   Якщо нижче - зменшує replicas
           ↓
5. Оновлює Deployment.spec.replicas
           ↓
6. Deployment створює/видаляє Pod'и
```

**YAML структура:**
```yaml
spec:
  scaleTargetRef:
    kind: Deployment
    name: backend               # Який Deployment масштабувати

  minReplicas: 2                # Мінімум Pod'ів (завжди)
  maxReplicas: 10               # Максимум Pod'ів (не більше)

  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70  # Тримати середній CPU ~70%
```

**Приклад роботи:**
```
Поточно: 2 Pod'и, CPU 80% (вище target 70%)

HPA обчислює:
desiredReplicas = 2 * (80 / 70) = 2.28 ≈ 3 Pod'и

HPA збільшує до 3 Pod'ів:
[Pod1] [Pod2] [Pod3]

Тепер: 3 Pod'и, CPU 53% (нижче 70%, все ОК)

---

Трафік зменшився: 3 Pod'и, CPU 20%

HPA обчислює:
desiredReplicas = 3 * (20 / 70) = 0.85 ≈ 1 Pod

Але minReplicas: 2, тому:
HPA зменшує до 2 Pod'ів (не нижче мінімуму)
[Pod1] [Pod2]
```

**Важливо:**
- Треба вказати `resources.requests` в Deployment
- Треба встановити Metrics Server
- HPA НЕ може масштабувати нижче minReplicas або вище maxReplicas

---

### 8. NetworkPolicy - Firewall для Pod'ів

**Що це простими словами:**
NetworkPolicy - це правила хто може з'єднуватися до твоїх Pod'ів і куди твої Pod'и можуть з'єднуватися.

**По замовчуванню в Kubernetes:**
```
Всі Pod'и можуть спілкуватися з усіма Pod'ами (без обмежень)

[Frontend] ↔ [Backend] ✓
[Frontend] ↔ [Database] ✓ (ПОГАНО! Frontend не має прямого доступу до DB)
[Backend] ↔ [Database] ✓
```

**З NetworkPolicy:**
```yaml
# Політика для Database Pod'ів
spec:
  podSelector:
    matchLabels:
      app: postgres              # Застосувати до postgres Pod'ів

  ingress:                       # Вхідний трафік (хто може з'єднуватися)
  - from:
    - podSelector:
        matchLabels:
          app: backend           # Тільки backend може з'єднуватися
    ports:
    - port: 5432
```

**Результат:**
```
[Frontend] → [Database] ✗ (блоковано)
[Backend] → [Database] ✓ (дозволено)
```

**Важливо:**
- CNI plugin має підтримувати NetworkPolicy (Calico, Cilium)
- Flannel НЕ підтримує!
- Коли створюєш NetworkPolicy - стає whitelist (дозволено тільки вказане)

---

## 🔗 Як все це працює разом

### Повний flow запиту

```
1. Користувач: https://api.example.com/congratulation
           ↓
2. DNS: api.example.com → 34.123.45.67 (LoadBalancer IP)
           ↓
3. LoadBalancer → Ingress Controller (nginx Pod)
           ↓
4. Ingress Controller читає Ingress rules:
   host: api.example.com → backend-service
           ↓
5. Ingress Controller: HTTP запит до backend-service:80
           ↓
6. Service (backend-service):
   - selector: app=backend
   - Знаходить Pod'и: [Pod1, Pod2, Pod3]
   - Вибирає один (round-robin): Pod2
           ↓
7. Service перенаправляє на Pod2:8383
           ↓
8. Pod2 (backend контейнер):
   - Читає env змінні з ConfigMap/Secret
   - Обробляє запит
   - Повертає JSON: {"message": "Congratulation"}
           ↓
9. Відповідь повертається назад користувачу
           ↓
10. Користувач бачить: {"message": "Congratulation"}
```

### Lifecycle Pod'а

```
1. Ти пишеш: kubectl apply -f deployment.yaml
           ↓
2. Deployment створює ReplicaSet
           ↓
3. ReplicaSet створює Pod
           ↓
4. Pod статус: Pending (чекає scheduler)
           ↓
5. Scheduler вибирає Node (яка має достатньо CPU/RAM)
           ↓
6. Pod статус: Pending (чекає образ)
           ↓
7. Kubelet на Node'і:
   - docker pull ghcr.io/user/backend:v1.0.0
           ↓
8. Pod статус: ContainerCreating
           ↓
9. Docker запускає контейнер
           ↓
10. Pod статус: Running (але ще не Ready)
           ↓
11. Readiness probe перевіряє: GET /health
    - Failed... чекаємо...
    - Failed... чекаємо...
    - Success! ✓
           ↓
12. Pod статус: Running + Ready ✓
           ↓
13. Service додає Pod до Endpoints
           ↓
14. Тепер трафік йде на цей Pod!
```

### Що відбувається коли Pod падає

```
1. Pod crash (exit code != 0)
           ↓
2. Kubelet бачить що процес завершився
           ↓
3. restartPolicy: Always → перезапустити
           ↓
4. Pod статус: CrashLoopBackOff (падає і рестартує)
           ↓
5. Kubernetes exponential backoff:
   - Спроба 1: відразу
   - Спроба 2: через 10 секунд
   - Спроба 3: через 20 секунд
   - Спроба 4: через 40 секунд
   - ...до 5 хвилин
           ↓
6. Якщо livenessProbe теж failed:
   Kubernetes: "Щось не так з Pod'ом, створюю новий"
           ↓
7. ReplicaSet бачить що Pod dead
           ↓
8. ReplicaSet створює НОВИЙ Pod (нова IP, нове ім'я)
           ↓
9. Старий Pod видаляється
```

---

## 🎯 Практичні сценарії

### Сценарій 1: Оновлення додатку (Rolling Update)

```bash
# 1. Побудувати новий Docker образ
docker build -t ghcr.io/user/backend:v1.1.0 .
docker push ghcr.io/user/backend:v1.1.0

# 2. Оновити Deployment
kubectl set image deployment/backend \
  backend=ghcr.io/user/backend:v1.1.0 \
  -n orchestra-staging

# 3. Дивитися процес
kubectl rollout status deployment/backend -n orchestra-staging

# Що відбувається:
# [Pod-v1.0] [Pod-v1.0] [Pod-v1.0]  ← старі
#     ↓
# [Pod-v1.0] [Pod-v1.0] [Pod-v1.0] [Pod-v1.1] ← створюється новий
#     ↓
# [Pod-v1.0] [Pod-v1.0] [Pod-v1.1✓]  ← новий готовий, старий видалений
#     ↓
# [Pod-v1.0] [Pod-v1.1✓] [Pod-v1.1✓]
#     ↓
# [Pod-v1.1✓] [Pod-v1.1✓] [Pod-v1.1✓]  ← всі оновлені!
```

### Сценарій 2: Відкат (Rollback)

```bash
# Упс! Нова версія багнута, відкатуємо:
kubectl rollout undo deployment/backend -n orchestra-staging

# Дивитися історію
kubectl rollout history deployment/backend -n orchestra-staging

# Відкат до конкретної версії
kubectl rollout undo deployment/backend --to-revision=3 -n orchestra-staging
```

### Сценарій 3: Масштабування

```bash
# Вручну збільшити кількість Pod'ів
kubectl scale deployment/backend --replicas=5 -n orchestra-staging

# З HPA - автоматично:
# - Багато трафіку → HPA додає Pod'и
# - Мало трафіку → HPA прибирає зайві Pod'и
```

### Сценарій 4: Debug проблемного Pod'а

```bash
# 1. Подивитися статус
kubectl get pods -n orchestra-staging
# NAME                       READY   STATUS             RESTARTS
# backend-abc123-xyz         0/1     CrashLoopBackOff   5

# 2. Подивитися деталі
kubectl describe pod backend-abc123-xyz -n orchestra-staging
# Дивись Events внизу!

# 3. Логи
kubectl logs backend-abc123-xyz -n orchestra-staging
kubectl logs backend-abc123-xyz -n orchestra-staging --previous  # Попередній crash

# 4. Exec в Pod (якщо він працює)
kubectl exec -it backend-abc123-xyz -n orchestra-staging -- sh
ls /app
cat /app/config.json
ps aux
```

### Сценарій 5: Тестування локально

```bash
# Port forward для локального доступу
kubectl port-forward svc/backend-service 8080:80 -n orchestra-staging

# В іншому терміналі:
curl http://localhost:8080/congratulation
```

---

## 🔬 Що очікувати після deploy

### 1. Одразу після `kubectl apply`

```bash
kubectl get all -n orchestra-staging

# Очікуй побачити:
NAME                           READY   STATUS              RESTARTS
pod/backend-abc123-xyz         0/1     ContainerCreating   0    ← Завантажує образ
pod/backend-abc123-www         0/1     Pending             0    ← Чекає scheduler

NAME                      TYPE        CLUSTER-IP      EXTERNAL-IP
service/backend-service   ClusterIP   10.96.0.50      <none>

NAME                      READY   UP-TO-DATE   AVAILABLE
deployment.apps/backend   0/2     2            0    ← 0 з 2 готові

NAME                                 DESIRED   CURRENT   READY
replicaset.apps/backend-abc123       2         2         0
```

### 2. Через 30-60 секунд (після завантаження образу)

```bash
kubectl get pods -n orchestra-staging

NAME                       READY   STATUS    RESTARTS   AGE
backend-abc123-xyz         1/1     Running   0          45s  ← Готовий!
backend-abc123-www         1/1     Running   0          45s  ← Готовий!

# READY 1/1 означає: 1 з 1 контейнерів готовий
# Readiness probe пройшла успішно!
```

### 3. Перевірити чи Service знайшов Pod'и

```bash
kubectl get endpoints backend-service -n orchestra-staging

NAME              ENDPOINTS                         AGE
backend-service   10.244.1.5:8383,10.244.2.8:8383   1m

# Має бути IP адреси Pod'ів!
# Якщо пусто - selector не співпадає з labels
```

### 4. Перевірити HPA

```bash
kubectl get hpa -n orchestra-staging

NAME          REFERENCE            TARGETS   MINPODS   MAXPODS   REPLICAS
backend-hpa   Deployment/backend   15%/70%   2         10        2

# TARGETS 15%/70% означає:
# - Поточний CPU: 15%
# - Target CPU: 70%
# - Все ОК, не треба масштабувати
```

### 5. Тестувати доступ

```bash
# Через port-forward
kubectl port-forward svc/backend-service 8080:80 -n orchestra-staging
curl http://localhost:8080/congratulation

# Має повернути:
{
  "message": "Congratulation"
}
```

---

## ❗ Типові помилки та як їх виправити

### 1. ImagePullBackOff

```
STATUS: ImagePullBackOff

Причина: Не може завантажити Docker образ

Перевірити:
1. Чи правильне ім'я образу в deployment.yaml?
2. Чи образ існує в registry?
3. Чи є доступ до registry (приватний → треба imagePullSecret)?

Виправлення:
- Перевір: docker pull ghcr.io/user/backend:tag
- Якщо приватний: kubectl create secret docker-registry
```

### 2. CrashLoopBackOff

```
STATUS: CrashLoopBackOff
RESTARTS: 5

Причина: Контейнер запускається і одразу падає

Перевірити логи:
kubectl logs pod-name -n namespace
kubectl logs pod-name -n namespace --previous

Типові причини:
- Помилка в коді (exit 1)
- Неправильні env змінні
- Не може з'єднатися до DB
- Port вже зайнятий
```

### 3. Endpoints пусті

```
kubectl get endpoints service-name -n namespace
# <none>

Причина: Service не знаходить Pod'и

Перевірити:
# Labels Service
kubectl get svc service-name -n namespace -o yaml | grep selector

# Labels Pod'ів
kubectl get pods -n namespace --show-labels

# Мають співпадати!

Виправлення:
- Виправити selector в service.yaml або labels в deployment.yaml
```

### 4. Ingress повертає 503

```
curl https://api.example.com
# 503 Service Unavailable

Причини:
1. Service не існує
2. Endpoints пусті (Pod'и не готові)
3. Неправильне ім'я Service в Ingress

Перевірити:
kubectl get svc -n namespace
kubectl get endpoints -n namespace
kubectl describe ingress -n namespace
```

### 5. HPA показує `<unknown>`

```
kubectl get hpa
# TARGETS: <unknown>/70%

Причина: Metrics Server не працює

Виправлення:
# Встановити Metrics Server
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# Перевірити
kubectl top nodes
kubectl top pods -n namespace
```

---

## 📚 Корисні команди (шпаргалка)

```bash
# === DEPLOY ===
kubectl apply -f file.yaml                    # Створити/оновити з файлу
kubectl apply -k directory/                   # Kustomize
kubectl delete -f file.yaml                   # Видалити

# === STATUS ===
kubectl get pods -n namespace                 # Список Pod'ів
kubectl get all -n namespace                  # Все
kubectl describe pod pod-name -n namespace    # Деталі
kubectl get events -n namespace --sort-by='.lastTimestamp'  # Події

# === LOGS ===
kubectl logs pod-name -n namespace            # Логи
kubectl logs -f pod-name -n namespace         # Follow (live)
kubectl logs pod-name -n namespace --previous # Попередній контейнер

# === DEBUG ===
kubectl exec -it pod-name -n namespace -- sh  # Shell в Pod
kubectl port-forward svc/service 8080:80      # Port forward
kubectl top pods -n namespace                 # Ресурси

# === SCALE ===
kubectl scale deployment/name --replicas=5    # Масштабувати

# === UPDATE ===
kubectl set image deployment/name container=image:tag  # Оновити образ
kubectl rollout status deployment/name        # Статус
kubectl rollout undo deployment/name          # Відкат
kubectl rollout restart deployment/name       # Перезапуск

# === CONFIG ===
kubectl get configmap name -n namespace -o yaml      # ConfigMap
kubectl get secret name -n namespace -o yaml         # Secret (encoded)
kubectl get secret name -n namespace -o json | jq -r '.data.KEY | @base64d'  # Decode

# === CONTEXT ===
kubectl config get-contexts                   # Список contexts
kubectl config use-context context-name       # Перемкнути
kubectl config view                           # Конфігурація
```

---

## 🎓 Висновок

### Ти вивчив:

✅ **Що таке Kubernetes** і навіщо він потрібен
✅ **Основні об'єкти**: Pod, Deployment, Service, Ingress
✅ **Як працює кожен компонент** і як вони взаємодіють
✅ **Як писати YAML файли** і що означає кожне поле
✅ **Як deploy'ити додаток** в Kubernetes
✅ **Як debug'ати проблеми** і читати логи
✅ **Практичні сценарії**: оновлення, rollback, масштабування

### Наступні кроки:

1. **Практика**: Deploy свій backend в Minikube/Kind локально
2. **Моніторинг**: Prometheus + Grafana для метрик
3. **Логування**: ELK Stack або Loki для централізованих логів
4. **GitOps**: ArgoCD або Flux для автоматичного deploy з Git
5. **Service Mesh**: Istio або Linkerd для advanced networking
6. **Security**: Network Policies, Pod Security Standards, RBAC

---

**Успіхів у вивченні Kubernetes! 🚀**

Якщо щось незрозуміло - перечитай конкретну секцію, поекспериментуй з командами, подивись логи. Kubernetes здається складним спочатку, але після практики стає інтуїтивним!
