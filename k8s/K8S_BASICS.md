# Kubernetes YAML - Повний навчальний курс

## 📚 Зміст

1. [Основи YAML](#основи-yaml)
2. [Структура K8s YAML файлу](#структура-k8s-yaml-файлу)
3. [Основні K8s ресурси](#основні-k8s-ресурси)
4. [Як правильно описувати що ти хочеш](#як-правильно-описувати-що-ти-хочеш)
5. [Практичні приклади](#практичні-приклади)
6. [Best Practices](#best-practices)

---

## Основи YAML

### Що таке YAML?

**YAML** (YAML Ain't Markup Language) - це формат серіалізації даних, зручний для читання людиною.

### Базовий синтаксис:

```yaml
# Це коментар (починається з #)

# Ключ-значення (key: value)
name: orchestra-backend
port: 8383

# Списки (lists)
colors:
  - red
  - green
  - blue

# Або в одну строку
colors: [red, green, blue]

# Об'єкти (nested structures)
person:
  name: John
  age: 30
  address:
    city: Kyiv
    street: Main St

# Багаторядковий текст
description: |
  Це багаторядковий текст.
  Кожна строка зберігається як окрема.

# Текст в одну строку
description: >
  Це теж багаторядковий текст,
  але він буде з'єднаний в одну строку.
```

### ⚠️ Важливо!

1. **Відступи (indentation)** - тільки пробіли, НЕ табуляція!
2. **2 пробіли** - стандарт для K8s
3. **Регістр важливий** - `Name` ≠ `name`
4. **Двокрапка + пробіл** - завжди `key: value`, не `key:value`

---

## Структура K8s YAML файлу

### Обов'язкові поля

Кожен K8s YAML файл має 4 обов'язкових поля:

```yaml
# 1. apiVersion - версія Kubernetes API
apiVersion: v1

# 2. kind - тип ресурсу (Deployment, Service, Pod, etc.)
kind: Pod

# 3. metadata - метадані (ім'я, labels, annotations)
metadata:
  name: my-app

# 4. spec - специфікація (що саме ти хочеш)
spec:
  containers:
    - name: app
      image: nginx
```

### Детально про кожне поле:

#### 1. **apiVersion**

Визначає яку версію API використовувати:

```yaml
# Для базових ресурсів
apiVersion: v1              # Pod, Service, ConfigMap, Secret

# Для Deployments, StatefulSets
apiVersion: apps/v1         # Deployment, StatefulSet, DaemonSet

# Для Ingress
apiVersion: networking.k8s.io/v1  # Ingress

# Для Jobs
apiVersion: batch/v1        # Job, CronJob
```

#### 2. **kind**

Тип ресурсу який ти створюєш:

```yaml
kind: Pod           # Один контейнер/група контейнерів
kind: Deployment    # Керування Pod'ами (replicas, rolling updates)
kind: Service       # Мережевий доступ до Pod'ів
kind: Ingress       # HTTP/HTTPS маршрутизація
kind: ConfigMap     # Конфігурація
kind: Secret        # Секрети (паролі, ключі)
kind: Namespace     # Ізоляція ресурсів
```

#### 3. **metadata**

Метадані ресурсу:

```yaml
metadata:
  # Обов'язкове - ім'я ресурсу
  name: my-backend

  # Необов'язкове - namespace (за замовчуванням "default")
  namespace: production

  # Labels - для селекції та організації
  labels:
    app: backend
    version: v1
    environment: prod

  # Annotations - додаткова інформація
  annotations:
    description: "Backend API service"
    maintainer: "devops@example.com"
```

**Labels vs Annotations:**

- **Labels** - для СЕЛЕКЦІЇ ресурсів (Service знаходить Pod'и по labels)
- **Annotations** - для ІНФОРМАЦІЇ (не використовуються для селекції)

#### 4. **spec**

Специфікація - ЩО ТИ ХОЧЕШ:

```yaml
spec:
  # Залежить від kind!
  # Для Deployment:
  replicas: 3
  selector:
    matchLabels:
      app: backend
  template:
    metadata:
      labels:
        app: backend
    spec:
      containers:
      - name: backend
        image: backend:v1
        ports:
        - containerPort: 8080
```

---

## Основні K8s ресурси

### 1. Pod

**Що це:** Найменша одиниця в K8s. Один або кілька контейнерів.

**Коли використовувати:** Майже ніколи напряму! Використовуй Deployment.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: backend-pod
  labels:
    app: backend
spec:
  containers:
  - name: backend
    image: orchestra-backend:v1
    ports:
    - containerPort: 8383
```

**Що тут написано:**
- "Kubernetes, створи мені Pod"
- "Назви його backend-pod"
- "Дай йому label app=backend"
- "Запусти контейнер з образу orchestra-backend:v1"
- "Відкрий порт 8383"

---

### 2. Deployment

**Що це:** Керує Pod'ами - створює, оновлює, масштабує.

**Коли використовувати:** Завжди для stateless застосунків!

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend-deployment
  labels:
    app: backend
spec:
  # Скільки копій Pod'ів ти хочеш
  replicas: 3

  # Як знайти Pod'и якими керувати
  selector:
    matchLabels:
      app: backend

  # Шаблон для створення Pod'ів
  template:
    metadata:
      labels:
        app: backend
    spec:
      containers:
      - name: backend
        image: orchestra-backend:v1
        ports:
        - containerPort: 8383
        # Ресурси
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"
```

**Що тут написано:**
- "Kubernetes, створи мені Deployment"
- "Назви його backend-deployment"
- "Створи 3 копії (replicas) Pod'ів"
- "Знайди Pod'и з label app=backend"
- "Кожен Pod має запускати контейнер backend"
- "Дай кожному Pod'у 128Mi RAM (мінімум) і 256Mi (максимум)"
- "Дай кожному Pod'у 0.1 CPU (мінімум) і 0.2 CPU (максимум)"

**Ключові концепції:**
- **replicas** - кількість копій
- **selector** - як знайти Pod'и
- **template** - як створювати Pod'i
- **resources.requests** - мінімум ресурсів (для scheduling)
- **resources.limits** - максимум ресурсів (не може перевищити)

---

### 3. Service

**Що це:** Надає постійний IP та DNS для доступу до Pod'ів.

**Коли використовувати:** Коли треба дати доступ до Pod'ів (всередині або зовні кластера).

```yaml
apiVersion: v1
kind: Service
metadata:
  name: backend-service
spec:
  # Тип сервісу
  type: ClusterIP  # або NodePort, LoadBalancer

  # Які Pod'и обслуговувати
  selector:
    app: backend

  # Порти
  ports:
  - name: http
    protocol: TCP
    port: 80         # Порт сервісу
    targetPort: 8383 # Порт контейнера
```

**Що тут написано:**
- "Kubernetes, створи мені Service"
- "Знайди всі Pod'и з label app=backend"
- "Створи внутрішній IP (ClusterIP)"
- "Коли хтось звертається на порт 80 сервісу"
- "Перенаправ на порт 8383 одного з Pod'ів"

**Типи Service:**

1. **ClusterIP** (за замовчуванням)
   - Доступ ТІЛЬКИ всередині кластера
   - Використання: внутрішні сервіси (база даних, internal API)

2. **NodePort**
   - Доступ ззовні через IP ноди + порт (30000-32767)
   - Використання: development, testing

3. **LoadBalancer**
   - Створює зовнішній Load Balancer (потрібен cloud provider)
   - Використання: production (зовнішній доступ)

---

### 4. ConfigMap

**Що це:** Зберігає конфігурацію (не секретну).

**Коли використовувати:** Для налаштувань, які можуть змінюватись.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: backend-config
data:
  # Ключ-значення
  PORT: "8383"
  LOG_LEVEL: "info"

  # Або цілий файл
  app.conf: |
    port=8383
    log_level=info
    cors_enabled=true
```

**Як використати в Pod:**

```yaml
spec:
  containers:
  - name: backend
    image: backend:v1
    # Як змінні оточення
    env:
    - name: PORT
      valueFrom:
        configMapKeyRef:
          name: backend-config
          key: PORT
    # Або як файл
    volumeMounts:
    - name: config
      mountPath: /etc/config
  volumes:
  - name: config
    configMap:
      name: backend-config
```

---

### 5. Secret

**Що це:** Зберігає секретні дані (паролі, ключі).

**Коли використовувати:** Для чутливих даних.

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: backend-secrets
type: Opaque
data:
  # Base64 encoded!
  DB_PASSWORD: cGFzc3dvcmQxMjM=
  API_KEY: YWJjZGVmZ2hpams=
stringData:
  # Звичайний текст (K8s закодує сам)
  DB_HOST: postgres.default.svc.cluster.local
```

**Як використати:**

```yaml
spec:
  containers:
  - name: backend
    env:
    - name: DB_PASSWORD
      valueFrom:
        secretKeyRef:
          name: backend-secrets
          key: DB_PASSWORD
```

---

## Як правильно описувати що ти хочеш

### Метод "Від задачі до YAML"

#### Крок 1: Опиши що хочеш СЛОВАМИ

Приклад:
> "Я хочу запустити мій backend на Go.
> Мені треба 3 копії для високої доступності.
> Кожна копія має використовувати не більше 256MB RAM.
> Треба щоб інші сервіси могли звертатися до нього на порту 80."

#### Крок 2: Розбий на компоненти

1. **Запустити backend** → Deployment
2. **3 копії** → `replicas: 3`
3. **256MB RAM** → `resources.limits.memory`
4. **Доступ на порту 80** → Service

#### Крок 3: Напиши YAML

```yaml
# 1. Deployment для backend
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend
spec:
  replicas: 3  # 3 копії
  selector:
    matchLabels:
      app: backend
  template:
    metadata:
      labels:
        app: backend
    spec:
      containers:
      - name: backend
        image: backend:v1
        resources:
          limits:
            memory: "256Mi"  # Не більше 256MB
        ports:
        - containerPort: 8383

---
# 2. Service для доступу
apiVersion: v1
kind: Service
metadata:
  name: backend
spec:
  selector:
    app: backend
  ports:
  - port: 80
    targetPort: 8383
```

### Шаблон питань для себе:

1. **Що я запускаю?** → Deployment (або StatefulSet для stateful)
2. **Скільки копій?** → `replicas: N`
3. **Які ресурси потрібні?** → `resources`
4. **Які порти?** → `containerPort` в Pod, `port` в Service
5. **Треба доступ ззовні?** → Service type: LoadBalancer або Ingress
6. **Є конфігурація?** → ConfigMap
7. **Є секрети?** → Secret
8. **Треба постійне сховище?** → PersistentVolume + PersistentVolumeClaim

---

## Практичні приклади

### Приклад 1: Простий Pod

**Задача:** Запусти nginx

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - name: nginx
    image: nginx:latest
    ports:
    - containerPort: 80
```

### Приклад 2: Production Deployment

**Задача:** Запусти backend з 3 репліками, health checks, ресурсами

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: backend
  labels:
    app: backend
    version: v1
spec:
  replicas: 3
  selector:
    matchLabels:
      app: backend
  template:
    metadata:
      labels:
        app: backend
        version: v1
    spec:
      containers:
      - name: backend
        image: backend:v1
        ports:
        - containerPort: 8383

        # Health checks
        livenessProbe:
          httpGet:
            path: /health
            port: 8383
          initialDelaySeconds: 30
          periodSeconds: 10

        readinessProbe:
          httpGet:
            path: /ready
            port: 8383
          initialDelaySeconds: 5
          periodSeconds: 5

        # Ресурси
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "200m"

        # Environment
        env:
        - name: PORT
          value: "8383"
        - name: LOG_LEVEL
          value: "info"
```

**Пояснення нових полів:**

- **livenessProbe** - "Чи живий Pod?" Якщо ні → перезапуск
- **readinessProbe** - "Чи готовий приймати трафік?" Якщо ні → не слати трафік
- **initialDelaySeconds** - скільки чекати після старту
- **periodSeconds** - як часто перевіряти

---

## Best Practices

### 1. Використовуй Labels розумно

```yaml
# ✓ Добре - структуровані labels
labels:
  app: backend
  component: api
  version: v1
  environment: production
  managed-by: helm

# ✗ Погано - мало інформації
labels:
  name: backend
```

### 2. Завжди вказуй resources

```yaml
# ✓ Добре
resources:
  requests:
    memory: "128Mi"
    cpu: "100m"
  limits:
    memory: "256Mi"
    cpu: "200m"

# ✗ Погано - немає лімітів
# Може з'їсти всі ресурси ноди!
```

### 3. Використовуй Health Checks

```yaml
# ✓ Добре - є перевірки
livenessProbe:
  httpGet:
    path: /health
    port: 8080
readinessProbe:
  httpGet:
    path: /ready
    port: 8080

# ✗ Погано - без перевірок
# K8s не знає чи працює застосунок
```

### 4. Ніколи не використовуй :latest

```yaml
# ✓ Добре - конкретна версія
image: backend:v1.2.3

# ✗ Погано - непередбачувано
image: backend:latest
```

### 5. Використовуй окремі файли

```
k8s/
├── backend/
│   ├── deployment.yaml
│   ├── service.yaml
│   └── configmap.yaml
└── frontend/
    ├── deployment.yaml
    ├── service.yaml
    └── ingress.yaml
```

### 6. Namespace для ізоляції

```yaml
# Розділяй environments
namespace: production
namespace: staging
namespace: development
```

### 7. ConfigMap для конфігурації

```yaml
# ✓ Добре - конфігурація окремо
env:
- name: PORT
  valueFrom:
    configMapKeyRef:
      name: config
      key: PORT

# ✗ Погано - хардкод
env:
- name: PORT
  value: "8080"
```

---

## Швидкий довідник

### Обов'язкові поля:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: my-app
spec:
  # тут специфікація
```

### Часті комбінації:

**Stateless app:**
- Deployment (для Pod'ів)
- Service (для доступу)
- ConfigMap (для конфігурації)
- Secret (для паролів)

**Web app з зовнішнім доступом:**
- Deployment
- Service (ClusterIP)
- Ingress (HTTP/HTTPS routing)

**Stateful app (база даних):**
- StatefulSet (замість Deployment)
- Service (Headless)
- PersistentVolumeClaim
- ConfigMap + Secret

---

## Наступні кроки

1. ✅ Вивчив основи YAML
2. ✅ Розумію структуру K8s файлів
3. ✅ Знаю основні ресурси
4. ➡️ Практика - створи файли для свого проекту
5. ➡️ Вивчи advanced теми (Ingress, StatefulSet, etc.)

Переходь до реальних прикладів для Orchestra проекту! →
