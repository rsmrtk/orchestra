package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/rsmrtk/orchestra/backend/internal/rest"
	"github.com/rsmrtk/orchestra/backend/internal/rest/services"
)

// App головна структура додатку
type App struct {
	db           *sql.DB
	restServices *services.Services
	restServer   *rest.Server
}

// Run запускає додаток
func Run() {
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

// initDB ініціалізує підключення до БД
func (a *App) initDB(ctx context.Context) error {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	log.Printf("Connecting to database: host=%s port=%s dbname=%s user=%s", dbHost, dbPort, dbName, dbUser)

	var err error
	a.db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	// Налаштовуємо connection pool
	a.db.SetMaxOpenConns(25)
	a.db.SetMaxIdleConns(5)
	a.db.SetConnMaxLifetime(5 * time.Minute)

	// Перевіряємо з'єднання
	if err := a.db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL!")
	return nil
}
