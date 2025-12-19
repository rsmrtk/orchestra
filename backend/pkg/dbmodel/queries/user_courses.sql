-- ============================================
-- CRUD Operations for user_courses table
-- ============================================

-- name: UserCourseAssign :one
INSERT INTO user_courses (user_id, course_id, expired_date, status)
VALUES ($1, $2, $3, $4)
RETURNING user_id, course_id, purchase_date, expired_date, status, created_at, updated_at;

-- name: UserCourseGetByUserAndCourse :one
SELECT user_id, course_id, purchase_date, expired_date, status, created_at, updated_at
FROM user_courses
WHERE user_id = $1 AND course_id = $2;

-- name: UserCourseGetByUserID :many
SELECT user_id, course_id, purchase_date, expired_date, status, created_at, updated_at
FROM user_courses
WHERE user_id = $1
ORDER BY purchase_date DESC;

-- name: UserCourseGetByCourseID :many
SELECT user_id, course_id, purchase_date, expired_date, status, created_at, updated_at
FROM user_courses
WHERE course_id = $1
ORDER BY purchase_date DESC
LIMIT $2 OFFSET $3;

-- name: UserCourseList :many
SELECT user_id, course_id, purchase_date, expired_date, status, created_at, updated_at
FROM user_courses
ORDER BY purchase_date DESC
LIMIT $1 OFFSET $2;

-- name: UserCourseUpdate :one
UPDATE user_courses
SET expired_date = $3,
    status       = $4,
    updated_at   = NOW()
WHERE user_id = $1 AND course_id = $2
RETURNING user_id, course_id, purchase_date, expired_date, status, created_at, updated_at;

-- name: UserCourseUpdateStatus :one
UPDATE user_courses
SET status     = $3,
    updated_at = NOW()
WHERE user_id = $1 AND course_id = $2
RETURNING user_id, course_id, purchase_date, expired_date, status, created_at, updated_at;

-- name: UserCourseRemove :exec
DELETE FROM user_courses
WHERE user_id = $1 AND course_id = $2;

-- name: UserCourseRemoveAllByUser :exec
DELETE FROM user_courses
WHERE user_id = $1;

-- name: UserCourseRemoveAllByCourse :exec
DELETE FROM user_courses
WHERE course_id = $1;

-- ============================================
-- Filtered Queries
-- ============================================

-- name: UserCourseGetActiveByUserID :many
SELECT user_id, course_id, purchase_date, expired_date, status, created_at, updated_at
FROM user_courses
WHERE user_id = $1 AND status = 'valid'
ORDER BY purchase_date DESC;

-- name: UserCourseGetExpiredByUserID :many
SELECT user_id, course_id, purchase_date, expired_date, status, created_at, updated_at
FROM user_courses
WHERE user_id = $1 AND status = 'expired'
ORDER BY purchase_date DESC;

-- name: UserCourseGetByStatus :many
SELECT user_id, course_id, purchase_date, expired_date, status, created_at, updated_at
FROM user_courses
WHERE status = $1
ORDER BY purchase_date DESC
LIMIT $2 OFFSET $3;

-- name: UserCourseGetExpiringBefore :many
SELECT user_id, course_id, purchase_date, expired_date, status, created_at, updated_at
FROM user_courses
WHERE expired_date IS NOT NULL
  AND expired_date <= $1
  AND status = 'valid'
ORDER BY expired_date ASC;

-- ============================================
-- Statistics
-- ============================================

-- name: UserCourseCountByUserID :one
SELECT COUNT(*) as count
FROM user_courses
WHERE user_id = $1;

-- name: UserCourseCountByCourseID :one
SELECT COUNT(*) as count
FROM user_courses
WHERE course_id = $1;

-- name: UserCourseCountByStatus :one
SELECT COUNT(*) as count
FROM user_courses
WHERE status = $1;

-- ============================================
-- Join Queries with users and courses tables
-- ============================================

-- name: UserCourseGetUsersWithDetails :many
SELECT
    u.user_id,
    u.first_name,
    u.last_name,
    u.email,
    u.phone,
    u.status as user_status,
    u.role,
    c.course_id,
    c.course_name,
    c.school_id,
    uc.purchase_date,
    uc.expired_date,
    uc.status as course_status
FROM user_courses uc
JOIN users u ON uc.user_id = u.user_id
JOIN courses c ON uc.course_id = c.course_id
WHERE uc.course_id = $1
ORDER BY uc.purchase_date DESC
LIMIT $2 OFFSET $3;

-- name: UserCourseGetUserCoursesWithDetails :many
SELECT
    uc.user_id,
    c.course_id,
    c.course_name,
    c.school_id,
    s.school_name,
    uc.purchase_date,
    uc.expired_date,
    uc.status as course_status,
    u.first_name,
    u.last_name,
    u.email
FROM user_courses uc
JOIN users u ON uc.user_id = u.user_id
JOIN courses c ON uc.course_id = c.course_id
JOIN schools s ON c.school_id = s.school_id
WHERE uc.user_id = $1
ORDER BY uc.purchase_date DESC;

-- name: UserCourseGetBySchoolID :many
SELECT
    uc.user_id,
    uc.course_id,
    uc.purchase_date,
    uc.expired_date,
    uc.status,
    uc.created_at,
    uc.updated_at
FROM user_courses uc
JOIN courses c ON uc.course_id = c.course_id
WHERE c.school_id = $1
ORDER BY uc.purchase_date DESC
LIMIT $2 OFFSET $3;

-- name: UserCourseCountBySchoolID :one
SELECT COUNT(*) as count
FROM user_courses uc
JOIN courses c ON uc.course_id = c.course_id
WHERE c.school_id = $1;
