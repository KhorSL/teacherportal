package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegisterTx(t *testing.T) {
	store := NewStore(testDB)

	teacher, err := createRandomTeacher()
	require.NoError(t, err)

	student1, err := createRandomStudent()
	require.NoError(t, err)
	student2, err := createRandomStudent()
	require.NoError(t, err)

	students := []string{student1.Email, student2.Email}

	err = store.RegisterTx(context.Background(), RegisterTxParams{
		Teacher:  teacher.Email,
		Students: students,
	})
	require.NoError(t, err)

	// check registration
	register1, err := testQueries.GetRegisterByStudentIdAndTeacherId(context.Background(), GetRegisterByStudentIdAndTeacherIdParams{
		StudentID: student1.ID,
		TeacherID: teacher.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, register1)

	require.Equal(t, register1.StudentID, student1.ID)
	require.Equal(t, register1.TeacherID, teacher.ID)

	register2, err := testQueries.GetRegisterByStudentIdAndTeacherId(context.Background(), GetRegisterByStudentIdAndTeacherIdParams{
		StudentID: student1.ID,
		TeacherID: teacher.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, register2)

	require.Equal(t, register2.StudentID, student1.ID)
	require.Equal(t, register2.TeacherID, teacher.ID)
}

func TestRegisterTx_TeacherNotExist(t *testing.T) {
	store := NewStore(testDB)

	notExistEmail := "not_exist@email.com"

	student1, err := createRandomStudent()
	require.NoError(t, err)
	student2, err := createRandomStudent()
	require.NoError(t, err)

	students := []string{student1.Email, student2.Email}

	err = store.RegisterTx(context.Background(), RegisterTxParams{
		Teacher:  notExistEmail,
		Students: students,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "teacher does not exist")
}

func TestRegisterTx_StudentNotExist(t *testing.T) {
	store := NewStore(testDB)

	teacher, err := createRandomTeacher()
	require.NoError(t, err)

	notExistEmail := "not_exist@email.com"

	student1, err := createRandomStudent()
	require.NoError(t, err)

	students := []string{student1.Email, notExistEmail}

	err = store.RegisterTx(context.Background(), RegisterTxParams{
		Teacher:  teacher.Email,
		Students: students,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, fmt.Sprintf("student %s does not exist", notExistEmail))
}

func TestRegisterTx_DuplicateRegisters(t *testing.T) {
	store := NewStore(testDB)

	teacher, err := createRandomTeacher()
	require.NoError(t, err)

	student1, err := createRandomStudent()
	require.NoError(t, err)

	students := []string{student1.Email}

	err = store.RegisterTx(context.Background(), RegisterTxParams{
		Teacher:  teacher.Email,
		Students: students,
	})
	require.NoError(t, err)

	err = store.RegisterTx(context.Background(), RegisterTxParams{
		Teacher:  teacher.Email,
		Students: students,
	})
	require.Error(t, err)
	require.ErrorContains(t, err, "previously registered")
}
