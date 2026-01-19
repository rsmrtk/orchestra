# 🚀 Швидкий старт Kubernetes для Orchestra Backend

> Покроковий практичний гайд для вивчення Kubernetes

---

## ⚡ Швидкий старт - Порядок роботи з файлами

**Шлях до файлів:** `/home/r/fff/deploy/orchestra/backend/deployments/staging/k8s/`

### 📝 Порядок вивчення та застосування:

```
Тиждень 1: Основи
├─ 1. namespace.yaml        ← Почни звідси! Створює ізольоване середовище
├─ 2. configmap.yaml        ← Конфігурація додатку (без секретів)
└─ 3. secret.yaml           ← Секрети (створи через kubectl, не файл!)

Тиждень 2: Основні ресурси (найважливіше!)
├─ 4. deployment.yaml       ← ⭐ ГОЛОВНИЙ! Запуск твого додатку
└─ 5. service.yaml          ← Мережевий доступ до Pod'ів

Тиждень 3: Розширені можливості
├─ 6. ingress.yaml          ← HTTP/HTTPS доступ ззовні (опційно)
├─ 7. hpa.yaml              ← Автомасштабування (опційно)
├─ 8. pdb.yaml              ← Захист від збоїв (опційно)
└─ 9. network-policy.yaml   ← Firewall rules (опційно)

Додаткові файли:
├─ kustomization.yaml       ← Для deploy всього одразу
├─ Makefile                 ← Автоматизація команд
├─ LINK.md                  ← Цей файл (швидкий старт)
├─ KUBERNETES_GUIDE_UA.md   ← Детальний теоретичний гайд
└─ README.md                ← Загальна інформація
```

### 🎯 Мінімальний набір для роботи:

Щоб запустити додаток, **обов'язково** потрібні:
1. ✅ `namespace.yaml` - ізоляція
2. ✅ `configmap.yaml` - конфігурація
3. ✅ Secret (створити через `kubectl create secret`)
4. ✅ `deployment.yaml` - запуск Pod'ів
5. ✅ `service.yaml` - мережевий доступ

Решта - опційно, для розширених можливостей.

### 📂 Швидкі посилання:

```bash
# Перейти до папки
cd /home/r/fff/deploy/orchestra/backend/deployments/staging/k8s

# Подивитися всі файли
ls -la

# Прочитати конкретний файл
cat namespace.yaml
cat deployment.yaml

# Застосувати конкретний файл
kubectl apply -f namespace.yaml
```

---

## 📋 Зміст

