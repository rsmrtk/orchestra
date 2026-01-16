# Orchestra Authentication System

Простий authentication система з бекендом на Go та фронтендом на React.

## Архітектура

### Backend (Go + Gin + PostgreSQL)

```
backend/
├── cmd/server/main.go           # Entry point
├── internal/
│   ├── app/app.go              # Application initialization
│   └── rest/
│       ├── controllers/        # HTTP handlers
│       ├── domain/            # Domain logic
│       │   └── customer/      # Customer business logic
│       ├── middlewares/       # Middlewares (CORS, etc)
│       ├── services/          # Services layer
│       └── server.go          # REST server setup
└── pkg/
    └── models/                # Database models & repositories
```

### Frontend (React + Vite)

```
frontend/
├── src/
│   ├── pages/
│   │   ├── Login.jsx         # Login page
│   │   ├── Register.jsx      # Registration page
│   │   ├── Congratulation.jsx # Success page
│   │   └── Auth.css          # Styles
│   ├── App.jsx              # Main app with routing
│   └── main.jsx             # Entry point
```

## Features

- **Login**: Перевіряє чи існує користувач з ім'ям та прізвищем
- **Register**: Реєструє нового користувача
- **Congratulation**: Сторінка успішного входу

## База даних

PostgreSQL з таблицею:

```sql
CREATE TABLE "orchestra-table" (
    customer_id VARCHAR(255) NOT NULL PRIMARY KEY,
    first_name  VARCHAR(255) NOT NULL,
    last_name   VARCHAR(255) NOT NULL
);
```

## Запуск проєкту

### Backend

1. Встановіть змінні середовища:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=your_db_username
export DB_PASSWORD=your_secure_password
export DB_NAME=orchestra
export PORT=8282
```

2. Запустіть бекенд:

```bash
cd backend
go run cmd/server/main.go
```

Сервер запуститься на `http://localhost:8282`

### Frontend

1. Створіть `.env` файл:

```bash
cd frontend
cp .env.example .env
```

2. Встановіть залежності та запустіть:

```bash
npm install
npm run dev
```

Фронтенд запуститься на `http://localhost:5173`

## API Endpoints

### POST /auth/login
Логін користувача

**Request:**
```json
{
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response (Success):**
```json
{
  "success": true,
  "message": "Login successful",
  "redirect": "/congratulation",
  "customer": {
    "customer_id": "uuid",
    "first_name": "John",
    "last_name": "Doe"
  }
}
```

**Response (Not Found):**
```json
{
  "success": false,
  "message": "You are not registered"
}
```

### POST /auth/register
Реєстрація користувача

**Request:**
```json
{
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response (Success):**
```json
{
  "success": true,
  "message": "Registration successful. Please login.",
  "redirect": "/login"
}
```

**Response (Already Exists):**
```json
{
  "success": false,
  "message": "Customer already exists. Please login.",
  "redirect": "/login"
}
```

### GET /health
Health check endpoint

### GET /congratulation
Повертає congratulation message

## Технології

### Backend
- Go 1.25+
- Gin web framework
- PostgreSQL
- Clean Architecture (Repository pattern, Service layer)

### Frontend
- React 19
- React Router DOM
- Vite
- CSS3 (з градієнтами та анімаціями)

## Структура коду

### Backend Layers:

1. **Models** (`pkg/models/`) - Структури даних та repository interfaces
2. **Domain** (`internal/rest/domain/`) - Бізнес логіка
3. **Services** (`internal/rest/services/`) - Service layer
4. **Controllers** (`internal/rest/controllers/`) - HTTP handlers
5. **Middlewares** (`internal/rest/middlewares/`) - CORS, error handling

### Frontend Pages:

1. **Login** - Форма логіну з перевіркою існування користувача
2. **Register** - Форма реєстрації з перевіркою дублікатів
3. **Congratulation** - Сторінка успішного входу

## Development

### Backend development:
```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

### Frontend development:
```bash
cd frontend
npm run dev
```

### Build for production:

Backend:
```bash
cd backend
go build -o bin/server cmd/server/main.go
```

Frontend:
```bash
cd frontend
npm run build
```
