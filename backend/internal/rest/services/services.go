package services

import (
	"database/sql"

	"github.com/rsmrtk/orchestra/backend/internal/rest/services/auth"
)

// Services містить всі сервіси додатку
type Services struct {
	Auth *auth.Service
}

// NewServices створює новий екземпляр Services
func NewServices(db *sql.DB) *Services {
	return &Services{
		Auth: auth.NewService(db),
	}
}
