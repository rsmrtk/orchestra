# Production Dockerfile - Multi-stage build
# Stage 1: Build
FROM golang:1.25-alpine AS builder

# Встановлюємо git для go mod
RUN apk add --no-cache git

WORKDIR /app

# Копіюємо dependency files
COPY go.mod go.sum ./backend/

# Завантажуємо залежності (кешується якщо go.mod не змінюється)
RUN go mod download
RUN go mod verify

# Копіюємо код
COPY . .

# Білдимо з оптимізацією для production
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -a -installsuffix cgo \
    -o /app/server \
    .

# Stage 2: Runtime
FROM alpine:3.19

# Додаємо ca-certificates для HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Створюємо non-root користувача
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /home/appuser

# Копіюємо бінарник
COPY --from=builder /app/server ./server

# Змінюємо власника
RUN chown -R appuser:appuser /home/appuser

# Переключаємось на non-root користувача
USER appuser

# Експонуємо порт
EXPOSE 8282

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8383/congratulation || exit 1

# Запускаємо сервер
CMD ["./server"]
