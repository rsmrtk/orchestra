-- name: GetCourse :one
SELECT * FROM courses
WHERE course_id = $1 LIMIT 1;

-- name: ListCourses :many
SELECT * FROM courses
ORDER BY course_name;

-- name: ListCoursesBySchool :many
SELECT * FROM courses
WHERE school_id = $1
ORDER BY course_name;

-- name: CreateCourse :one
INSERT INTO courses (school_id, course_name)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateCourse :one
UPDATE courses
SET course_name = $2, updated_at = NOW()
WHERE course_id = $1
RETURNING *;

-- name: DeleteCourse :exec
DELETE FROM courses
WHERE course_id = $1;
