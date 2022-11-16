package querycustom

const GetStudentsEmailForNotificationSQL string = `-- name: GetStudentsEmailForNotification :many
SELECT s.email FROM student s WHERE
(s.id IN (SELECT r.student_id FROM register r WHERE r.teacher_id = %d)
OR s.email IN (%s))
AND is_suspended = false;`
