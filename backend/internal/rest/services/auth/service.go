package auth

import (
	"database/sql"

	"github.com/rsmrtk/orchestra/backend/pkg/models"
)

// Service сервіс для аутентифікації
type Service struct {
	customerRepo models.CustomerRepository
}

// NewService створює новий auth сервіс
func NewService(db *sql.DB) *Service {
	return &Service{
		customerRepo: models.NewCustomerRepository(db),
	}
}

// GetCustomerRepo повертає репозиторій клієнтів
func (s *Service) GetCustomerRepo() models.CustomerRepository {
	return s.customerRepo
}
