# Kubernetes - Візуальний Guide

## 🎨 Як все працює разом

### Архітектура Orchestra в Kubernetes

```
┌─────────────────────────────────────────────────────────────────┐
│                         KUBERNETES CLUSTER                       │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │                    INGRESS CONTROLLER                       │ │
│  │                    (nginx-ingress)                          │ │
│  │  ┌───────────────┐  ┌──────────────────┐                  │ │
│  │  │ / → Frontend  │  │ /api → Backend   │                  │ │
│  │  └───────┬───────┘  └────────┬─────────┘                  │ │
│  └──────────┼────────────────────┼────────────────────────────┘ │
│             │                    │                               │
│  ┌──────────▼──────────┐  ┌─────▼───────────┐                  │
│  │  Frontend Service   │  │ Backend Service  │                  │
│  │   ClusterIP: 80     │  │  ClusterIP: 80   │                  │
│  └──────────┬──────────┘  └─────┬────────────┘                  │
│             │  selects by       │  selects by                   │
│             │  app=orchestra    │  app=orchestra                │
│             │  component=       │  component=backend            │
│             │  frontend         │                                │
│  ┌──────────▼──────────┐  ┌─────▼────────────┐                  │
│  │ Frontend Deployment │  │Backend Deployment│                  │
│  │   replicas: 2       │  │  replicas: 3     │                  │
│  └──────────┬──────────┘  └─────┬────────────┘                  │
│             │                    │                               │
│   ┌─────────┴────────┐   ┌──────┴──────────────┐               │
│   │                  │   │                      │               │
│ ┌─▼──┐  ┌────┐     ┌─▼──┐ ┌────┐  ┌────┐                       │
│ │Pod1│  │Pod2│     │Pod1│ │Pod2│  │Pod3│                       │
│ │    │  │    │     │    │ │    │  │    │                       │
│ │Nginx│  │Nginx│     │ Go │ │ Go │  │ Go │                       │
│ │5173│  │5173│     │8383│ │8383│  │8383│                       │
│ └────┘  └────┘     └────┘ └────┘  └────┘                       │
│                                                                  │
│                    ┌──────────────┐                              │
│                    │  ConfigMap   │                              │
│                    │backend-config│                              │
│                    └──────┬───────┘                              │
│                           │ provides env vars                    │
│                           │                                      │
│                    ┌──────▼───────┐                              │
│                    │    Secret    │                              │
│                    │(if needed)   │                              │
│                    └──────────────┘                              │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘

                           ┌────────────┐
                           │   User     │
                           │  Browser   │
                           └──────┬─────┘
                                  │
                                  │ http://cluster-ip/
                                  │
                           ┌──────▼─────┐
                           │  Ingress   │
                           └────────────┘
```

## 🔄 Request Flow

### 1. User открывает Frontend

```
1. Browser Request
   │
   │  GET http://orchestra.example.com/
   │
   ▼
2. Ingress
   │
   │  Rule: "/" → orchestra-frontend:80
   │
   ▼
3. Frontend Service
   │
   │  selector: app=orchestra, component=frontend
   │  Finds: Pod1, Pod2
   │  LoadBalance: round-robin
   │
   ▼
4. Frontend Pod (random: Pod1)
   │
   │  Nginx serves: /usr/share/nginx/html/index.html
   │
   ▼
5. React App loads in browser
```

### 2. React App викликає Backend API

```
1. React App
   │
   │  fetch('http://orchestra-backend/congratulation')
   │  (or через Ingress: /api/congratulation)
   │
   ▼
2. Kubernetes DNS
   │
   │  Resolve: orchestra-backend → ClusterIP (e.g., 10.96.1.100)
   │
   ▼
3. Backend Service
   │
   │  selector: app=orchestra, component=backend
   │  Finds: Pod1, Pod2, Pod3
   │  LoadBalance: round-robin
   │
   ▼
4. Backend Pod (random: Pod2)
   │
   │  Go Gin handles: GET /congratulation
   │  Returns: {"message": "Congratulation"}
   │
   ▼
5. Response to React App
   │
   │  JSON response
   │
   ▼
6. React updates UI
```

## 📦 Deployment Process

### Як створюється Deployment

```
kubectl apply -f deployment.yaml
            │
            ▼
┌───────────────────────────┐
│   Kubernetes API Server   │
│  "Create Deployment"      │
└─────────────┬─────────────┘
              │
              ▼
┌───────────────────────────┐
│  Deployment Controller    │
│  "Need 3 replicas"        │
└─────────────┬─────────────┘
              │
              ▼
┌───────────────────────────┐
│  ReplicaSet               │
│  "Create 3 Pods"          │
└─────────────┬─────────────┘
              │
        ┌─────┴─────┬─────┐
        │           │     │
        ▼           ▼     ▼
    ┌────┐      ┌────┐ ┌────┐
    │Pod1│      │Pod2│ │Pod3│
    └────┘      └────┘ └────┘
```

### Rolling Update

