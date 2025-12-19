package services

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rsmrtk/orchestra/backend/internal/rest/services/auth"
)

// Services містить всі сервіси додатку
type Services struct {
	Auth *auth.Service
}

// NewServices створює новий екземпляр Services
func NewServices(db *pgxpool.Pool) *Services {
	return &Services{
		Auth: auth.NewService(db),
	}
}
