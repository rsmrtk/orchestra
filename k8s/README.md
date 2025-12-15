# Kubernetes Configuration для Orchestra

## 📚 Документація

### Навчальні матеріали:
- **[K8S_BASICS.md](K8S_BASICS.md)** ⭐ - Повний навчальний курс по Kubernetes YAML
  - Основи YAML синтаксису
  - Структура K8s файлів
  - Основні ресурси (Pod, Deployment, Service, etc.)
  - Як правильно описувати що ти хочеш
  - Best practices

### YAML файли проекту:
- **[backend/](backend/)** - Backend конфігурація
  - `deployment.yaml` - Backend Deployment з детальними коментарями
  - `service.yaml` - Backend Service
  - `configmap.yaml` - Конфігурація
- **[frontend/](frontend/)** - Frontend конфігурація
  - `deployment.yaml` - Frontend Deployment
  - `service.yaml` - Frontend Service
  - `ingress.yaml` - HTTP/HTTPS маршрутизація
- **[common/](common/)** - Спільні ресурси
  - `namespace.yaml` - Namespaces для різних environments

## 🚀 Швидкий старт

### Крок 1: Встановити Kubernetes

Вибери один з варіантів:

**Локально (для навчання):**
- **minikube**: https://minikube.sigs.k8s.io/docs/start/
  ```bash
  # MacOS
  brew install minikube
  minikube start

  # Linux
  curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
  sudo install minikube-linux-amd64 /usr/local/bin/minikube
  minikube start
  ```

- **kind** (Kubernetes in Docker): https://kind.sigs.k8s.io/
  ```bash
  brew install kind
  kind create cluster --name orchestra
  ```

- **Docker Desktop**: Включи Kubernetes в налаштуваннях

**Cloud (для production):**
- **GKE** (Google Kubernetes Engine)
- **EKS** (Amazon Elastic Kubernetes Service)
- **AKS** (Azure Kubernetes Service)

### Крок 2: Перевірити підключення

```bash
# Перевірити чи працює
kubectl cluster-info

# Перевірити ноди
kubectl get nodes
```

### Крок 3: Зібіл Docker образи

```bash
# З корня проекту
cd ../

# Backend
docker build -f backend/Dockerfile -t orchestra-backend:v1 backend/

# Frontend
docker build -f frontend/Dockerfile -t orchestra-frontend:v1 frontend/

# Для minikube - завантаж образи в minikube
minikube image load orchestra-backend:v1
minikube image load orchestra-frontend:v1

# Або для kind
kind load docker-image orchestra-backend:v1 --name orchestra
kind load docker-image orchestra-frontend:v1 --name orchestra
```

### Крок 4: Створити ресурси

```bash
cd k8s/

# 1. Backend
kubectl apply -f backend/configmap.yaml
kubectl apply -f backend/deployment.yaml
kubectl apply -f backend/service.yaml

# 2. Frontend
kubectl apply -f frontend/deployment.yaml
kubectl apply -f frontend/service.yaml

# 3. Ingress (опціонально)
kubectl apply -f frontend/ingress.yaml
```

**Або все разом:**
```bash
kubectl apply -f backend/
kubectl apply -f frontend/
```

### Крок 5: Перевірити

```bash
# Переглянути всі ресурси
kubectl get all

# Переглянути Pod'и
kubectl get pods

# Переглянути Services
kubectl get services

# Переглянути Deployments
kubectl get deployments

# Детальна інформація
kubectl describe deployment orchestra-backend
kubectl describe pod <pod-name>

# Логи
kubectl logs -f <pod-name>
kubectl logs -l app=orchestra,component=backend
```

### Крок 6: Доступ до застосунку

**Варіант 1: Port Forwarding (найпростіше)**
```bash
# Backend
kubectl port-forward svc/orchestra-backend 8383:80

# Frontend
kubectl port-forward svc/orchestra-frontend 8080:80

# Тепер відкрий:
# Backend: http://localhost:8383/congratulation
# Frontend: http://localhost:8080
```

**Варіант 2: NodePort (для локального доступу)**
```bash
# Для minikube
minikube service orchestra-frontend

# Це відкриє браузер з правильним URL
```

