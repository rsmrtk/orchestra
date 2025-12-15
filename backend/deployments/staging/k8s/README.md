# 📦 Kubernetes Deployment для Orchestra Backend (Staging)

Цей каталог містить всі Kubernetes YAML файли для деплоймента Orchestra Backend в staging середовище.

## 📋 Зміст

- [Структура файлів](#-структура-файлів)
- [Передумови](#-передумови)
- [Швидкий старт](#-швидкий-старт)
- [Детальний гайд](#-детальний-гайд)
- [Kubernetes концепції](#-kubernetes-концепції)
- [Troubleshooting](#-troubleshooting)
- [Best Practices](#-best-practices)

---

## 📁 Структура файлів

```
k8s/
├── namespace.yaml          # Namespace для ізоляції ресурсів
├── configmap.yaml          # Конфігурація додатку (env змінні)
├── secret.yaml             # Секрети (паролі, токени)
├── deployment.yaml         # Управління Pod'ами
├── service.yaml            # Мережевий доступ до Pod'ів
├── ingress.yaml            # HTTP/HTTPS маршрутизація ззовні
├── hpa.yaml                # Автоматичне масштабування
├── pdb.yaml                # Бюджет збоїв Pod'ів
├── network-policy.yaml     # Правила мережевої безпеки
├── kustomization.yaml      # Kustomize конфігурація
└── README.md               # Цей файл
```

---

## ✅ Передумови

### 1. Kubernetes кластер
Тобі потрібен робочий Kubernetes кластер. Варіанти:

**Локально (для тестування):**
```bash
# Minikube
minikube start --cpus=4 --memory=8192
minikube addons enable ingress
minikube addons enable metrics-server

# Kind
kind create cluster --config=kind-config.yaml

# Docker Desktop
# Включити Kubernetes в настройках Docker Desktop
```

**Cloud (для production):**
- GKE (Google Kubernetes Engine)
- EKS (Amazon Elastic Kubernetes Service)
- AKS (Azure Kubernetes Service)
- DigitalOcean Kubernetes
- Linode Kubernetes Engine

### 2. Kubectl
```bash
# Перевірити версію
kubectl version --client

# Встановити (якщо нема):
# Linux
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x kubectl
sudo mv kubectl /usr/local/bin/

# macOS
brew install kubectl

# Windows
choco install kubernetes-cli
```

### 3. Docker образ
Треба побудувати Docker образ твого backend:

```bash
cd /home/r/fff/deploy/orchestra/backend

# Build образу
docker build -f deployments/staging/Dockerfile -t your-registry/orchestra-backend:staging .

# Push до registry (якщо не локальний кластер)
docker push your-registry/orchestra-backend:staging
```

**Приклади registry:**
- Docker Hub: `docker.io/username/orchestra-backend:staging`
- GitHub Container Registry: `ghcr.io/username/orchestra-backend:staging`
- Google Container Registry: `gcr.io/project-id/orchestra-backend:staging`

### 4. Налаштувати kubectl context
```bash
# Подивитися доступні contexts
kubectl config get-contexts

# Перемкнутися на потрібний context
kubectl config use-context <context-name>

# Перевірити з'єднання
kubectl cluster-info
kubectl get nodes
```

---

## 🚀 Швидкий старт

### Крок 1: Оновити Docker образ

Відкрий `deployment.yaml` і змінинапиши мені для /home/r/fff/deploy/orchestra/backend/deployments/staging/k8s yaml файли, та навчи мене їх писати, розуміти як їх робити і шо я маю очікувати, які поля і навіщо їх писати, тобто розпиши все по максимуму стосовно кубернетеса:

```yaml
spec:
  template:
    spec:
      containers:
      - name: backend
        image: your-registry/orchestra-backend:staging  # ⬅️ ЗМІНИ ЦЕ
```

### Крок 2: Оновити секрети

Відкрий `secret.yaml` та оновинапиші реальні секрети (закодовані в base64):

```bash
# Закодувати секрет
echo -n "твій-секретний-ключ" | base64

# Результат вставити в secret.yaml
```

**ВАЖЛИВО:** НЕ комітити реальні секрети в Git!

### Крок 3: Deploy!

```bash
# Метод 1: Використати kustomize (рекомендовано)
kubectl apply -k /home/r/fff/deploy/orchestra/backend/deployments/staging/k8s

# Метод 2: Apply окремі файли
kubectl apply -f namespace.yaml
kubectl apply -f configmap.yaml
kubectl apply -f secret.yaml
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml
kubectl apply -f ingress.yaml
kubectl apply -f hpa.yaml
kubectl apply -f pdb.yaml
kubectl apply -f network-policy.yaml
```

### Крок 4: Перевірити

```bash
# Перевірити чи все створилося
kubectl get all -n orchestra-staging

# Дивитися статус Pod'ів
kubectl get pods -n orchestra-staging -w

# Дивитися логи
kubectl logs -f deployment/orchestra-backend -n orchestra-staging
```

### Крок 5: Тестувати

```bash
# Port forward для локального доступу
kubectl port-forward -n orchestra-staging service/orchestra-backend-service 8080:80

# Тестувати API
curl http://localhost:8080/congratulation
```

---

## 📚 Детальний гайд

### 1. Namespace (namespace.yaml)

**Що це?**
Namespace - це віртуальний кластер всередині фізичного кластера. Дозволяє ізолювати ресурси різних проектів.

**Навіщо?**
- Ізоляція staging від production
- Окремий доступ та квоти
- Організація ресурсів

**Команди:**
```bash
# Створити
kubectl apply -f namespace.yaml

# Подивитися
kubectl get namespaces
kubectl describe namespace orchestra-staging

# Видалити (УВАГА: видалить ВСІ ресурси всередині!)
kubectl delete namespace orchestra-staging
```

---

### 2. ConfigMap (configmap.yaml)

**Що це?**
ConfigMap зберігає конфігураційні дані (env змінні, файли конфігурації).

**Навіщо?**
- Відокремити конфігурацію від коду
- Змінювати конфігурацію без пересборки образу
- Різна конфігурація для staging/production

**Як оновити:**
```bash
# Редагувати файл configmap.yaml, потім:
kubectl apply -f configmap.yaml

# Перезапустити Pod'и щоб підхопили нову конфігурацію
kubectl rollout restart deployment/orchestra-backend -n orchestra-staging
```

**Подивитися значення:**
```bash
kubectl get configmap orchestra-backend-config -n orchestra-staging -o yaml
```

---

### 3. Secret (secret.yaml)

**Що це?**
Secret зберігає чутливі дані (паролі, API ключі, токени).

**ВАЖЛИВО:**
- Дані в base64 (НЕ шифрування!)
- НЕ комітити реальні секрети в Git
- Використовуй Secret management (Vault, Sealed Secrets, External Secrets)

**Як створити секрет:**
```bash
# Метод 1: З командної строки (рекомендовано)
kubectl create secret generic orchestra-backend-secret \
  --from-literal=API_KEY=твій-ключ \
  --from-literal=DB_PASSWORD=твій-пароль \
  -n orchestra-staging

# Метод 2: З файлу
kubectl create secret generic orchestra-backend-secret \
  --from-file=tls.crt=path/to/cert.crt \
  --from-file=tls.key=path/to/cert.key \
  -n orchestra-staging

# Метод 3: З YAML файлу (треба base64)
echo -n "твій-секрет" | base64  # Закодувати
kubectl apply -f secret.yaml
```

**Подивитися секрет:**
```bash
# Подивитися (закодоване)
kubectl get secret orchestra-backend-secret -n orchestra-staging -o yaml

# Декодувати
kubectl get secret orchestra-backend-secret -n orchestra-staging \
  -o jsonpath='{.data.API_KEY}' | base64 -d
```

---

### 4. Deployment (deployment.yaml)

**Що це?**
Deployment управляє Pod'ами: створює, перезапускає, оновлює, масштабує.

**Ключові поля:**

```yaml
spec:
  replicas: 2                      # Кількість Pod'ів

  strategy:
    type: RollingUpdate            # Як оновлювати (RollingUpdate або Recreate)
    rollingUpdate:
      maxUnavailable: 1            # Скільки можуть бути недоступні під час оновлення
      maxSurge: 1                  # Скільки додаткових можна створити

  template:
    spec:
      containers:
      - name: backend
        image: your-registry/orchestra-backend:staging
        imagePullPolicy: Always     # Always, IfNotPresent, Never

        ports:
        - containerPort: 8383       # Порт контейнера

        resources:
          requests:                 # Мінімум (для scheduler)
            memory: "128Mi"
            cpu: "100m"
          limits:                   # Максимум (для захисту)
            memory: "256Mi"
            cpu: "200m"

        livenessProbe:              # Чи живий контейнер
          httpGet:
            path: /congratulation
            port: 8383

        readinessProbe:             # Чи готовий приймати трафік
          httpGet:
            path: /congratulation
            port: 8383
```

**Команди:**
```bash
# Створити
kubectl apply -f deployment.yaml

# Статус
kubectl get deployments -n orchestra-staging
kubectl describe deployment orchestra-backend -n orchestra-staging

# Pod'и
kubectl get pods -n orchestra-staging
kubectl describe pod <pod-name> -n orchestra-staging

# Логи
kubectl logs -f deployment/orchestra-backend -n orchestra-staging
kubectl logs <pod-name> -n orchestra-staging

# Масштабувати
kubectl scale deployment orchestra-backend --replicas=5 -n orchestra-staging

# Оновити образ
kubectl set image deployment/orchestra-backend \
  backend=your-registry/orchestra-backend:new-tag \
  -n orchestra-staging

# Статус оновлення
kubectl rollout status deployment/orchestra-backend -n orchestra-staging

# Відкат
kubectl rollout undo deployment/orchestra-backend -n orchestra-staging

# Історія
kubectl rollout history deployment/orchestra-backend -n orchestra-staging

# Перезапустити
kubectl rollout restart deployment/orchestra-backend -n orchestra-staging

# Exec в Pod
kubectl exec -it <pod-name> -n orchestra-staging -- sh
```

---

### 5. Service (service.yaml)

**Що це?**
Service - стабільна точка доступу до Pod'ів. Pod'и мають динамічні IP, Service - статичний.

**Типи:**
- **ClusterIP** (default): доступний тільки всередині кластера
- **NodePort**: доступний на кожній Node:port
- **LoadBalancer**: створює зовнішній Load Balancer (cloud)

**DNS імена:**
```
<service-name>.<namespace>.svc.cluster.local

Приклади:
- orchestra-backend-service                                    (той самий namespace)
- orchestra-backend-service.orchestra-staging                  (з іншого namespace)
- orchestra-backend-service.orchestra-staging.svc.cluster.local (повне FQDN)
```

**Команди:**
```bash
# Створити
kubectl apply -f service.yaml

# Статус
kubectl get services -n orchestra-staging
kubectl describe service orchestra-backend-service -n orchestra-staging

# Endpoints (IP Pod'ів)
kubectl get endpoints orchestra-backend-service -n orchestra-staging

# Port forward
kubectl port-forward svc/orchestra-backend-service 8080:80 -n orchestra-staging
curl http://localhost:8080/congratulation

# Тестувати з test Pod'а
kubectl run test --image=alpine -n orchestra-staging --rm -it -- sh
apk add curl
curl http://orchestra-backend-service:80/congratulation
```

---

### 6. Ingress (ingress.yaml)

**Що це?**
Ingress - HTTP/HTTPS маршрутизація ззовні кластера. Працює як reverse proxy.

**Передумови:**
Потрібен Ingress Controller (nginx, traefik, istio).

```bash
# Встановити nginx-ingress
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/cloud/deploy.yaml

# Перевірити
kubectl get pods -n ingress-nginx
kubectl get svc -n ingress-nginx
# Дивись EXTERNAL-IP - це адреса для DNS
```

**Налаштувати DNS:**
```
staging-api.orchestra.example.com -> <EXTERNAL-IP ingress controller>
```

**SSL/TLS (опційно):**
```bash
# Встановити cert-manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.13.0/cert-manager.yaml

# Створити ClusterIssuer (вже є в ingress.yaml)
kubectl apply -f ingress.yaml
```

**Команди:**
```bash
# Створити
kubectl apply -f ingress.yaml

# Статус
kubectl get ingress -n orchestra-staging
kubectl describe ingress orchestra-backend-ingress -n orchestra-staging

# SSL сертифікат (якщо є cert-manager)
kubectl get certificate -n orchestra-staging

# Тестувати
curl https://staging-api.orchestra.example.com/api/congratulation
```

---

### 7. HPA - HorizontalPodAutoscaler (hpa.yaml)

**Що це?**
HPA автоматично масштабує кількість Pod'ів базуючись на метриках (CPU, Memory).

**Передумови:**
Metrics Server має бути встановлений.

```bash
# Встановити Metrics Server
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# Перевірити
kubectl get deployment metrics-server -n kube-system
kubectl top nodes
kubectl top pods -n orchestra-staging
```

**Ключові поля:**
```yaml
spec:
  minReplicas: 2                   # Мінімум Pod'ів
  maxReplicas: 10                  # Максимум Pod'ів

  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70     # Масштабувати коли CPU > 70%
```

**Команди:**
```bash
# Створити
kubectl apply -f hpa.yaml

# Статус
kubectl get hpa -n orchestra-staging
kubectl describe hpa orchestra-backend-hpa -n orchestra-staging

# Дивитися в реальному часі
kubectl get hpa -n orchestra-staging --watch

# Load testing (тестувати масштабування)
# Apache Bench
ab -n 100000 -c 100 http://your-service/endpoint

# hey
hey -z 5m -c 50 http://your-service/endpoint
```

---

### 8. PDB - PodDisruptionBudget (pdb.yaml)

**Що це?**
PDB гарантує що мінімальна кількість Pod'ів залишається доступною під час оновлень/maintenance.

**Навіщо?**
- Висока доступність під час drain Node
- Контроль над оновленнями кластера
- Запобігти downtime

**Ключові поля:**
```yaml
spec:
  minAvailable: 1                  # Мінімум 1 Pod завжди доступний
  # АБО
  maxUnavailable: 1                # Максимум 1 Pod може бути недоступним
```

**Команди:**
```bash
# Створити
kubectl apply -f pdb.yaml

# Статус
kubectl get pdb -n orchestra-staging
kubectl describe pdb orchestra-backend-pdb -n orchestra-staging

# Подивитися disruptionsAllowed
kubectl get pdb orchestra-backend-pdb -n orchestra-staging -o jsonpath='{.status}'

# Тестувати (drain Node)
kubectl drain <node-name> --ignore-daemonsets
# PDB блокуватиме drain якщо це порушить minAvailable
```

---

### 9. NetworkPolicy (network-policy.yaml)

**Що це?**
NetworkPolicy - firewall rules для Pod'ів. Контролює вхідний/вихідний трафік.

**Передумови:**
CNI plugin має підтримувати NetworkPolicy (Calico, Cilium, Weave). Flannel НЕ підтримує.

**Ключові концепції:**
```yaml
spec:
  podSelector:                     # До яких Pod'ів застосовується
    matchLabels:
      app: orchestra-backend

  policyTypes:
  - Ingress                        # Вхідний трафік
  - Egress                         # Вихідний трафік

  ingress:                         # Хто може з'єднуватися ДО наших Pod'ів
  - from:
    - podSelector:
        matchLabels:
          app: frontend
    ports:
    - port: 8383

  egress:                          # Куди наші Pod'и можуть з'єднуватися
  - to:
    - podSelector:
        matchLabels:
          app: postgres
    ports:
    - port: 5432
```

**Команди:**
```bash
# Створити
kubectl apply -f network-policy.yaml

# Статус
kubectl get networkpolicy -n orchestra-staging
kubectl describe networkpolicy orchestra-backend-netpol -n orchestra-staging

# Тестувати
kubectl run test --image=alpine -n orchestra-staging --rm -it -- sh
apk add curl
curl http://orchestra-backend-service:80
# Якщо NetworkPolicy блокує - з'єднання не вийде
```

---

## 🧠 Kubernetes концепції

### Pod
- Найменша одиниця деплойменту
- Один або більше контейнерів
- Спільний namespace (network, storage)
- Ephemeral (тимчасовий) - може бути перезапущений будь-коли

### ReplicaSet
- Гарантує що певна кількість Pod'ів працює
- Автоматично створює нові Pod'и якщо старі падають
- Зазвичай не створюється вручну (Deployment створює його)

### Deployment
- Управляє ReplicaSet
- Rolling updates, rollbacks
- Декларативний стан
- Найчастіший спосіб deploy додатків

### Service
- Стабільна мережева точка доступу до Pod'ів
- Load balancing між Pod'ами
- Service Discovery через DNS

### Ingress
- HTTP/HTTPS маршрутизація
- SSL termination
- Virtual hosting (по domain/path)

### ConfigMap & Secret
- ConfigMap: конфігурація (не секретна)
- Secret: секретна інформація
- Можуть бути env змінними або файлами

### Namespace
- Віртуальний кластер
- Ізоляція ресурсів
- RBAC, квоти

### Labels & Selectors
- Labels: key-value мітки на ресурсах
- Selectors: як знайти ресурси по labels
- Критично для зв'язку Deployment -> Service -> Pod

---

## 🔧 Troubleshooting

### Pod не запускається

```bash
# Подивитися статус
kubectl get pods -n orchestra-staging

# Можливі статуси:
# - Pending: чекає на scheduler (можливо недостатньо ресурсів)
# - ContainerCreating: завантажує образ
# - Running: працює
# - CrashLoopBackOff: падає і перезапускається
# - ImagePullBackOff: не може завантажити образ
# - Error: помилка

# Деталі
kubectl describe pod <pod-name> -n orchestra-staging
# Дивись Events внизу!

# Логи
kubectl logs <pod-name> -n orchestra-staging
kubectl logs <pod-name> -n orchestra-staging --previous  # Логи попереднього контейнера

# Exec в Pod
kubectl exec -it <pod-name> -n orchestra-staging -- sh
```

**Типові проблеми:**

1. **ImagePullBackOff**
   - Неправильне ім'я образу
   - Образ не існує
   - Немає доступу до registry
   - Треба imagePullSecrets

2. **CrashLoopBackOff**
   - Додаток падає при старті
   - Дивись логи: `kubectl logs`
   - Перевір env змінні, config

3. **Pending**
   - Недостатньо ресурсів (CPU/Memory) на Node'ах
   - `kubectl describe node` - подивись Allocatable
   - Зменш resources.requests або додай Node

4. **OOMKilled (Out of Memory)**
   - Pod використав більше пам'яті ніж limits
   - Збільш resources.limits.memory

---

### Service не працює

```bash
# Перевірити Service
kubectl get svc -n orchestra-staging
kubectl describe svc orchestra-backend-service -n orchestra-staging

# Перевірити Endpoints (IP Pod'ів)
kubectl get endpoints orchestra-backend-service -n orchestra-staging

# Якщо Endpoints пусті:
# - Selector не співпадає з labels Pod'ів
# - Pod'и не готові (readinessProbe failed)
```

**Перевірити labels:**
```bash
# Labels Service selector
kubectl get svc orchestra-backend-service -n orchestra-staging -o jsonpath='{.spec.selector}'

# Labels Pod'ів
kubectl get pods -n orchestra-staging --show-labels

# Має співпадати!
```

---

### Ingress не працює

```bash
# Перевірити Ingress
kubectl get ingress -n orchestra-staging
kubectl describe ingress orchestra-backend-ingress -n orchestra-staging

# Перевірити Ingress Controller
kubectl get pods -n ingress-nginx
kubectl get svc -n ingress-nginx
# Дивись EXTERNAL-IP

# Перевірити DNS
nslookup staging-api.orchestra.example.com
# Має вказувати на EXTERNAL-IP

# Логи Ingress Controller
kubectl logs -n ingress-nginx deployment/ingress-nginx-controller
```

---

### HPA не масштабує

```bash
# Перевірити HPA
kubectl get hpa -n orchestra-staging
kubectl describe hpa orchestra-backend-hpa -n orchestra-staging

# Якщо метрики <unknown>:
# - Metrics Server не встановлений або не працює
kubectl get deployment metrics-server -n kube-system
kubectl top pods -n orchestra-staging

# Якщо не працює - встанови Metrics Server

# Перевірити чи вказані resources.requests в Deployment
# HPA потребує requests для обчислення %
```

---

## ✨ Best Practices

### 1. Resources
- Завжди вказуй `resources.requests` та `resources.limits`
- requests = мінімум (для scheduler)
- limits = максимум (для захисту)

### 2. Health Checks
- Завжди додавай `livenessProbe` та `readinessProbe`
- liveness = чи живий (для перезапуску)
- readiness = чи готовий (для трафіку)

### 3. Labels
- Використовуй стандартні labels: `app`, `environment`, `version`
- Селектори мають співпадати з labels Pod'ів

### 4. Secrets
- НЕ комітити секрети в Git
- Використовуй Secret management (Vault, Sealed Secrets)
- Або створюй через `kubectl create secret`

### 5. High Availability
- `replicas >= 2` для production
- Додавай PodDisruptionBudget
- Використовуй Pod Anti-Affinity (розміщувати на різних Node'ах)

### 6. Security
- Використовуй NetworkPolicy
- `runAsNonRoot: true`
- `readOnlyRootFilesystem: true` (якщо можливо)
- Сканування образів на вразливості

### 7. Моніторинг
- Логи: централізоване логування (ELK, Loki)
- Метрики: Prometheus + Grafana
- Tracing: Jaeger, Zipkin
- Alerts: на важливі метрики

### 8. CI/CD
- Автоматичне тестування
- Автоматичний build образів
- Tag образи (не `:latest` в production!)
- Canary/Blue-Green deployments

### 9. Gitops
- Всі YAML в Git
- Автоматичний deploy через GitOps (ArgoCD, Flux)
- Review changes через Pull Requests

### 10. Namespace Organization
- Окремі namespace для staging/production
- ResourceQuotas для namespace
- LimitRanges для default requests/limits

---

## 📖 Корисні команди

### Загальні

```bash
# Всі ресурси в namespace
kubectl get all -n orchestra-staging

# Конкретні типи
kubectl get pods,svc,deploy,ingress -n orchestra-staging

# Широкий вивід (більше інфо)
kubectl get pods -n orchestra-staging -o wide

# YAML вивід
kubectl get deployment orchestra-backend -n orchestra-staging -o yaml

# JSON Path
kubectl get pods -n orchestra-staging -o jsonpath='{.items[*].metadata.name}'

# Watch (auto-refresh)
kubectl get pods -n orchestra-staging --watch

# Labels
kubectl get pods -n orchestra-staging --show-labels
kubectl get pods -n orchestra-staging -l app=orchestra-backend
```

### Debugging

```bash
# Describe (деталі + Events)
kubectl describe pod <pod-name> -n orchestra-staging

# Логи
kubectl logs <pod-name> -n orchestra-staging
kubectl logs <pod-name> -n orchestra-staging -f  # Follow
kubectl logs <pod-name> -n orchestra-staging --previous  # Попередній контейнер
kubectl logs <pod-name> -n orchestra-staging -c <container-name>  # Multi-container Pod

# Exec
kubectl exec <pod-name> -n orchestra-staging -- ls /app
kubectl exec -it <pod-name> -n orchestra-staging -- sh

# Port forward
kubectl port-forward <pod-name> 8080:8383 -n orchestra-staging
kubectl port-forward svc/orchestra-backend-service 8080:80 -n orchestra-staging

# Copy files
kubectl cp <pod-name>:/path/to/file ./local-file -n orchestra-staging
kubectl cp ./local-file <pod-name>:/path/to/file -n orchestra-staging

# Top (ресурси)
kubectl top nodes
kubectl top pods -n orchestra-staging
kubectl top pods -n orchestra-staging --sort-by=memory
```

### Edit

```bash
# Edit в редакторі
kubectl edit deployment orchestra-backend -n orchestra-staging

# Patch (JSON)
kubectl patch deployment orchestra-backend -n orchestra-staging -p '{"spec":{"replicas":5}}'

# Scale
kubectl scale deployment orchestra-backend --replicas=3 -n orchestra-staging

# Set image
kubectl set image deployment/orchestra-backend backend=new-image:tag -n orchestra-staging
```

### Rollout

```bash
# Статус deploy
kubectl rollout status deployment/orchestra-backend -n orchestra-staging

# Історія
kubectl rollout history deployment/orchestra-backend -n orchestra-staging

# Undo (rollback)
kubectl rollout undo deployment/orchestra-backend -n orchestra-staging
kubectl rollout undo deployment/orchestra-backend --to-revision=2 -n orchestra-staging

# Restart
kubectl rollout restart deployment/orchestra-backend -n orchestra-staging

# Pause/Resume (для batch змін)
kubectl rollout pause deployment/orchestra-backend -n orchestra-staging
kubectl set image deployment/orchestra-backend backend=new-image:tag -n orchestra-staging
kubectl set resources deployment/orchestra-backend -c=backend --limits=memory=512Mi -n orchestra-staging
kubectl rollout resume deployment/orchestra-backend -n orchestra-staging
```

### Cleanup

```bash
# Видалити конкретний ресурс
kubectl delete pod <pod-name> -n orchestra-staging
kubectl delete deployment orchestra-backend -n orchestra-staging

# Видалити з файлу
kubectl delete -f deployment.yaml

# Видалити з kustomize
kubectl delete -k /home/r/fff/deploy/orchestra/backend/deployments/staging/k8s

# Видалити namespace (видалить ВСІ ресурси!)
kubectl delete namespace orchestra-staging

# Force delete (якщо застряг)
kubectl delete pod <pod-name> -n orchestra-staging --grace-period=0 --force
```

---

## 📚 Додаткові ресурси

### Документація
- [Kubernetes Official Docs](https://kubernetes.io/docs/)
- [kubectl Cheat Sheet](https://kubernetes.io/docs/reference/kubectl/cheatsheet/)
- [Kustomize Docs](https://kubectl.docs.kubernetes.io/guides/)

### Інтерактивні туторіали
- [Kubernetes Tutorial](https://kubernetes.io/docs/tutorials/)
- [Play with Kubernetes](https://labs.play-with-k8s.com/)
- [Katacoda Kubernetes](https://www.katacoda.com/courses/kubernetes)

### Книги (українською/російською)
- "Kubernetes в действии" (Марко Лукша)
- "Kubernetes Patterns" (Bilgin Ibryam, Roland Huß)

### YouTube канали
- [TechWorld with Nana](https://www.youtube.com/c/TechWorldwithNana)
- [Just me and Opensource](https://www.youtube.com/c/wenkatn-justmeandopensource)

### Інструменти
- [k9s](https://k9scli.io/) - TUI для Kubernetes
- [kubectx/kubens](https://github.com/ahmetb/kubectx) - Швидке перемикання context/namespace
- [stern](https://github.com/stern/stern) - Multi-pod логи
- [Lens](https://k8slens.dev/) - Kubernetes IDE

---

## 🤝 Допомога

Якщо щось не працює:

1. Прочитай секцію [Troubleshooting](#-troubleshooting)
2. Використай `kubectl describe` та `kubectl logs`
3. Перевір Events: `kubectl get events -n orchestra-staging --sort-by='.lastTimestamp'`
4. Погугли помилку + "kubernetes"
5. Запитай в [Kubernetes Slack](https://kubernetes.slack.com/)

---

## 📝 License

Цей проект для внутрішнього використання Orchestra.

---

**Зроблено з ❤️ для Orchestra Backend Team**
