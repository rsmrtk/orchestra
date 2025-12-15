# Orchestra Deploy Practice

Простий проект для практики деплойменту з Go бекендом та React фронтендом.

## Структура проекту

```
orchestra/
├── backend/     # Go сервер з Gin
└── frontend/    # React додаток
```

## Backend (Go + Gin)

### Запуск backend

```bash
cd backend
go run main.go
```

Сервер запуститься на `http://localhost:8080`

### Endpoints

- `GET /congratulation` - повертає JSON з повідомленням "Congratulation"

### Компіляція

```bash
cd backend
go build -o server
./server
```

## Frontend (React + Vite)

### Запуск frontend

```bash
cd frontend
npm run dev
```

Додаток запуститься на `http://localhost:5173`

### Збірка для production

```bash
cd frontend
npm run build
```

## Як використовувати

1. Запустіть backend сервер:
   ```bash
   cd backend
   go run main.go
   ```

2. В іншому терміналі запустіть frontend:
   ```bash
   cd frontend
   npm run dev
   ```

3. Відкрийте браузер за адресою `http://localhost:5173`

4. Ви побачите повідомлення "Congratulation" отримане з API

## Технології

- **Backend**: Go 1.25, Gin, CORS middleware
- **Frontend**: React, Vite, JavaScript

## Особливості

- CORS налаштований для локальної розробки
- Автоматичне завантаження повідомлення при відкритті сторінки
- Кнопка для повторного завантаження даних
- Обробка помилок та стану завантаження
