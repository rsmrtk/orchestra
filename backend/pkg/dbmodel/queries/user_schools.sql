-- ============================================
-- CRUD Operations for user_schools table
-- ============================================

-- name: UserSchoolAdd :one
INSERT INTO user_schools (user_id, school_id)
VALUES ($1, $2)
RETURNING user_id, school_id, created_at, updated_at;

-- name: UserSchoolRemove :exec
DELETE FROM user_schools
WHERE user_id = $1 AND school_id = $2;

-- name: UserSchoolGetByUserID :many
SELECT user_id, school_id, created_at, updated_at
FROM user_schools
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UserSchoolGetBySchoolID :many
SELECT user_id, school_id, created_at, updated_at
FROM user_schools
WHERE school_id = $1
ORDER BY created_at DESC;

-- name: UserSchoolList :many
SELECT user_id, school_id, created_at, updated_at
FROM user_schools
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: UserSchoolRemoveAllByUser :exec
DELETE FROM user_schools
WHERE user_id = $1;

-- name: UserSchoolRemoveAllBySchool :exec
DELETE FROM user_schools
WHERE school_id = $1;

-- ============================================
-- Additional Queries
-- ============================================

-- name: UserSchoolCheck :one
SELECT EXISTS(
    SELECT 1
    FROM user_schools
    WHERE user_id = $1 AND school_id = $2
) as exists;

-- name: UserSchoolCountBySchoolID :one
SELECT COUNT(*) as count
FROM user_schools
WHERE school_id = $1;

-- name: UserSchoolCountByUserID :one
SELECT COUNT(*) as count
FROM user_schools
WHERE user_id = $1;

-- ============================================
-- Join Queries with users and schools tables
-- ============================================

-- name: UserSchoolGetUsersWithDetails :many
SELECT
    u.user_id,
    u.first_name,
    u.last_name,
    u.email,
    u.phone,
    u.status,
    u.role,
    s.school_id,
    s.school_name,
    us.created_at as enrolled_at
FROM user_schools us
JOIN users u ON us.user_id = u.user_id
JOIN schools s ON us.school_id = s.school_id
WHERE us.school_id = $1
ORDER BY us.created_at DESC
LIMIT $2 OFFSET $3;

-- name: UserSchoolGetWithSchoolName :many
SELECT
    us.user_id,
    us.school_id,
    s.school_name,
    us.created_at,
    us.updated_at
FROM user_schools us
JOIN schools s ON us.school_id = s.school_id
WHERE us.user_id = $1
ORDER BY us.created_at DESC;