**Варіант 3: Ingress (production-like)**
```bash
# Встановити Nginx Ingress Controller
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/cloud/deploy.yaml

# Почекати поки запуститься
kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=120s

# Застосувати Ingress
kubectl apply -f frontend/ingress.yaml

# Дізнатись IP
kubectl get ingress

# Для minikube
minikube tunnel

# Відкрити в браузері
http://<INGRESS-IP>
```

## 📁 Структура файлів

```
k8s/
├── K8S_BASICS.md                  ← Навчальний курс (ПОЧИНАЙ ТУТ!)
├── README.md                      ← Цей файл
│
├── backend/
│   ├── deployment.yaml            ← Backend Pods (ДЕТАЛЬНІ КОМЕНТАРІ)
│   ├── service.yaml               ← Backend доступ
│   └── configmap.yaml             ← Конфігурація
│
├── frontend/
│   ├── deployment.yaml            ← Frontend Pods
│   ├── service.yaml               ← Frontend доступ
│   └── ingress.yaml               ← Зовнішній HTTP доступ
│
└── common/
    └── namespace.yaml             ← Namespaces (dev/staging/prod)
```

## 🎯 Послідовність навчання

### Початківець:

1. **Прочитай [K8S_BASICS.md](K8S_BASICS.md)** - основи YAML та K8s
2. **Вивчи [backend/deployment.yaml](backend/deployment.yaml)** - детальні коментарі
3. **Вивчи [backend/service.yaml](backend/service.yaml)** - як працює Service
4. **Спробуй запустити** - швидкий старт вище
5. **Експериментуй** - зміни replicas, resources, etc.

### Середній рівень:

1. **ConfigMap** - [backend/configmap.yaml](backend/configmap.yaml)
2. **Ingress** - [frontend/ingress.yaml](frontend/ingress.yaml)
3. **Namespaces** - [common/namespace.yaml](common/namespace.yaml)
4. **Масштабування** - `kubectl scale`
5. **Rolling updates** - `kubectl rollout`

### Просунутий:

1. **Resource Quotas** - ліміти на namespace
2. **Network Policies** - ізоляція мережі
3. **StatefulSets** - для stateful застосунків
4. **Persistent Volumes** - для даних
5. **Helm charts** - package manager

## 📝 Корисні команди

### Основні операції:

```bash
# Створити ресурси
kubectl apply -f deployment.yaml
kubectl apply -f backend/  # всі файли в директорії

# Переглянути ресурси
kubectl get all
kubectl get pods
kubectl get deployments
kubectl get services
kubectl get ingress

# Детальна інформація
kubectl describe pod <pod-name>
kubectl describe deployment <deployment-name>

# Логи
kubectl logs <pod-name>
kubectl logs -f <pod-name>  # follow
kubectl logs -l app=backend  # по labels

# Видалити
kubectl delete -f deployment.yaml
kubectl delete deployment <name>
kubectl delete pod <name>
```

### Налагодження:

```bash
# Виконати команду в Pod
kubectl exec <pod-name> -- ls -la
kubectl exec -it <pod-name> -- bash

# Port forwarding
kubectl port-forward pod/<pod-name> 8080:8383
kubectl port-forward svc/<service-name> 8080:80

# Копіювати файли
kubectl cp <pod-name>:/path/to/file ./local-file
kubectl cp ./local-file <pod-name>:/path/to/file

# Переглянути events
kubectl get events
kubectl get events --sort-by='.lastTimestamp'

# Перевірити чому Pod не запускається
kubectl describe pod <pod-name>
kubectl logs <pod-name>
```

### Масштабування:

```bash
# Змінити кількість replicas
kubectl scale deployment orchestra-backend --replicas=5

# Autoscaling (HPA)
kubectl autoscale deployment orchestra-backend --min=2 --max=10 --cpu-percent=80

# Переглянути HPA
kubectl get hpa
```

### Оновлення:

```bash
# Оновити образ
kubectl set image deployment/orchestra-backend backend=orchestra-backend:v2

# Переглянути статус rollout
kubectl rollout status deployment/orchestra-backend

# Переглянути історію
kubectl rollout history deployment/orchestra-backend

# Rollback
kubectl rollout undo deployment/orchestra-backend
kubectl rollout undo deployment/orchestra-backend --to-revision=2
```

### Secrets:

