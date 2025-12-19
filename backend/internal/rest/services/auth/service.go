package auth

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rsmrtk/orchestra/backend/pkg/models"
)

// Service сервіс для аутентифікації
type Service struct {
	customerRepo models.CustomerRepository
}

// NewService створює новий auth сервіс
func NewService(db *pgxpool.Pool) *Service {
	return &Service{
		customerRepo: models.NewCustomerRepository(db),
	}
}

// GetCustomerRepo повертає репозиторій клієнтів
func (s *Service) GetCustomerRepo() models.CustomerRepository {
	return s.customerRepo
}
