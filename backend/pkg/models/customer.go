package models

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
)

// Customer представляє клієнта в системі
type Customer struct {
	CustomerID string `json:"customer_id" db:"customer_id"`
	FirstName  string `json:"first_name" db:"first_name"`
	LastName   string `json:"last_name" db:"last_name"`
}

// CustomerRepository інтерфейс для роботи з клієнтами
type CustomerRepository interface {
	GetByName(ctx context.Context, firstName, lastName string) (*Customer, error)
	Create(ctx context.Context, customer *Customer) error
	Exists(ctx context.Context, firstName, lastName string) (bool, error)
}

// customerRepo реалізація CustomerRepository
type customerRepo struct {
	db *sql.DB
}

// NewCustomerRepository створює новий репозиторій для Customer
func NewCustomerRepository(db *sql.DB) CustomerRepository {
	return &customerRepo{db: db}
}

// GetByName повертає клієнта за ім'ям та прізвищем
func (r *customerRepo) GetByName(ctx context.Context, firstName, lastName string) (*Customer, error) {
	query := `SELECT customer_id, first_name, last_name FROM "orchestra-table" WHERE first_name = $1 AND last_name = $2`

	customer := &Customer{}
	err := r.db.QueryRowContext(ctx, query, firstName, lastName).Scan(
		&customer.CustomerID,
		&customer.FirstName,
		&customer.LastName,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return customer, nil
}

// Create створює нового клієнта в БД
func (r *customerRepo) Create(ctx context.Context, customer *Customer) error {
	// Генеруємо UUID якщо не заданий
	if customer.CustomerID == "" {
		customer.CustomerID = uuid.New().String()
	}

	query := `INSERT INTO "orchestra-table" (customer_id, first_name, last_name) VALUES ($1, $2, $3)`

	_, err := r.db.ExecContext(ctx, query, customer.CustomerID, customer.FirstName, customer.LastName)
	if err != nil {
		return fmt.Errorf("failed to create customer: %w", err)
	}

	return nil
}

// Exists перевіряє чи існує клієнт з таким ім'ям та прізвищем
func (r *customerRepo) Exists(ctx context.Context, firstName, lastName string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM "orchestra-table" WHERE first_name = $1 AND last_name = $2)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, firstName, lastName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if customer exists: %w", err)
	}

	return exists, nil
}
