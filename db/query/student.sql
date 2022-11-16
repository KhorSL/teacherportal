-- name: CreateStudent :exec
INSERT INTO student (
  full_name,
  email
) VALUES (
  ?, ?
);

-- name: GetStudentById :one
SELECT * FROM student
WHERE id = ? AND is_active = 1 LIMIT 1;

-- name: GetStudentByEmail :one
SELECT * FROM student
WHERE email = ? AND is_active = 1 LIMIT 1;

-- name: SuspendStudentByEmail :exec
UPDATE student
SET is_suspended = 1, suspended_at = CURRENT_TIMESTAMP(), updated_at = CURRENT_TIMESTAMP()
WHERE email = ?;

-- name: GetNotSuspendedStudentByEmail :one
SELECT * FROM student
WHERE email = ? AND is_suspended = 0 AND is_active = 1 LIMIT 1;