```
Before Update (v1):
┌────┐ ┌────┐ ┌────┐
│ v1 │ │ v1 │ │ v1 │
└────┘ └────┘ └────┘
   ▲      ▲      ▲
   └──────┴──────┘ Service

kubectl set image deployment/backend backend=backend:v2

Step 1: Create new Pod
┌────┐ ┌────┐ ┌────┐ ┌────┐
│ v1 │ │ v1 │ │ v1 │ │ v2 │ ← новий
└────┘ └────┘ └────┘ └────┘
   ▲      ▲      ▲      ▲
   └──────┴──────┴──────┘ Service

Step 2: Terminate old Pod
┌────┐ ┌────┐ ┌────┐
│ v1 │ │ v1 │ │ v2 │
└────┘ └────┘ └────┘
   ▲      ▲      ▲
   └──────┴──────┘ Service

Step 3: Continue...
┌────┐ ┌────┐ ┌────┐ ┌────┐
│ v1 │ │ v2 │ │ v2 │ │ v2 │ ← новий
└────┘ └────┘ └────┘ └────┘
   ▲      ▲      ▲      ▲
   └──────┴──────┴──────┘ Service

Finally (all v2):
┌────┐ ┌────┐ ┌────┐
│ v2 │ │ v2 │ │ v2 │
└────┘ └────┘ └────┘
   ▲      ▲      ▲
   └──────┴──────┘ Service
```

## 🌐 Service Types

### ClusterIP (internal only)

```
┌─────────────────────────────────┐
│        Kubernetes Cluster        │
│                                  │
│  ┌────────────┐                 │
│  │  Service   │                 │
│  │ ClusterIP  │ ← Internal IP   │
│  │ 10.96.1.10 │                 │
│  └─────┬──────┘                 │
│        │                         │
│   ┌────┴────┬────┐              │
│   │         │    │              │
│  Pod1     Pod2  Pod3            │
│                                  │
│  Доступ ТІЛЬКИ всередині кластера│
│                                  │
└─────────────────────────────────┘
```

### NodePort (external access via node IP)

```
┌─────────────────────────────────┐
│        Kubernetes Cluster        │
│                                  │
│  ┌────────────┐                 │
│  │  Service   │                 │
│  │  NodePort  │                 │
│  │  30080     │ ← Port on nodes │
│  └─────┬──────┘                 │
│        │                         │
│   ┌────┴────┬────┐              │
│   │         │    │              │
│  Pod1     Pod2  Pod3            │
│                                  │
└───────────────┬─────────────────┘
                │
         Доступ через:
         http://<NodeIP>:30080
```

### LoadBalancer (cloud provider)

```
                ┌─────────────┐
                │   Cloud LB  │ ← External IP
                │ 34.123.45.67│
                └──────┬──────┘
                       │
┌──────────────────────┼──────────┐
│   Kubernetes Cluster │          │
│               ┌──────▼──────┐   │
│               │  Service    │   │
│               │LoadBalancer │   │
│               └──────┬──────┘   │
│                      │           │
│                ┌─────┴─────┬───┐│
│                │           │   ││
│              Pod1        Pod2  Pod3
│                                  │
└─────────────────────────────────┘
```

## 🔍 Labels and Selectors

### Як Service знаходить Pod'и

```
Deployment YAML:
metadata:
  labels:                  ┌─ Це labels Deployment
    app: backend           │
    version: v1            │
template:                  │
  metadata:                │
    labels:                │  Service selector:
      app: backend    ◄────┼─  matchLabels:
      component: api  ◄────┤     app: backend
                           │     component: api
                           │
Pod's які створюються:    │
┌───────────────────┐     │
│ Pod1              │     │
│ labels:           │     │
│   app: backend ◄──┼─────┘
│   component: api ◄┼─────┘
└───────────────────┘
     ▲
     │ Service знаходить Pod
     │ по цим labels!
```

### Label Matching

```
Service Selector:
  app: backend
  component: api

┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│ Pod1         │  │ Pod2         │  │ Pod3         │
│ app: backend │  │ app: backend │  │ app: frontend│
│ component:api│  │ component:api│  │ component:ui │
└──────────────┘  └──────────────┘  └──────────────┘
       ✓                 ✓                  ✗
    MATCH             MATCH            NO MATCH
```

## 💾 ConfigMap та Environment Variables

### Як ConfigMap потрапляє в Pod

```
kubectl apply -f configmap.yaml
            │
            ▼
┌────────────────────────┐
│      ConfigMap         │
│  name: backend-config  │
│  data:                 │
│    PORT: "8383"        │
│    LOG_LEVEL: "info"   │
└────────────┬───────────┘
             │
             │ Referenced in Deployment
             │
    ┌────────▼─────────────────┐
    │    Deployment            │
    │  env:                    │
    │  - name: PORT            │
    │    valueFrom:            │
    │      configMapKeyRef:    │
    │        name: backend-cfg │
    │        key: PORT          │
    └────────┬─────────────────┘
             │
             │ Creates Pod with env
             │
    ┌────────▼─────────┐
    │      Pod         │
    │  ENV:            │
    │    PORT=8383     │ ← From ConfigMap!
    │    LOG_LEVEL=info│
    └──────────────────┘
```

