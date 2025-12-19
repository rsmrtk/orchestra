-- name: GetCustomerByName :one
SELECT customer_id, first_name, last_name, created_at
FROM "orchestra-table"
WHERE first_name = $1 AND last_name = $2
LIMIT 1;

-- name: CustomerExists :one
SELECT EXISTS(SELECT 1 FROM "orchestra-table" WHERE first_name = $1 AND last_name = $2);

-- name: CreateCustomer :one
INSERT INTO "orchestra-table" (customer_id, first_name, last_name, created_at)
VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
RETURNING customer_id, first_name, last_name, created_at;

-- name: ListCustomers :many
SELECT customer_id, first_name, last_name, created_at
FROM "orchestra-table"
ORDER BY created_at DESC;

-- name: GetCustomerByID :one
SELECT customer_id, first_name, last_name, created_at
FROM "orchestra-table"
WHERE customer_id = $1
LIMIT 1;
