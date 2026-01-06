package models

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rsmrtk/orchestra/backend/pkg/dbq"
)

// Customer представляє клієнта в системі
type Customer struct {
	CustomerID string `json:"customer_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
}

// CustomerRepository інтерфейс для роботи з клієнтами
type CustomerRepository interface {
	GetByName(ctx context.Context, firstName, lastName string) (*Customer, error)
	Create(ctx context.Context, customer *Customer) error
	Exists(ctx context.Context, firstName, lastName string) (bool, error)
}

// customerRepo реалізація CustomerRepository з використанням sqlc
type customerRepo struct {
	db      *pgxpool.Pool
	queries *dbq.Queries
}

// NewCustomerRepository створює новий репозиторій для Customer
func NewCustomerRepository(db *pgxpool.Pool) CustomerRepository {
	return &customerRepo{
		db:      db,
		queries: dbq.New(db),
	}
}

// GetByName повертає клієнта за ім'ям та прізвищем
func (r *customerRepo) GetByName(ctx context.Context, firstName, lastName string) (*Customer, error) {
	result, err := r.queries.GetCustomerByName(ctx, dbq.GetCustomerByNameParams{
		FirstName: firstName,
		LastName:  lastName,
	})

	if err == pgx.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &Customer{
		CustomerID: result.CustomerID,
		FirstName:  result.FirstName,
		LastName:   result.LastName,
	}, nil
}

// Create створює нового клієнта в БД
func (r *customerRepo) Create(ctx context.Context, customer *Customer) error {
	// Генеруємо UUID якщо не заданий
	if customer.CustomerID == "" {
		customer.CustomerID = uuid.New().String()
	}

	result, err := r.queries.CreateCustomer(ctx, dbq.CreateCustomerParams{
		CustomerID: customer.CustomerID,
		FirstName:  customer.FirstName,
		LastName:   customer.LastName,
	})

	if err != nil {
		return err
	}

	// Оновлюємо customer з результатом
	customer.CustomerID = result.CustomerID

	return nil
}

// Exists перевіряє чи існує клієнт з таким ім'ям та прізвищем
func (r *customerRepo) Exists(ctx context.Context, firstName, lastName string) (bool, error) {
	return r.queries.CustomerExists(ctx, dbq.CustomerExistsParams{
		FirstName: firstName,
		LastName:  lastName,
	})
}

// this changes for fun