## 🏥 Health Checks

### Liveness Probe

```
Pod Started
    │
    │ initialDelaySeconds: 30s
    ▼
┌─────────────┐
│   Wait 30s  │
└──────┬──────┘
       │
       │ Every periodSeconds: 10s
       ▼
┌──────────────────┐
│ GET /health      │
│ Timeout: 5s      │
└────┬──────┬──────┘
     │      │
  200 OK   500 Error
     │      │
     ▼      ▼
   ✓ OK   ✗ Fail
          │
          │ failureThreshold: 3
          ▼
    ┌─────────────┐
    │ Fail 3 times│
    └──────┬──────┘
           │
           ▼
    ┌──────────────┐
    │ Restart Pod  │ ← K8s перезапускає!
    └──────────────┘
```

### Readiness Probe

```
Pod Started
    │
    │ initialDelaySeconds: 5s
    ▼
┌─────────────┐
│   Wait 5s   │
└──────┬──────┘
       │
       │ Every periodSeconds: 5s
       ▼
┌──────────────────┐
│ GET /ready       │
└────┬──────┬──────┘
     │      │
  200 OK   503 Not Ready
     │      │
     ▼      ▼
   ✓ OK   ✗ Not Ready
   │      │
   │      │ Service REMOVES Pod from endpoints
   │      │ (не шле трафік)
   │      │
   │      ▼
   │  ┌─────────────────────┐
   │  │ Pod not in Service  │
   │  └─────────────────────┘
   │
   │ Service ADDS Pod to endpoints
   │ (шле трафік)
   │
   ▼
┌──────────────────────┐
│ Pod receives traffic │
└──────────────────────┘
```

## 🎯 Namespace Isolation

```
┌───────────────────────────────────────────────────────────┐
│                   Kubernetes Cluster                       │
│                                                            │
│  ┌──────────────────┐  ┌──────────────────┐              │
│  │ namespace: dev   │  │namespace: prod   │              │
│  │                  │  │                  │              │
│  │  Backend (v1)    │  │  Backend (v2)    │              │
│  │  Frontend        │  │  Frontend        │              │
│  │  ConfigMap       │  │  ConfigMap       │              │
│  │                  │  │  Secret          │              │
│  └──────────────────┘  └──────────────────┘              │
│         ▲                      ▲                          │
│         │                      │                          │
│    DNS: backend.dev      DNS: backend.prod               │
│                                                            │
│  ┌──────────────────┐                                     │
│  │namespace: staging│                                     │
│  │                  │                                     │
│  │  Backend (v2-rc) │                                     │
│  │  Frontend        │                                     │
│  └──────────────────┘                                     │
│                                                            │
└───────────────────────────────────────────────────────────┘
```

## 🔐 Resource Limits

### Як працюють requests і limits

```
Pod Definition:
resources:
  requests:          ← Для scheduling
    cpu: 100m
    memory: 128Mi
  limits:           ← Maximum
    cpu: 200m
    memory: 256Mi

Node with 4 CPUs, 8Gi RAM:

┌───────────────────────────────────┐
│           Node Resources           │
│  4000m CPU          8Gi RAM        │
├───────────────────────────────────┤
│                                    │
│  Pod1: requests 100m/128Mi         │
│  ├─────────┐                       │
│  │ ░░░░░   │ ← Can use up to       │
│  └─────────┘    200m/256Mi         │
│                                    │
│  Pod2: requests 100m/128Mi         │
│  ├─────────┐                       │
│  │ ░░░░░   │                       │
│  └─────────┘                       │
│                                    │
│  ... more pods ...                 │
│                                    │
│  Available: 3800m/7.7Gi            │
│  Reserved: 200m/256Mi              │
└───────────────────────────────────┘

If Pod exceeds:
- Memory limit → OOMKilled (killed)
- CPU limit → Throttled (slowed down)
```

## 🎨 Complete Flow: Development → Production

```
Developer
    │
    │ 1. Writes code
    ▼
┌──────────────┐
│  Git Push    │
└──────┬───────┘
       │
       │ 2. CI/CD Pipeline
       ▼
┌──────────────┐
│ Docker Build │
│ backend:v2   │
└──────┬───────┘
       │
       │ 3. Push to Registry
       ▼
┌──────────────────┐
│ Container        │
│ Registry         │
└──────┬───────────┘
       │
       │ 4. Update Deployment
       │ kubectl set image ...
       ▼
┌────────────────────────────┐
│   Kubernetes Cluster       │
│                            │
│  Deployment Controller     │
│  ┌──────────────────────┐ │
│  │ Rolling Update       │ │
│  │ v1 → v2              │ │
│  └──────────────────────┘ │
│                            │
│  ┌────┐ ┌────┐ ┌────┐    │
│  │ v1 │→│ v2 │→│ v2 │→   │
│  └────┘ └────┘ └────┘    │
│                            │
└────────────────────────────┘
```

---

**Ці діаграми допоможуть візуалізувати як працює Kubernetes!** 📊
