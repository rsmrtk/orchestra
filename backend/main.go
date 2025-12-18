package main

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
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
			"message": "CongratulationV1",
		})
	})

	// Запускаємо сервер на порту 8080
	router.Run(":8383")
}
