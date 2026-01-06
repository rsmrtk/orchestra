package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/rsmrtk/orchestra/backend/internal/rest"
	"github.com/rsmrtk/orchestra/backend/internal/rest/services"
)

// App головна структура додатку
type App struct {
	db           *pgxpool.Pool
	restServices *services.Services
	restServer   *rest.Server
}

// Run запускає додаток
func Run() {
	// Завантажуємо .env файл (ігноруємо помилку, якщо файл не існує)
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	ctx := context.Background()

	app := &App{}

	// Ініціалізуємо БД
	if err := app.initDB(ctx); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer app.db.Close()

	// Ініціалізуємо сервіси
	app.restServices = services.NewServices(app.db)

	// Ініціалізуємо REST сервер
	port := os.Getenv("PORT")
	if port == "" {
		port = "8282"
	}

	server, err := rest.NewServer(rest.ServerOptions{
		Services: app.restServices,
		Port:     port,
	})
	if err != nil {
		log.Fatalf("Failed to create REST server: %v", err)
	}
	app.restServer = server

	// Запускаємо сервер
	log.Printf("Starting Orchestra backend on port %s", port)
	if err := app.restServer.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initDB ініціалізує підключення до БД через pgxpool
func (a *App) initDB(ctx context.Context) error {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)

	log.Printf("Connecting to database: host=%s port=%s dbname=%s user=%s", dbHost, dbPort, dbName, dbUser)

	// Парсимо конфігурацію
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Printf("ERROR: Failed to parse database config: %v", err)
		return fmt.Errorf("failed to parse database config: %w", err)
	}
	log.Println("Database config parsed successfully")

	// Налаштовуємо connection pool
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 5 * time.Minute
	config.MaxConnIdleTime = 90 * time.Second

	// Створюємо pool з timeout
	log.Println("Creating database connection pool...")
	connectCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	a.db, err = pgxpool.NewWithConfig(connectCtx, config)
	if err != nil {
		log.Printf("ERROR: Failed to create connection pool: %v", err)
		return fmt.Errorf("failed to create connection pool: %w", err)
	}
	log.Println("Connection pool created successfully")

	// Перевіряємо з'єднання з timeout
	log.Println("Pinging database...")
	pingCtx, pingCancel := context.WithTimeout(ctx, 10*time.Second)
	defer pingCancel()

	if err := a.db.Ping(pingCtx); err != nil {
		log.Printf("ERROR: Failed to ping database: %v", err)
		return fmt.Errorf("failed to ping database: %w", err)
	}
	log.Println("Database ping successful!")

	log.Println("Successfully connected to PostgreSQL!")
	return nil
}
