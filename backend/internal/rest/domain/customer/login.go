package customer

import (
	"context"

	"github.com/rsmrtk/orchestra/backend/pkg/models"
)

// LoginRequest запит для логіну
type LoginRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// LoginResponse відповідь після логіну
type LoginResponse struct {
	Success  bool             `json:"success"`
	Message  string           `json:"message"`
	Redirect string           `json:"redirect,omitempty"`
	Customer *models.Customer `json:"customer,omitempty"`
}

// Login виконує логін клієнта
func Login(ctx context.Context, repo models.CustomerRepository, req LoginRequest) (*LoginResponse, error) {
	// Перевіряємо чи існує клієнт
	customer, err := repo.GetByName(ctx, req.FirstName, req.LastName)
	if err != nil {
		return nil, err
	}

	// Якщо клієнт не знайдений
	if customer == nil {
		return &LoginResponse{
			Success: false,
			Message: "You are not registered",
		}, nil
	}

	// Клієнт знайдений - успішний логін
	return &LoginResponse{
		Success:  true,
		Message:  "Login successful",
		Redirect: "/congratulation",
		Customer: customer,
	}, nil
}