```bash
# Створити Secret
kubectl create secret generic db-secret \
  --from-literal=username=admin \
  --from-literal=password=secret123

# Переглянути Secret (закодовано в base64)
kubectl get secret db-secret -o yaml

# Decode Secret
kubectl get secret db-secret -o jsonpath='{.data.password}' | base64 -d
```

## 🐛 Troubleshooting

### Pod не запускається

```bash
# 1. Переглянути статус
kubectl get pods

# 2. Детальна інформація
kubectl describe pod <pod-name>

# 3. Логи
kubectl logs <pod-name>

# 4. Попередні логи (якщо Pod перезапустився)
kubectl logs <pod-name> --previous
```

**Частіпричини:**
- **ImagePullBackOff** - не може завантажити образ
  - Перевір назву образу
  - Для локальних образів: завантаж в minikube/kind
- **CrashLoopBackOff** - застосунок падає після старту
  - Перевір логи: `kubectl logs <pod-name>`
  - Перевір команду запуску
- **Pending** - не може знайти ноду
  - Недостатньо ресурсів
  - Перевір resources requests/limits

### Service не працює

```bash
# 1. Перевірити Service
kubectl get svc <service-name>
kubectl describe svc <service-name>

# 2. Перевірити Endpoints
kubectl get endpoints <service-name>

# 3. Перевірити селектор
kubectl get pods --show-labels
```

**Причини:**
- Selector не збігається з labels Pod'ів
- Pod'и не running
- Неправильні порти

### Ingress не працює

```bash
# 1. Перевірити чи встановлений Ingress Controller
kubectl get pods -n ingress-nginx

# 2. Переглянути Ingress
kubectl describe ingress <ingress-name>

# 3. Логи Ingress Controller
kubectl logs -n ingress-nginx -l app.kubernetes.io/name=ingress-nginx
```

## 🔄 Workflow для різних environments

### Development:

```bash
# Використати namespace
kubectl create namespace orchestra-dev
kubectl config set-context --current --namespace=orchestra-dev

# Застосувати конфігурацію
kubectl apply -f backend/ -n orchestra-dev
kubectl apply -f frontend/ -n orchestra-dev

# Port forwarding для доступу
kubectl port-forward svc/orchestra-frontend 8080:80 -n orchestra-dev
```

### Staging:

```bash
kubectl create namespace orchestra-staging
kubectl apply -f backend/ -n orchestra-staging
kubectl apply -f frontend/ -n orchestra-staging
```

### Production:

```bash
kubectl create namespace orchestra-prod

# Для production - використай Ingress з SSL
kubectl apply -f common/namespace.yaml
kubectl apply -f backend/ -n orchestra-prod
kubectl apply -f frontend/ -n orchestra-prod
kubectl apply -f frontend/ingress.yaml -n orchestra-prod
```

## 📚 Додаткові ресурси

### Офіційна документація:
- https://kubernetes.io/docs/
- https://kubernetes.io/docs/tutorials/

### Навчальні ресурси:
- https://kubernetes.io/docs/tutorials/kubernetes-basics/
- https://www.katacoda.com/courses/kubernetes

### Інструменти:
- **kubectl** - CLI для K8s
- **k9s** - Terminal UI для K8s
- **lens** - Desktop IDE для K8s
- **helm** - Package manager для K8s

## 💡 Поради

1. **Починай з [K8S_BASICS.md](K8S_BASICS.md)** - там все детально пояснено
2. **Читай коментарі в YAML файлах** - вони дуже детальні
3. **Експериментуй** - змінюй значення і дивись що станеться
4. **Використовуй `kubectl explain`** - вбудована документація:
   ```bash
   kubectl explain deployment
   kubectl explain deployment.spec
   kubectl explain deployment.spec.template.spec.containers
   ```
5. **Дивись логи** - `kubectl logs -f <pod-name>`
6. **Використовуй labels** - для організації та селекції
7. **Namespace для ізоляції** - dev/staging/prod

## ❓ Питання?

Якщо щось незрозуміло:
1. Прочитай [K8S_BASICS.md](K8S_BASICS.md)
2. Подивись коментарі в конкретному YAML файлі
3. Використай `kubectl explain <resource>`
4. Перевір офіційну документацію

**Happy Kubernetes learning!** 🚀
