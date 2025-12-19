-- name: GetSchool :one
SELECT * FROM schools
WHERE school_id = $1 LIMIT 1;

-- name: GetSchoolByName :one
SELECT * FROM schools
WHERE school_name = $1 LIMIT 1;

-- name: ListSchools :many
SELECT * FROM schools
ORDER BY school_name;

-- name: CreateSchool :one
INSERT INTO schools (school_name)
VALUES ($1)
RETURNING *;

-- name: UpdateSchool :one
UPDATE schools
SET school_name = $2, updated_at = NOW()
WHERE school_id = $1
RETURNING *;

-- name: DeleteSchool :exec
DELETE FROM schools
WHERE school_id = $1;
