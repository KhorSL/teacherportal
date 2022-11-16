package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetCommonStudentsEmail(t *testing.T) {
	teacher1, err := createRandomTeacher()
	require.NoError(t, err)

	teacher2, err := createRandomTeacher()
	require.NoError(t, err)

	commonStudent1, err := createRandomStudent()
	require.NoError(t, err)

	commonStudent2, err := createRandomStudent()
	require.NoError(t, err)

	t1Student, err := createRandomStudent()
	require.NoError(t, err)

	t2Student, err := createRandomStudent()
	require.NoError(t, err)

	// Teacher 1
	arg := CreateRegisterParams{TeacherID: teacher1.ID, StudentID: commonStudent1.ID}
	err = testQueries.CreateRegister(context.Background(), arg)
	require.NoError(t, err)

	arg = CreateRegisterParams{TeacherID: teacher1.ID, StudentID: commonStudent2.ID}
	err = testQueries.CreateRegister(context.Background(), arg)
	require.NoError(t, err)

	arg = CreateRegisterParams{TeacherID: teacher1.ID, StudentID: t1Student.ID}
	err = testQueries.CreateRegister(context.Background(), arg)
	require.NoError(t, err)

	// Teacher 2
	arg = CreateRegisterParams{TeacherID: teacher2.ID, StudentID: commonStudent1.ID}
	err = testQueries.CreateRegister(context.Background(), arg)
	require.NoError(t, err)

	arg = CreateRegisterParams{TeacherID: teacher2.ID, StudentID: commonStudent2.ID}
	err = testQueries.CreateRegister(context.Background(), arg)
	require.NoError(t, err)

	arg = CreateRegisterParams{TeacherID: teacher1.ID, StudentID: t2Student.ID}
	err = testQueries.CreateRegister(context.Background(), arg)
	require.NoError(t, err)

	commonStudentArg := GetCommonStudentsEmailParams{
		Email: fmt.Sprintf("'%s','%s'", teacher1.Email, teacher2.Email),
		Count: int64(2),
	}

	emails, err := testQueries.GetCommonStudentsEmail(context.Background(), commonStudentArg)
	require.NoError(t, err)
	require.NotEmpty(t, emails)

	require.Contains(t, emails, commonStudent1.Email)
	require.Contains(t, emails, commonStudent2.Email)
	require.NotContains(t, emails, t1Student.Email)
	require.NotContains(t, emails, t2Student.Email)
}

func TestGetCommonStudentsEmail_CommonStudentsBtwThreeTeachers(t *testing.T) {
	teacher1, err := createRandomTeacher()
	require.NoError(t, err)

	teacher2, err := createRandomTeacher()
	require.NoError(t, err)

	teacher3, err := createRandomTeacher()
	require.NoError(t, err)

	commonStudent1, err := createRandomStudent()
	require.NoError(t, err)

	t1t2Student2, err := createRandomStudent()
	require.NoError(t, err)

	t1Student, err := createRandomStudent()
	require.NoError(t, err)

	t2Student, err := createRandomStudent()
	require.NoError(t, err)

	// Teacher 1
	arg := CreateRegisterParams{TeacherID: teacher1.ID, StudentID: commonStudent1.ID}
	err = testQueries.CreateRegister(context.Background(), arg)
	require.NoError(t, err)

	arg = CreateRegisterParams{TeacherID: teacher1.ID, StudentID: t1t2Student2.ID}
	err = testQueries.CreateRegister(context.Background(), arg)
	require.NoError(t, err)

	arg = CreateRegisterParams{TeacherID: teacher1.ID, StudentID: t1Student.ID}
	err = testQueries.CreateRegister(context.Background(), arg)
	require.NoError(t, err)

	// Teacher 2
	arg = CreateRegisterParams{TeacherID: teacher2.ID, StudentID: commonStudent1.ID}
	err = testQueries.CreateRegister(context.Background(), arg)
	require.NoError(t, err)

	arg = CreateRegisterParams{TeacherID: teacher2.ID, StudentID: t1t2Student2.ID}
	err = testQueries.CreateRegister(context.Background(), arg)
	require.NoError(t, err)

	arg = CreateRegisterParams{TeacherID: teacher1.ID, StudentID: t2Student.ID}
	err = testQueries.CreateRegister(context.Background(), arg)
	require.NoError(t, err)

	// Teacher 3
	arg = CreateRegisterParams{TeacherID: teacher3.ID, StudentID: commonStudent1.ID}
	err = testQueries.CreateRegister(context.Background(), arg)
	require.NoError(t, err)

	commonStudentArg := GetCommonStudentsEmailParams{
		Email: fmt.Sprintf("'%s','%s', '%s'", teacher1.Email, teacher2.Email, teacher3.Email),
		Count: int64(3),
	}

	emails, err := testQueries.GetCommonStudentsEmail(context.Background(), commonStudentArg)
	require.NoError(t, err)
	require.NotEmpty(t, emails)

	require.Contains(t, emails, commonStudent1.Email)
	require.NotContains(t, emails, t1t2Student2.Email)
	require.NotContains(t, emails, t1Student.Email)
	require.NotContains(t, emails, t2Student.Email)
}
