package db

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomRegister(t *testing.T) (Register, error) {
	teacher, err := createRandomTeacher()
	require.NoError(t, err)

	student, err := createRandomStudent()
	require.NoError(t, err)

	arg1 := CreateRegisterParams{
		TeacherID: teacher.ID,
		StudentID: student.ID,
	}

	err = testQueries.CreateRegister(context.Background(), arg1)
	if err != nil {
		log.Fatal("cannot create random student:", err)
		return Register{}, err
	}

	arg2 := GetRegisterByStudentIdAndTeacherIdParams{
		StudentID: student.ID,
		TeacherID: teacher.ID,
	}

	register, err := testQueries.GetRegisterByStudentIdAndTeacherId(context.Background(), arg2)

	return register, err
}

func TestCreateRegister(t *testing.T) {
	teacher, err := createRandomTeacher()
	require.NoError(t, err)

	student, err := createRandomStudent()
	require.NoError(t, err)

	arg := CreateRegisterParams{
		TeacherID: teacher.ID,
		StudentID: student.ID,
	}

	err = testQueries.CreateRegister(context.Background(), arg)

	require.NoError(t, err)
}

func TestGetRegisterByStudentIdAndTeacherId(t *testing.T) {
	register1, err := createRandomRegister(t)

	require.NoError(t, err)
	require.NotEmpty(t, register1)

	arg := GetRegisterByStudentIdAndTeacherIdParams{
		StudentID: register1.StudentID,
		TeacherID: register1.TeacherID,
	}

	register2, err := testQueries.GetRegisterByStudentIdAndTeacherId(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, register2)

	require.Equal(t, register1.ID, register2.ID)
	require.Equal(t, register1.StudentID, register2.StudentID)
	require.Equal(t, register1.TeacherID, register2.TeacherID)
	require.WithinDuration(t, register1.CreatedAt, register2.CreatedAt, time.Second)
}

func TestGetRegisterByTeacherId(t *testing.T) {
	var lastRegister Register
	for i := 0; i < 10; i++ {
		lastRegister, _ = createRandomRegister(t)
	}

	registers, err := testQueries.GetRegisterByTeacherId(context.Background(), lastRegister.TeacherID)
	require.NoError(t, err)
	require.NotEmpty(t, registers)

	for _, transfer := range registers {
		require.NotEmpty(t, transfer)
		require.Equal(t, lastRegister.StudentID, transfer.StudentID)
		require.Equal(t, lastRegister.TeacherID, transfer.TeacherID)
	}
}

func TestGetRegisterByStudentId(t *testing.T) {
	var lastRegister Register
	for i := 0; i < 10; i++ {
		lastRegister, _ = createRandomRegister(t)
	}

	registers, err := testQueries.GetRegisterByStudentId(context.Background(), lastRegister.StudentID)
	require.NoError(t, err)
	require.NotEmpty(t, registers)

	for _, transfer := range registers {
		require.NotEmpty(t, transfer)
		require.Equal(t, lastRegister.StudentID, transfer.StudentID)
		require.Equal(t, lastRegister.TeacherID, transfer.TeacherID)
	}
}
