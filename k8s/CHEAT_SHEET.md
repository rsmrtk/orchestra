# Kubernetes Cheat Sheet

## 🎯 Швидкий довідник команд

### kubectl основи

```bash
# Версія
kubectl version

# Інформація про кластер
kubectl cluster-info

# Ноди
kubectl get nodes
kubectl describe node <node-name>
```

### Основні ресурси

```bash
# Переглянути всі ресурси
kubectl get all
kubectl get all -n <namespace>

# Pod'и
kubectl get pods
kubectl get pods -o wide  # з IP та нодами
kubectl get pods --show-labels
kubectl get pods -l app=backend  # по labels
kubectl describe pod <pod-name>
kubectl logs <pod-name>
kubectl logs -f <pod-name>  # follow
kubectl logs <pod-name> --previous  # попередній контейнер
kubectl logs -l app=backend  # всі Pod'и з label

# Deployments
kubectl get deployments
kubectl get deploy  # скорочення
kubectl describe deployment <name>
kubectl edit deployment <name>

# Services
kubectl get services
kubectl get svc  # скорочення
kubectl describe svc <name>
kubectl get endpoints

# ConfigMaps
kubectl get configmaps
kubectl get cm  # скорочення
kubectl describe cm <name>

# Secrets
kubectl get secrets
kubectl describe secret <name>
kubectl get secret <name> -o yaml

# Ingress
kubectl get ingress
kubectl get ing  # скорочення
kubectl describe ingress <name>
```

### Створення/Оновлення/Видалення

```bash
# Apply (створити або оновити)
kubectl apply -f file.yaml
kubectl apply -f directory/

# Create (тільки створити, помилка якщо існує)
kubectl create -f file.yaml

# Replace (замінити існуючий)
kubectl replace -f file.yaml

# Delete
kubectl delete -f file.yaml
kubectl delete pod <name>
kubectl delete deployment <name>
kubectl delete all --all  # ВСЕ! (обережно)
```

### Налагодження

```bash
# Виконати команду в Pod
kubectl exec <pod-name> -- ls -la
kubectl exec -it <pod-name> -- bash
kubectl exec -it <pod-name> -- sh

# Port forwarding
kubectl port-forward pod/<pod-name> 8080:8383
kubectl port-forward svc/<service-name> 8080:80
kubectl port-forward deployment/<name> 8080:8383

# Копіювати файли
kubectl cp <pod-name>:/path/to/file ./local
kubectl cp ./local <pod-name>:/path/to/file

# Events
kubectl get events
kubectl get events --sort-by='.lastTimestamp'
kubectl get events -w  # watch

# Top (потрібен metrics-server)
kubectl top nodes
kubectl top pods
```

### Масштабування

```bash
# Manual scaling
kubectl scale deployment <name> --replicas=5

# Autoscaling
kubectl autoscale deployment <name> --min=2 --max=10 --cpu-percent=80
kubectl get hpa  # horizontal pod autoscaler
```

### Оновлення та Rollback

```bash
# Оновити образ
kubectl set image deployment/<name> <container-name>=<image>:<tag>

# Rollout статус
kubectl rollout status deployment/<name>
kubectl rollout history deployment/<name>

# Pause/Resume
kubectl rollout pause deployment/<name>
kubectl rollout resume deployment/<name>

# Rollback
kubectl rollout undo deployment/<name>
kubectl rollout undo deployment/<name> --to-revision=2

# Restart (zero-downtime)
kubectl rollout restart deployment/<name>
```

### Namespaces

```bash
# Список namespaces
kubectl get namespaces
kubectl get ns  # скорочення

# Створити namespace
kubectl create namespace <name>

# Працювати з namespace
kubectl get pods -n <namespace>
kubectl get all -n <namespace>
kubectl apply -f file.yaml -n <namespace>

# Змінити default namespace
kubectl config set-context --current --namespace=<name>
kubectl config view --minify | grep namespace

# Видалити namespace (і всі ресурси в ньому!)
kubectl delete namespace <name>
```

### Labels та Selectors

```bash
# Додати label
kubectl label pod <name> environment=production
kubectl label pod <name> key=value --overwrite

# Видалити label
kubectl label pod <name> environment-

# Фільтрувати по labels
kubectl get pods -l app=backend
kubectl get pods -l 'app=backend,version=v1'
kubectl get pods -l 'app in (backend,frontend)'
kubectl get pods -l 'environment!=production'

# Показати labels
kubectl get pods --show-labels
```

### Contexts

```bash
# Список contexts
kubectl config get-contexts
kubectl config current-context

# Переключитись між кластерами
kubectl config use-context <context-name>

# Set namespace
kubectl config set-context --current --namespace=<name>
```

### Створення ресурсів з командної строки

```bash
# Deployment
kubectl create deployment nginx --image=nginx:latest

# Service
kubectl expose deployment nginx --port=80 --type=LoadBalancer

# ConfigMap
kubectl create configmap app-config \
  --from-literal=key1=value1 \
  --from-literal=key2=value2

kubectl create configmap app-config \
  --from-file=config.yaml

# Secret
kubectl create secret generic db-secret \
  --from-literal=username=admin \
  --from-literal=password=pass123

kubectl create secret generic db-secret \
  --from-file=./username.txt \
  --from-file=./password.txt

# Namespace
kubectl create namespace <name>
```

### Інформація та документація

```bash
# Вбудована документація
kubectl explain pod
kubectl explain pod.spec
kubectl explain deployment.spec.template.spec.containers

# API resources
kubectl api-resources
kubectl api-resources --namespaced=true
kubectl api-resources --namespaced=false

# API versions
kubectl api-versions
```

### Корисні прапорці

