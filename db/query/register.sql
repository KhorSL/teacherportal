-- name: CreateRegister :exec
INSERT INTO register (
  teacher_id,
  student_id
) VALUES (
  ?, ?
);

-- name: GetRegisterByTeacherId :many
SELECT * FROM register
WHERE teacher_id = ?;

-- name: GetRegisterByStudentId :many
SELECT * FROM register
WHERE student_id = ?;

-- name: GetRegisterByStudentIdAndTeacherId :one
SELECT * FROM register
WHERE student_id = ? AND teacher_id = ?;