package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rsmrtk/orchestra/backend/internal/rest/controllers"
	"github.com/rsmrtk/orchestra/backend/internal/rest/middlewares"
	"github.com/rsmrtk/orchestra/backend/internal/rest/services"
)

// Server REST сервер
type Server struct {
	addr   string
	engine *gin.Engine
}

// ServerOptions опції для створення сервера
type ServerOptions struct {
	Services *services.Services
	Port     string
}

// NewServer створює новий REST сервер
func NewServer(opts ServerOptions) (*Server, error) {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(middlewares.CORSMiddleware())

	// Health check endpoint
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Auth endpoints
	authController := controllers.NewAuthController(opts.Services.Auth)
	auth := engine.Group("/auth")
	{
		auth.POST("/login", authController.Login)
		auth.POST("/register", authController.Register)
	}

	// Congratulation endpoint для фронтенду
	engine.GET("/congratulation", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Congratulation",
		})
	})

	addr := fmt.Sprintf(":%s", opts.Port)
	return &Server{
		addr:   addr,
		engine: engine,
	}, nil
}

// Start запускає сервер
func (s *Server) Start() error {
	return s.engine.Run(s.addr)
}