```bash
# Вивід
-o wide                    # додаткова інформація
-o yaml                    # YAML формат
-o json                    # JSON формат
-o jsonpath='{.spec}'      # JSONPath query
--show-labels              # показати labels
--sort-by='.metadata.name' # сортувати
-w, --watch                # watch для змін

# Селекція
-l, --selector             # label selector
-n, --namespace            # namespace
--all-namespaces           # всі namespaces
-A                         # скорочення для --all-namespaces

# Інше
--dry-run=client           # не виконувати, тільки показати
--force                    # force операцію
-f                         # файл
-R                         # рекурсивно
```

## 📝 Шаблони YAML

### Мінімальний Pod

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nginx
spec:
  containers:
  - name: nginx
    image: nginx:latest
```

### Мінімальний Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:latest
        ports:
        - containerPort: 80
```

### Мінімальний Service

```yaml
apiVersion: v1
kind: Service
metadata:
  name: nginx
spec:
  selector:
    app: nginx
  ports:
  - port: 80
    targetPort: 80
```

### ConfigMap

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  KEY: "value"
  config.json: |
    {
      "key": "value"
    }
```

### Secret

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: app-secrets
type: Opaque
data:
  username: YWRtaW4=  # base64 encoded
  password: cGFzczEyMw==
stringData:
  # Звичайний текст (K8s encode автоматично)
  api-key: my-secret-key
```

## 🔍 Troubleshooting Flow

### Pod не запускається:

```bash
1. kubectl get pods                       # Статус?
2. kubectl describe pod <name>            # Events? Error?
3. kubectl logs <name>                    # Логи застосунку
4. kubectl logs <name> --previous         # Логи попереднього запуску
5. kubectl get events | grep <name>       # Events кластера
```

### Service не працює:

```bash
1. kubectl get svc <name>                 # Існує?
2. kubectl get endpoints <name>           # Є endpoints?
3. kubectl get pods -l <selector>         # Pod'и running?
4. kubectl describe svc <name>            # Selector правильний?
```

### Ingress не працює:

```bash
1. kubectl get ingress                    # ADDRESS заповнений?
2. kubectl describe ingress <name>        # Events? Error?
3. kubectl get svc                        # Services існують?
4. kubectl logs -n ingress-nginx -l app.kubernetes.io/name=ingress-nginx
```

## 💾 Збереження конфігурації

```bash
# Експорт існуючого ресурсу
kubectl get deployment nginx -o yaml > deployment.yaml
kubectl get service nginx -o yaml > service.yaml

# Dry run для генерації YAML
kubectl create deployment nginx --image=nginx --dry-run=client -o yaml > deployment.yaml
kubectl expose deployment nginx --port=80 --dry-run=client -o yaml > service.yaml
```

## 🎭 Контексти та кластери

```bash
# Додати кластер
kubectl config set-cluster my-cluster \
  --server=https://1.2.3.4 \
  --certificate-authority=ca.crt

# Додати користувача
kubectl config set-credentials my-user \
  --client-certificate=client.crt \
  --client-key=client.key

# Додати context
kubectl config set-context my-context \
  --cluster=my-cluster \
  --user=my-user \
  --namespace=default

# Використати context
kubectl config use-context my-context
```

## 🔐 RBAC (Role-Based Access Control)

```bash
# Переглянути ролі
kubectl get roles
kubectl get clusterroles

# Переглянути role bindings
kubectl get rolebindings
kubectl get clusterrolebindings

# Перевірити дозволи
kubectl auth can-i create deployments
kubectl auth can-i create pods --namespace=dev
kubectl auth can-i '*' '*'  # all permissions
```

## ⚡ Швидкі alias

Додай в `~/.bashrc` або `~/.zshrc`:

```bash
alias k='kubectl'
alias kg='kubectl get'
alias kd='kubectl describe'
alias kl='kubectl logs'
alias kex='kubectl exec -it'
alias kap='kubectl apply -f'
alias kdel='kubectl delete'
alias kgp='kubectl get pods'
alias kgs='kubectl get svc'
alias kgd='kubectl get deployments'
alias kn='kubectl config set-context --current --namespace'
```

Тепер можна:
```bash
k get pods
kg svc
kl <pod-name> -f
```

## 📊 Корисні JSONPath запити

```bash
# Отримати назви всіх Pod'ів
kubectl get pods -o jsonpath='{.items[*].metadata.name}'

# Отримати IP всіх Pod'ів
kubectl get pods -o jsonpath='{.items[*].status.podIP}'

# Отримати образи всіх контейнерів
kubectl get pods -o jsonpath='{.items[*].spec.containers[*].image}'

# Custom columns
kubectl get pods -o custom-columns=NAME:.metadata.name,STATUS:.status.phase,IP:.status.podIP
```

## 🚀 Швидкі команди для Orchestra

```bash
# Застосувати всю конфігурацію
kubectl apply -f k8s/backend/
kubectl apply -f k8s/frontend/

# Переглянути все
kubectl get all -l app=orchestra

# Логи backend
kubectl logs -l app=orchestra,component=backend -f

# Логи frontend
kubectl logs -l app=orchestra,component=frontend -f

# Port forward backend
kubectl port-forward svc/orchestra-backend 8383:80

# Port forward frontend
kubectl port-forward svc/orchestra-frontend 8080:80

# Масштабувати backend
kubectl scale deployment orchestra-backend --replicas=5

# Оновити backend
kubectl set image deployment/orchestra-backend backend=orchestra-backend:v2

# Перезапустити
kubectl rollout restart deployment/orchestra-backend
kubectl rollout restart deployment/orchestra-frontend
```

---

**Зберігай це як закладку!** 📖
