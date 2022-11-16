// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: notification.sql

package db

import (
	"context"
	"fmt"
    qc "github.com/khorsl/teacherportal/db/querycustom"
)

type GetStudentsEmailForNotificationParams struct {
    TeacherID  int64  `json:"teacher_id"`
    StudentEmails string `json:"student_emails"`
}

func (q *Queries) GetStudentsEmailForNotification(ctx context.Context, arg GetStudentsEmailForNotificationParams) ([]string, error) {
    // getStudentsEmailForNotification := fmt.Sprintf(
    //     `-- name: GetStudentsEmailForNotification :many
    //     SELECT s.email FROM student s WHERE
    //     s.id IN (SELECT r.student_id FROM register r WHERE r.teacher_id = %d)
    //     OR s.email IN (%s)
    //     AND is_suspended = false`, arg.TeacherID, arg.StudentEmails,
    // )
    getStudentsEmailForNotification := fmt.Sprintf(
        qc.GetStudentsEmailForNotificationSQL, arg.TeacherID, arg.StudentEmails,
    )
    rows, err := q.db.QueryContext(ctx, getStudentsEmailForNotification)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    items := []string{}
    for rows.Next() {
        var email string
        if err := rows.Scan(&email); err != nil {
            return nil, err
        }
        items = append(items, email)
    }
    if err := rows.Close(); err != nil {
        return nil, err
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return items, nil
}
