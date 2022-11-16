-- name: CreateTeacher :exec
INSERT INTO teacher (
  full_name,
  email
) VALUES (
  ?, ?
);

-- name: GetTeacherById :one
SELECT * FROM teacher
WHERE id = ? AND is_active = 1 LIMIT 1;

-- name: GetTeacherByEmail :one
SELECT * FROM teacher
WHERE email = ? AND is_active = 1 LIMIT 1;
