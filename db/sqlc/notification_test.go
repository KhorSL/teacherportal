package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetStudentsEmailForNotification(t *testing.T) {
	teacher, err := createRandomTeacher()
	require.NoError(t, err)

	studentR1, err := createRandomStudent()
	require.NoError(t, err)

	studentR2, err := createRandomStudent()
	require.NoError(t, err)

	studentNr3, err := createRandomStudent()
	require.NoError(t, err)

	studentNr4, err := createRandomStudent()
	require.NoError(t, err)

	studentNotMentioned, err := createRandomStudent()
	require.NoError(t, err)

	createRegisterArg := CreateRegisterParams{TeacherID: teacher.ID, StudentID: studentR1.ID}
	err2 := testQueries.CreateRegister(context.Background(), createRegisterArg)
	require.NoError(t, err2)

	createRegisterArg = CreateRegisterParams{TeacherID: teacher.ID, StudentID: studentR2.ID}
	err2 = testQueries.CreateRegister(context.Background(), createRegisterArg)
	require.NoError(t, err2)

	studentEmails := fmt.Sprintf("'%s','%s'", studentNr3.Email, studentNr4.Email)
	notificationArg := GetStudentsEmailForNotificationParams{
		TeacherID:     teacher.ID,
		StudentEmails: studentEmails,
	}
	emails, err3 := testQueries.GetStudentsEmailForNotification(context.Background(), notificationArg)
	require.NoError(t, err3)
	require.NotEmpty(t, emails)

	require.Len(t, emails, 4)
	require.Contains(t, emails, studentR1.Email)
	require.Contains(t, emails, studentR2.Email)
	require.Contains(t, emails, studentNr3.Email)
	require.Contains(t, emails, studentNr4.Email)
	require.NotContains(t, emails, studentNotMentioned.Email)
}

func TestGetStudentsEmailForNotification_RegisteredAndMentioned(t *testing.T) {
	teacher, err := createRandomTeacher()
	require.NoError(t, err)

	studentR1, err := createRandomStudent()
	require.NoError(t, err)

	createRegisterArg := CreateRegisterParams{TeacherID: teacher.ID, StudentID: studentR1.ID}
	err2 := testQueries.CreateRegister(context.Background(), createRegisterArg)
	require.NoError(t, err2)

	notificationArg := GetStudentsEmailForNotificationParams{
		TeacherID:     teacher.ID,
		StudentEmails: fmt.Sprintf("'%s'", studentR1.Email),
	}
	emails, err3 := testQueries.GetStudentsEmailForNotification(context.Background(), notificationArg)
	require.NoError(t, err3)
	require.NotEmpty(t, emails)

	require.Len(t, emails, 1)
	require.Contains(t, emails, studentR1.Email)
}

func TestGetStudentsEmailForNotification_SuspendShouldNotSend(t *testing.T) {
	teacher, err := createRandomTeacher()
	require.NoError(t, err)

	studentSusR1, err := createRandomStudent()
	require.NoError(t, err)

	studentR2, err := createRandomStudent()
	require.NoError(t, err)

	studentNr3, err := createRandomStudent()
	require.NoError(t, err)

	studentSus4, err := createRandomStudent()
	require.NoError(t, err)

	studentNr5, err := createRandomStudent()
	require.NoError(t, err)

	studentR6, err := createRandomStudent()
	require.NoError(t, err)

	createRegisterArg := CreateRegisterParams{TeacherID: teacher.ID, StudentID: studentSusR1.ID}
	err = testQueries.CreateRegister(context.Background(), createRegisterArg)
	require.NoError(t, err)

	createRegisterArg = CreateRegisterParams{TeacherID: teacher.ID, StudentID: studentR2.ID}
	err = testQueries.CreateRegister(context.Background(), createRegisterArg)
	require.NoError(t, err)

	createRegisterArg = CreateRegisterParams{TeacherID: teacher.ID, StudentID: studentR6.ID}
	err = testQueries.CreateRegister(context.Background(), createRegisterArg)
	require.NoError(t, err)

	err = testQueries.SuspendStudentByEmail(context.Background(), studentSusR1.Email)
	require.NoError(t, err)
	err = testQueries.SuspendStudentByEmail(context.Background(), studentSus4.Email)
	require.NoError(t, err)

	notificationArg := GetStudentsEmailForNotificationParams{
		TeacherID:     teacher.ID,
		StudentEmails: fmt.Sprintf("'%s','%s','%s'", studentNr3.Email, studentSus4.Email, studentNr5.Email),
	}

	emails, err := testQueries.GetStudentsEmailForNotification(context.Background(), notificationArg)
	require.NoError(t, err)
	require.NotEmpty(t, emails)

	require.Len(t, emails, 4)
	require.NotContains(t, emails, studentSusR1.Email)
	require.Contains(t, emails, studentR2.Email)
	require.Contains(t, emails, studentNr3.Email)
	require.NotContains(t, emails, studentSus4.Email)
	require.Contains(t, emails, studentNr5.Email)
	require.Contains(t, emails, studentR6.Email)
}
