package customer

import (
	"context"

	"github.com/rsmrtk/orchestra/backend/pkg/models"
)

// RegisterRequest запит для реєстрації
type RegisterRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// RegisterResponse відповідь після реєстрації
type RegisterResponse struct {
	Success  bool   `json:"success"`
	Message  string `json:"message"`
	Redirect string `json:"redirect,omitempty"`
}

// Register виконує реєстрацію нового клієнта
func Register(ctx context.Context, repo models.CustomerRepository, req RegisterRequest) (*RegisterResponse, error) {
	// Перевіряємо чи вже існує клієнт з таким ім'ям та прізвищем
	exists, err := repo.Exists(ctx, req.FirstName, req.LastName)
	if err != nil {
		return nil, err
	}

	// Якщо клієнт вже існує - перенаправляємо на логін
	if exists {
		return &RegisterResponse{
			Success:  false,
			Message:  "Customer already exists. Please login.",
			Redirect: "/login",
		}, nil
	}

	// Створюємо нового клієнта
	customer := &models.Customer{
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	err = repo.Create(ctx, customer)
	if err != nil {
		return nil, err
	}

	// Успішна реєстрація - перенаправляємо на логін
	return &RegisterResponse{
		Success:  true,
		Message:  "Registration successful. Please login.",
		Redirect: "/login",
	}, nil
}
