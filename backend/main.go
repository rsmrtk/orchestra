package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

// initDB ініціалізує підключення до PostgreSQL використовуючи змінні середовища
func initDB() error {
	// Читаємо секрети з Kubernetes (змінні середовища)
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Встановлюємо значення за замовчуванням якщо змінні не задані
	//if dbHost == "" {
	//	dbHost = "localhost"
	//}
	//if dbPort == "" {
	//	dbPort = "5432"
	//}
	//if dbUser == "" {
	//	dbUser = "postgres"
	//}
	//if dbName == "" {
	//	dbName = "orchestra"
	//}

	// Формуємо connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	log.Printf("Підключаюсь до БД: host=%s port=%s dbname=%s user=%s", dbHost, dbPort, dbName, dbUser)

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("помилка відкриття з'єднання: %w", err)
	}

	// Налаштовуємо connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Перевіряємо з'єднання
	if err := db.Ping(); err != nil {
		return fmt.Errorf("помилка ping до БД: %w", err)
	}

	log.Println("Успішно підключено до PostgreSQL!")
	return nil
}

func main() {
	// Ініціалізуємо підключення до БД
	if err := initDB(); err != nil {
		log.Fatalf("Не вдалося підключитися до БД: %v", err)
	}
	defer db.Close()

	// Створюємо Gin router
	router := gin.Default()

	// Налаштовуємо CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:5173"} // React dev servers
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	router.Use(cors.New(config))

	// GET endpoint який повертає "Congratulation"
	router.GET("/congratulation", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "CongratulationV2",
		})
	})

	// Endpoint для перевірки підключення до БД
	router.GET("/health/db", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "error",
				"error":  "База даних недоступна",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Підключення до БД успішне",
		})
	})

	// Endpoint для виконання простого запиту до БД
	router.GET("/db/version", func(c *gin.Context) {
		var version string
		err := db.QueryRow("SELECT version()").Scan(&version)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  fmt.Sprintf("Помилка запиту: %v", err),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": version,
		})
	})

	// Endpoint для отримання поточного часу з БД
	router.GET("/db/time", func(c *gin.Context) {
		var currentTime time.Time
		err := db.QueryRow("SELECT NOW()").Scan(&currentTime)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status": "error",
				"error":  fmt.Sprintf("Помилка запиту: %v", err),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":      "ok",
			"server_time": currentTime.Format(time.RFC3339),
		})
	})

	// Запускаємо сервер на порту 8282
	log.Println("Сервер запущено на порту 8282")
	router.Run(":8282")
}
