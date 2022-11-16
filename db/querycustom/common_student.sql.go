package querycustom

const GetCommonStudentsEmailSQL string = `-- name: GetCommonStudentsEmail :many
SELECT s.email FROM student s WHERE s.id IN (
SELECT r.student_id FROM register r WHERE r.teacher_id
IN (SELECT t.id FROM teacher t WHERE t.email IN (%s))
GROUP BY r.student_id
HAVING COUNT(r.student_id) = %d)`
