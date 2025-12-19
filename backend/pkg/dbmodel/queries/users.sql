-- ============================================
-- CRUD Operations for users table
-- ============================================

-- name: UserCreate :one
INSERT INTO users (first_name, last_name, email, phone, status, role, school, password_hash)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING user_id, first_name, last_name, email, phone, status, role, school, created_at, updated_at;

-- name: UserGetByID :one
SELECT user_id, first_name, last_name, email, phone, status, role, school, created_at, updated_at
FROM users
WHERE user_id = $1;

-- name: UserGetByEmail :one
SELECT user_id, first_name, last_name, email, phone, status, role, school, created_at, updated_at
FROM users
WHERE email = $1;

-- name: UserList :many
SELECT user_id, first_name, last_name, email, phone, status, role, school, created_at, updated_at
FROM users
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UserUpdate :one
UPDATE users
SET first_name    = $2,
    last_name     = $3,
    email         = $4,
    phone         = $5,
    status        = $6,
    role          = $7,
    school        = $8,
    updated_at    = NOW()
WHERE user_id = $1
RETURNING user_id, first_name, last_name, email, phone, status, role, school, created_at, updated_at;

-- name: UserDelete :exec
DELETE FROM users
WHERE user_id = $1;

-- ============================================
-- Additional Queries
-- ============================================

-- name: UserListByStatus :many
SELECT user_id, first_name, last_name, email, phone, status, role, school, created_at, updated_at
FROM users
WHERE status = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UserListByRole :many
SELECT user_id, first_name, last_name, email, phone, status, role, school, created_at, updated_at
FROM users
WHERE role = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UserListBySchool :many
SELECT user_id, first_name, last_name, email, phone, status, role, school, created_at, updated_at
FROM users
WHERE school = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UserSearch :many
SELECT user_id, first_name, last_name, email, phone, status, role, school, created_at, updated_at
FROM users
WHERE
    first_name ILIKE '%' || $1 || '%' OR
    last_name ILIKE '%' || $1 || '%' OR
    email ILIKE '%' || $1 || '%'
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- ============================================
-- Authentication
-- ============================================

-- name: UserGetPasswordHash :one
SELECT user_id, password_hash
FROM users
WHERE email = $1;

-- name: UserUpdatePassword :exec
UPDATE users
SET password_hash = $2,
    updated_at = NOW()
WHERE user_id = $1;

-- ============================================
-- Statistics
-- ============================================

-- name: UserCountByStatus :one
SELECT COUNT(*) as count
FROM users
WHERE status = $1;

-- name: UserCountByRole :one
SELECT COUNT(*) as count
FROM users
WHERE role = $1;

-- name: UserCountTotal :one
SELECT COUNT(*) as count
FROM users;
