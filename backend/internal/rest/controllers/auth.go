package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rsmrtk/orchestra/backend/internal/rest/domain/customer"
	"github.com/rsmrtk/orchestra/backend/internal/rest/services/auth"
)

// AuthController контролер для аутентифікації
type AuthController struct {
	authService *auth.Service
}

// NewAuthController створює новий AuthController
func NewAuthController(authService *auth.Service) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// Login обробляє запит на логін
func (ctrl *AuthController) Login(c *gin.Context) {
	var req customer.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid login request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
		})
		return
	}

	// Викликаємо доменну логіку
	resp, err := customer.Login(c.Request.Context(), ctrl.authService.GetCustomerRepo(), req)
	if err != nil {
		log.Printf("Login error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal server error",
		})
		return
	}

	// Повертаємо результат
	if resp.Success {
		c.JSON(http.StatusOK, resp)
	} else {
		c.JSON(http.StatusUnauthorized, resp)
	}
}

// Register обробляє запит на реєстрацію
func (ctrl *AuthController) Register(c *gin.Context) {
	var req customer.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Invalid register request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request data",
		})
		return
	}

	// Викликаємо доменну логіку
	resp, err := customer.Register(c.Request.Context(), ctrl.authService.GetCustomerRepo(), req)
	if err != nil {
		log.Printf("Register error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal server error",
		})
		return
	}

	// Повертаємо результат
	if resp.Success {
		c.JSON(http.StatusCreated, resp)
	} else {
		c.JSON(http.StatusConflict, resp)
	}
}
