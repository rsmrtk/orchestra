# Orchestra - Quick Start Guide 🚀

Простий гайд для швидкого запуску всіх компонентів разом.

## Що входить

- **Frontend**: React + Vite (порт 5252)
- **Backend**: Go + Gin (порт 8383)
- **Database**: PostgreSQL 17 (порт 2828)

## Швидкий старт

### 1. Запустити все одразу

```bash
cd /home/r/fff/deploy/orchestra
./start.sh
```

АБО

```bash
cd /home/r/fff/deploy/orchestra
docker compose -f backend/deployments/production/docker-compose.yml up -d
```

### 2. Відкрити в браузері

- Frontend: http://localhost:5252
- Backend API: http://localhost:8383/health

### 3. Переглянути логи

```bash
./logs.sh
```

АБО

```bash
docker compose -f backend/deployments/production/docker-compose.yml logs -f
```

### 4. Зупинити всі сервіси

```bash
./stop.sh
```

АБО

```bash
docker compose -f backend/deployments/production/docker-compose.yml down
```

## Конфігурація

Всі налаштування в файлі `.env`:

```bash
cd backend/deployments/production
nano .env
```

### Доступні параметри:

```env
# Database
DB_NAME=orchestra-db
DB_USER=rsmrtk
DB_PASSWORD=rsmrtk
DB_PORT=2828

# Backend
BACKEND_PORT=8383
GIN_MODE=release

# Frontend
FRONTEND_PORT=5252

# JWT
JWT_SECRET=your-secret-here
JWT_DURATION=24h
```

### Як змінити порти:

1. Відредагуйте `.env` файл
2. Перезапустіть сервіси:

```bash
cd backend/deployments/production
docker compose down
docker compose up -d
```

## Корисні команди

### Перевірити статус контейнерів

```bash
docker compose -f backend/deployments/production/docker-compose.yml ps
```

### Перезапустити конкретний сервіс

```bash
cd backend/deployments/production

# Frontend
docker compose restart frontend

# Backend
docker compose restart backend

# Database
docker compose restart postgres
```

### Пересібрати після змін коду

```bash
cd backend/deployments/production
docker compose up --build -d
```

### Підключитись до бази даних

```bash
docker exec -it orchestra-db psql -U rsmrtk -d orchestra-db
```

### Очистити все (включно з даними БД)

```bash
cd backend/deployments/production
docker compose down -v
```

⚠️ **Увага**: Команда з `-v` видалить ВСІ дані з бази даних!

## Тестування

### Зареєструвати нового користувача

```bash
curl -X POST http://localhost:8383/auth/register \
  -H "Content-Type: application/json" \
  -d '{"first_name":"Test","last_name":"User"}'
```

### Залогінитись

```bash
curl -X POST http://localhost:8383/auth/login \
  -H "Content-Type: application/json" \
  -d '{"first_name":"Test","last_name":"User"}'
```

## Troubleshooting

### Порт вже зайнятий

Змініть порт в `.env` файлі:

```env
BACKEND_PORT=8484  # Замість 8383
FRONTEND_PORT=5353  # Замість 5252
```

### База даних не стартує

Перевірте логи:

```bash
docker compose -f backend/deployments/production/docker-compose.yml logs postgres
```

### Frontend не може підключитись до backend

Перевірте що backend працює:

```bash
curl http://localhost:8383/health
```

Повинно повернути: `{"status":"ok"}`

## Структура проекту

```
orchestra/
├── backend/
│   └── deployments/
│       └── production/
│           ├── docker-compose.yml  # Головний файл
│           ├── .env                # Ваші налаштування
│           └── .env.example        # Приклад
├── frontend/
│   ├── Dockerfile                  # Frontend образ
│   └── nginx/
│       └── nginx.conf              # Nginx конфігурація
├── start.sh                        # Запустити все
├── stop.sh                         # Зупинити все
└── logs.sh                         # Переглянути логи
```

## Готово! 🎉

Тепер ви можете запустити весь проект однією командою! 🚀