1. [Порядок вивчення](#-порядок-вивчення)
2. [Тиждень 1: Основи](#-тиждень-1-основи)
3. [Тиждень 2: Основні ресурси](#-тиждень-2-основні-ресурси)
4. [Тиждень 3: Розширені можливості](#-тиждень-3-розширені-можливості)
5. [Корисні команди](#-корисні-команди)
6. [Швидкий deploy всього](#-швидкий-deploy-всього)

---

## 🎯 Порядок вивчення

```
1. namespace.yaml        → Ізоляція середовища
2. configmap.yaml        → Конфігурація
3. secret.yaml           → Секрети
4. deployment.yaml       → ⭐ Запуск Pod'ів (найважливіший!)
5. service.yaml          → Мережевий доступ
6. ingress.yaml          → HTTP/HTTPS ззовні
7. hpa.yaml              → Автомасштабування
8. pdb.yaml              → Захист від збоїв
9. network-policy.yaml   → Firewall rules
```

---

## 📅 Тиждень 1: Основи

### 1️⃣ Namespace

**Що це:** Ізоляція середовищ (staging окремо від production)

```bash
# Створити
kubectl apply -f namespace.yaml

# Подивитися всі namespaces
kubectl get namespaces
kubectl get ns  # коротка форма

# Деталі
kubectl describe namespace orchestra-staging

# Список ресурсів (поки пусто)
kubectl get all -n orchestra-staging
```

---

### 2️⃣ ConfigMapї

**Що це:** Конфігурація додатку (env змінні, файли)

```bash
# Створити
kubectl apply -f configmap.yaml

# Подивитися
kubectl get configmap -n orchestra-staging
kubectl get cm -n orchestra-staging  # коротка форма

# Деталі + всі ключі
kubectl describe configmap orchestra-backend-config -n orchestra-staging

# У форматі YAML (всі значення)
kubectl get configmap orchestra-backend-config -n orchestra-staging -o yaml

# Експеримент: змінити значення
# 1. Відредагуй configmap.yaml (наприклад GIN_MODE: "release")
# 2. kubectl apply -f configmap.yaml
# 3. kubectl rollout restart deployment/orchestra-backend -n orchestra-staging
```

**Ключові поля:**
- `data:` - ключі та значення конфігурації
- Використовується в Deployment через `configMapKeyRef`

---

### 3️⃣ Secret

**Що це:** Секретна інформація (паролі, токени, API ключі)

⚠️ **ВАЖЛИВО:** Secret НЕ зашифрований за замовчуванням (тільки base64)!

```bash
# ✅ РЕКОМЕНДОВАНИЙ СПОСІБ - через kubectl
kubectl create secret generic orchestra-backend-secret \
  --from-literal=API_KEY=test-key-123 \
  --from-literal=DB_PASSWORD=test-pass-456 \
  -n orchestra-staging

# Подивитися
kubectl get secret -n orchestra-staging
kubectl describe secret orchestra-backend-secret -n orchestra-staging

# Подивитися закодовані значення (base64)
kubectl get secret orchestra-backend-secret -n orchestra-staging -o yaml

# Декодувати конкретний ключ
kubectl get secret orchestra-backend-secret -n orchestra-staging \
  -o jsonpath='{.data.API_KEY}' | base64 -d
echo  # newline

# Вручну закодувати/декодувати
echo -n "my-secret-value" | base64          # Закодувати
echo "bXktc2VjcmV0LXZhbHVl" | base64 -d    # Декодувати
```

**Не робіть так:**
- ❌ Не комітити реальні секрети в Git
- ❌ Не вважати base64 за шифрування
- ✅ Використовувати Sealed Secrets, External Secrets, Vault для production

---

## 📅 Тиждень 2: Основні ресурси

### 4️⃣ Deployment ⭐ **НАЙВАЖЛИВІШИЙ!**

**Що це:** Управління Pod'ами (створення, масштабування, оновлення)

#### Підготовка

```bash
# СПОЧАТКУ: замінити Docker image на тестовий
# Відредагуй deployment.yaml, знайди рядок:
#   image: your-registry/orchestra-backend:staging
# Заміни на:
#   image: nginx:alpine
```

#### Базові операції

```bash
# Створити
kubectl apply -f deployment.yaml

# Статус
kubectl get deployments -n orchestra-staging
kubectl get deploy -n orchestra-staging  # коротка форма

# Дивитися Pod'и у реальному часі
kubectl get pods -n orchestra-staging -w  # -w = watch, Ctrl+C для виходу

# Деталі Deployment
kubectl describe deployment orchestra-backend -n orchestra-staging

# Деталі конкретного Pod'а
kubectl get pods -n orchestra-staging  # скопіюй ім'я Pod'а
kubectl describe pod <pod-name> -n orchestra-staging

# Логи
kubectl logs -f deployment/orchestra-backend -n orchestra-staging
# або конкретного Pod'а:
kubectl logs -f <pod-name> -n orchestra-staging
# попередній контейнер (якщо Pod перезапускався):
kubectl logs <pod-name> -n orchestra-staging --previous
```

#### Масштабування

```bash
# Збільшити до 3 Pod'ів
kubectl scale deployment orchestra-backend --replicas=3 -n orchestra-staging
kubectl get pods -n orchestra-staging

# Зменшити до 2 Pod'ів
kubectl scale deployment orchestra-backend --replicas=2 -n orchestra-staging
kubectl get pods -n orchestra-staging -w  # 1 Pod видалиться

# Назад до 1
kubectl scale deployment orchestra-backend --replicas=1 -n orchestra-staging
```

#### Оновлення (Rolling Update)

```bash
# Оновити образ (симуляція deploy нової версії)
kubectl set image deployment/orchestra-backend \
  backend=nginx:1.25 \
  -n orchestra-staging

# Дивитися процес оновлення
kubectl rollout status deployment/orchestra-backend -n orchestra-staging

# Дивитися як Pod'и оновлюються по черзі
kubectl get pods -n orchestra-staging -w
```

#### Rollback (відкат)

```bash
# Відкатити до попередньої версії
kubectl rollout undo deployment/orchestra-backend -n orchestra-staging

# Історія deploy'ів
kubectl rollout history deployment/orchestra-backend -n orchestra-staging

# Відкат до конкретної ревізії
kubectl rollout undo deployment/orchestra-backend --to-revision=2 -n orchestra-staging
```

#### Перезапуск

```bash
# Перезапустити всі Pod'и (без зміни конфігурації)
kubectl rollout restart deployment/orchestra-backend -n orchestra-staging
```

#### Debug

```bash
# Shell всередині Pod'а (як SSH)
kubectl exec -it <pod-name> -n orchestra-staging -- sh
ls /
ps aux
exit

# Скопіювати файл з Pod'а
kubectl cp <pod-name>:/path/to/file ./local-file -n orchestra-staging

# Скопіювати файл в Pod
kubectl cp ./local-file <pod-name>:/path/to/file -n orchestra-staging
```

#### Troubleshooting

```bash
# STATUS: ImagePullBackOff (не може завантажити образ)
kubectl describe pod <pod-name> -n orchestra-staging
# Дивись Events внизу - там причина
# Рішення: заміни image на nginx:alpine

# STATUS: CrashLoopBackOff (Pod падає)
kubectl logs <pod-name> -n orchestra-staging
kubectl logs <pod-name> -n orchestra-staging --previous
# Дивись помилку в логах

# STATUS: Pending (довго чекає)
kubectl describe pod <pod-name> -n orchestra-staging
# Можливо недостатньо ресурсів на Node
```

**Ключові поля в deployment.yaml:**
- `spec.replicas` - кількість Pod'ів
- `spec.template.spec.containers[].image` - Docker образ
- `spec.template.spec.containers[].resources` - CPU/Memory requests/limits
- `spec.template.spec.containers[].livenessProbe` - чи живий контейнер
- `spec.template.spec.containers[].readinessProbe` - чи готовий приймати трафік

---

### 5️⃣ Service

**Що це:** Стабільний мережевий доступ до Pod'ів (Load Balancer + DNS)

```bash
# Створити
kubectl apply -f service.yaml

# Подивитися
kubectl get services -n orchestra-staging
kubectl get svc -n orchestra-staging  # коротка форма

# Деталі
kubectl describe service orchestra-backend-service -n orchestra-staging

# ⭐ ДУЖЕ ВАЖЛИВО: Перевірити Endpoints
kubectl get endpoints orchestra-backend-service -n orchestra-staging
# МАЄ показати IP адреси Pod'ів!
# Якщо пусто → selector не співпадає з labels Pod'ів
```

#### Тестування

```bash
# Метод 1: Port-forward (локальний доступ)
kubectl port-forward svc/orchestra-backend-service 8080:80 -n orchestra-staging
# У ІНШОМУ терміналі:
curl http://localhost:8080
# Ctrl+C для зупинки port-forward

# Метод 2: З іншого Pod'а (тестування DNS)
kubectl run test --image=alpine -n orchestra-staging --rm -it -- sh
apk add curl
curl http://orchestra-backend-service:80
# Працює! Service знайдений через DNS
exit

# Метод 3: Перевірка DNS резолюції
kubectl run test --image=alpine -n orchestra-staging --rm -it -- sh
apk add bind-tools
nslookup orchestra-backend-service
# Має показати Cluster IP Service
exit
```

#### Troubleshooting

```bash
# Проблема: Endpoints пусті
# Перевірити labels Pod'ів
kubectl get pods -n orchestra-staging --show-labels

# Перевірити selector Service
kubectl get svc orchestra-backend-service -n orchestra-staging -o yaml | grep -A 5 selector

# Labels Pod'ів МАЮТЬ співпадати з selector Service!
```

**DNS імена:**
- `orchestra-backend-service` - в тому самому namespace
- `orchestra-backend-service.orchestra-staging` - з іншого namespace
- `orchestra-backend-service.orchestra-staging.svc.cluster.local` - повне FQDN

**Ключові поля в service.yaml:**
- `spec.type` - ClusterIP, NodePort, LoadBalancer
- `spec.selector` - як знайти Pod'и (по labels)
- `spec.ports[].port` - порт Service
- `spec.ports[].targetPort` - порт контейнера

---

## 📅 Тиждень 3: Розширені можливості

### 6️⃣ Ingress

**Що це:** HTTP/HTTPS маршрутизація ззовні кластера (Reverse Proxy)

#### Підготовка: Встановити Ingress Controller

```bash
# Для Minikube
minikube addons enable ingress

# АБО для звичайного кластера
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.8.1/deploy/static/provider/cloud/deploy.yaml

# Перевірити
kubectl get pods -n ingress-nginx
kubectl get svc -n ingress-nginx
# Дивись EXTERNAL-IP ingress-nginx-controller
```

#### Базові операції

```bash
# Створити
kubectl apply -f ingress.yaml

# Подивитися
kubectl get ingress -n orchestra-staging
kubectl get ing -n orchestra-staging  # коротка форма

# Деталі
kubectl describe ingress orchestra-backend-ingress -n orchestra-staging

# Дізнатися IP Ingress Controller
kubectl get svc -n ingress-nginx
# Дивись EXTERNAL-IP
```

#### Тестування (якщо є DNS)

```bash
# Якщо DNS налаштований:
curl http://staging-api.orchestra.example.com/api/congratulation

# Для локального тесту БЕЗ DNS:
# Додати в /etc/hosts:
echo "<EXTERNAL-IP> staging-api.orchestra.example.com" | sudo tee -a /etc/hosts
curl http://staging-api.orchestra.example.com/api/congratulation
```

**Ключові поля в ingress.yaml:**
- `spec.rules[].host` - домен
- `spec.rules[].http.paths[].path` - шлях (/api, /app)
- `spec.rules[].http.paths[].backend.service` - куди перенаправити
- `spec.tls` - SSL/HTTPS налаштування

---

### 7️⃣ HPA (Horizontal Pod Autoscaler)

**Що це:** Автоматичне масштабування на основі CPU/Memory

#### Підготовка: Встановити Metrics Server

```bash
# Для Minikube
minikube addons enable metrics-server

# АБО для звичайного кластера
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# Перевірити
kubectl get deployment metrics-server -n kube-system
kubectl top nodes
kubectl top pods -n orchestra-staging
```

#### Базові операції

```bash
# Створити
kubectl apply -f hpa.yaml

# Подивитися
kubectl get hpa -n orchestra-staging

# Деталі
kubectl describe hpa orchestra-backend-hpa -n orchestra-staging

# Дивитися у реальному часі
kubectl get hpa -n orchestra-staging --watch
```

#### Тестування (Load Test)

```bash
# Запустити load generator
kubectl run -it --rm load-generator --image=busybox -n orchestra-staging -- /bin/sh
while true; do wget -q -O- http://orchestra-backend-service; done

# У ІНШОМУ терміналі дивись:
kubectl get hpa -n orchestra-staging --watch
kubectl get pods -n orchestra-staging --watch
# Має з'явитися більше Pod'ів!

# Ctrl+C в load generator щоб зупинити
# Pod'и поступово зменшаться назад до minReplicas
```

**Ключові поля в hpa.yaml:**
- `spec.minReplicas` - мінімум Pod'ів (завжди)
- `spec.maxReplicas` - максимум Pod'ів
- `spec.metrics[].target.averageUtilization` - threshold (наприклад 70% CPU)

---

### 8️⃣ PDB (Pod Disruption Budget)

**Що це:** Гарантія мінімальної кількості Pod'ів під час maintenance

```bash
# Створити
kubectl apply -f pdb.yaml

# Подивитися
kubectl get pdb -n orchestra-staging
kubectl get poddisruptionbudget -n orchestra-staging  # повна назва

# Деталі
kubectl describe pdb orchestra-backend-pdb -n orchestra-staging
```

---

### 9️⃣ NetworkPolicy

**Що це:** Firewall rules для Pod'ів (хто може з'єднуватися)

⚠️ **Вимога:** CNI plugin з підтримкою NetworkPolicy (Calico, Cilium). Flannel НЕ підтримує!

```bash
# Створити
kubectl apply -f network-policy.yaml

# Подивитися
kubectl get networkpolicy -n orchestra-staging
kubectl get netpol -n orchestra-staging  # коротка форма

# Деталі
kubectl describe networkpolicy orchestra-backend-netpol -n orchestra-staging
```

---

## 🛠 Корисні команди

### Загальні

```bash
# Все в namespace
kubectl get all -n orchestra-staging

# Конкретні типи
kubectl get pods,svc,deploy,ing -n orchestra-staging

# Широкий вивід (більше інфо)
kubectl get pods -n orchestra-staging -o wide

# YAML вивід
kubectl get deployment orchestra-backend -n orchestra-staging -o yaml

# JSON вивід
kubectl get deployment orchestra-backend -n orchestra-staging -o json

# Watch (auto-refresh)
kubectl get pods -n orchestra-staging --watch

# Labels
kubectl get pods -n orchestra-staging --show-labels
kubectl get pods -n orchestra-staging -l app=orchestra-backend
```

### Debug

```bash
# Describe (деталі + Events)
kubectl describe <resource> <name> -n orchestra-staging

# Логи
kubectl logs <pod-name> -n orchestra-staging
kubectl logs <pod-name> -n orchestra-staging -f            # Follow (live)
kubectl logs <pod-name> -n orchestra-staging --previous     # Попередній контейнер
kubectl logs <pod-name> -n orchestra-staging --tail=100     # Останні 100 рядків

# Exec (shell в Pod)
kubectl exec <pod-name> -n orchestra-staging -- ls /app
kubectl exec -it <pod-name> -n orchestra-staging -- sh

# Port forward (локальний доступ)
kubectl port-forward <pod-name> 8080:8383 -n orchestra-staging
kubectl port-forward svc/<service-name> 8080:80 -n orchestra-staging

# Copy files
kubectl cp <pod-name>:/path/to/file ./local-file -n orchestra-staging
kubectl cp ./local-file <pod-name>:/path/to/file -n orchestra-staging

# Ресурси (CPU/Memory)
kubectl top nodes
kubectl top pods -n orchestra-staging
kubectl top pods -n orchestra-staging --sort-by=memory
kubectl top pods -n orchestra-staging --sort-by=cpu
```

### Управління

```bash
# Apply
kubectl apply -f file.yaml
kubectl apply -f directory/
kubectl apply -k directory/  # Kustomize

# Delete
kubectl delete -f file.yaml
kubectl delete pod <pod-name> -n orchestra-staging
kubectl delete deployment <name> -n orchestra-staging

# Edit (відкриє редактор)
kubectl edit deployment orchestra-backend -n orchestra-staging

# Patch (JSON)
kubectl patch deployment orchestra-backend -n orchestra-staging \
  -p '{"spec":{"replicas":5}}'

# Scale
kubectl scale deployment orchestra-backend --replicas=3 -n orchestra-staging

# Set image
kubectl set image deployment/orchestra-backend \
  backend=new-image:tag \
  -n orchestra-staging

# Rollout
kubectl rollout status deployment/orchestra-backend -n orchestra-staging
kubectl rollout history deployment/orchestra-backend -n orchestra-staging
kubectl rollout undo deployment/orchestra-backend -n orchestra-staging
kubectl rollout restart deployment/orchestra-backend -n orchestra-staging
```

### Events та Troubleshooting

```bash
# Події (відсортовані по часу)
kubectl get events -n orchestra-staging --sort-by='.lastTimestamp'

# Всі події в кластері
kubectl get events --all-namespaces --sort-by='.lastTimestamp'

# Перевірка здоров'я кластера
kubectl cluster-info
kubectl get nodes
kubectl get componentstatuses
kubectl get cs  # коротка форма

# Версія
kubectl version --short
```

### Context та Namespace

```bash
# Поточний context
kubectl config current-context

# Список contexts
kubectl config get-contexts

# Перемкнути context
kubectl config use-context <context-name>

# Встановити default namespace
kubectl config set-context --current --namespace=orchestra-staging

# Тепер можна без -n orchestra-staging
kubectl get pods
```

---

## 🚀 Швидкий deploy всього

### Метод 1: Окремі файли (по порядку)

```bash
cd /home/r/fff/deploy/orchestra/backend/deployments/staging/k8s

# 1. Основа
kubectl apply -f namespace.yaml
kubectl apply -f configmap.yaml

# 2. Secret (через kubectl, не файл!)
kubectl create secret generic orchestra-backend-secret \
  --from-literal=API_KEY=your-key \
  --from-literal=DB_PASSWORD=your-password \
  -n orchestra-staging

# 3. Pod'и та мережа
kubectl apply -f deployment.yaml
kubectl apply -f service.yaml

# 4. Розширені (опційно)
kubectl apply -f ingress.yaml
kubectl apply -f hpa.yaml
kubectl apply -f pdb.yaml
kubectl apply -f network-policy.yaml

# 5. Перевірити
kubectl get all -n orchestra-staging
kubectl get pods -n orchestra-staging -w
```

### Метод 2: Kustomize (все одразу)

```bash
cd /home/r/fff/deploy/orchestra/backend/deployments/staging/k8s

# Deploy всього одною командою
kubectl apply -k .

# Перевірити
kubectl get all -n orchestra-staging
```

### Метод 3: Makefile (якщо є)

```bash
cd /home/r/fff/deploy/orchestra/backend/deployments/staging/k8s

# Подивитися доступні команди
make help

# Deploy
make deploy

# Видалити
make delete
```

---

## 🧹 Cleanup (видалення всього)

```bash
# Видалити весь namespace (і все всередині)
kubectl delete namespace orchestra-staging

# Або окремо кожен ресурс
kubectl delete -f deployment.yaml
kubectl delete -f service.yaml
kubectl delete -f ingress.yaml
kubectl delete -f hpa.yaml
kubectl delete -f pdb.yaml
kubectl delete -f network-policy.yaml
kubectl delete -f configmap.yaml
kubectl delete secret orchestra-backend-secret -n orchestra-staging
kubectl delete -f namespace.yaml

# Через Kustomize
kubectl delete -k .
```

---

## 📊 Моніторинг та логи

```bash
# Real-time моніторинг
kubectl get pods -n orchestra-staging -w
kubectl get hpa -n orchestra-staging -w
kubectl top pods -n orchestra-staging --watch

# Агреговані логи всіх Pod'ів Deployment
kubectl logs -f deployment/orchestra-backend -n orchestra-staging

# Stern (якщо встановлений) - логи багатьох Pod'ів
stern orchestra-backend -n orchestra-staging

# Events (що відбувається)
kubectl get events -n orchestra-staging --watch
```

---

## ✅ Чеклист перед Production

- [ ] `resources.requests` та `resources.limits` вказані
- [ ] `livenessProbe` та `readinessProbe` налаштовані
- [ ] `replicas >= 2` для high availability
- [ ] PodDisruptionBudget створений
- [ ] HPA налаштований (якщо потрібно)
- [ ] Secrets НЕ в Git (використовуй External Secrets)
- [ ] NetworkPolicy налаштовані
- [ ] Ingress з SSL (cert-manager)
- [ ] Моніторинг та alerts налаштовані
- [ ] Backup стратегія є

---

## 📚 Додаткові ресурси

- **Детальний гайд:** `KUBERNETES_GUIDE_UA.md` (в цій же папці)
- **README:** `README.md` (в цій же папці)
- **Офіційна документація:** https://kubernetes.io/docs/
- **kubectl шпаргалка:** https://kubernetes.io/docs/reference/kubectl/cheatsheet/

---

**Зроблено з ❤️ для Orchestra Backend Team**
