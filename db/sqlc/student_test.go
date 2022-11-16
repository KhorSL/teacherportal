package db

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/khorsl/teacherportal/util"
	"github.com/stretchr/testify/require"
)

func createRandomStudent() (Student, error) {
	arg := CreateStudentParams{
		FullName: util.RandomName(),
		Email:    util.RandomEmail(),
	}

	err := testQueries.CreateStudent(context.Background(), arg)
	if err != nil {
		log.Fatal("cannot create random student:", err)
		return Student{}, err
	}

	student, err := testQueries.GetStudentByEmail(context.Background(), arg.Email)

	return student, err
}

func TestCreateStudent(t *testing.T) {
	arg := CreateStudentParams{
		FullName: util.RandomName(),
		Email:    util.RandomEmail(),
	}

	err := testQueries.CreateStudent(context.Background(), arg)

	require.NoError(t, err)
}

func TestGetStudentById(t *testing.T) {
	student1, _ := createRandomStudent()
	student2, err := testQueries.GetStudentById(context.Background(), student1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, student2)

	require.Equal(t, student1.ID, student2.ID)
	require.Equal(t, student1.FullName, student2.FullName)
	require.Equal(t, student1.Email, student2.Email)

	require.Equal(t, student1.IsActive, student2.IsActive)
	require.Equal(t, true, student2.IsActive)

	require.Equal(t, student1.IsSuspended, student2.IsSuspended)
	require.Equal(t, false, student2.IsSuspended)

	require.WithinDuration(t, student1.SuspendedAt, student2.SuspendedAt, time.Second)
	require.Zero(t, student2.SuspendedAt)

	require.WithinDuration(t, student1.CreatedAt, student2.CreatedAt, time.Second)
	require.WithinDuration(t, student1.UpdatedAt, student2.UpdatedAt, time.Second)

	require.WithinDuration(t, student1.DeletedAt, student2.DeletedAt, time.Second)
	require.Zero(t, student2.DeletedAt)
}

func TestGetStudentByEmail(t *testing.T) {
	student1, _ := createRandomStudent()
	student2, err := testQueries.GetStudentByEmail(context.Background(), student1.Email)

	require.NoError(t, err)
	require.NotEmpty(t, student2)

	require.Equal(t, student1.ID, student2.ID)
	require.Equal(t, student1.FullName, student2.FullName)
	require.Equal(t, student1.Email, student2.Email)
	require.Equal(t, student1.IsActive, student2.IsActive)
	require.WithinDuration(t, student1.CreatedAt, student2.CreatedAt, time.Second)
	require.WithinDuration(t, student1.UpdatedAt, student2.UpdatedAt, time.Second)
	require.WithinDuration(t, student1.DeletedAt, student2.DeletedAt, time.Second)
}

func TestSuspendStudentByEmail(t *testing.T) {
	student1, err := createRandomStudent()
	require.NoError(t, err)

	err = testQueries.SuspendStudentByEmail(context.Background(), student1.Email)
	require.NoError(t, err)

	student2, err := testQueries.GetStudentById(context.Background(), student1.ID)
	require.NoError(t, err)

	require.Equal(t, student1.ID, student2.ID)
	require.Equal(t, true, student2.IsSuspended)
	require.WithinDuration(t, time.Now(), student2.UpdatedAt, time.Second)
}

func TestGetNotSuspendedStudentByEmail_NotSuspend(t *testing.T) {
	student1, _ := createRandomStudent()
	student2, err := testQueries.GetNotSuspendedStudentByEmail(context.Background(), student1.Email)

	require.NoError(t, err)
	require.NotEmpty(t, student2)

	require.Equal(t, student1.ID, student2.ID)
	require.Equal(t, student1.FullName, student2.FullName)
	require.Equal(t, student1.Email, student2.Email)
	require.Equal(t, student1.IsActive, student2.IsActive)
	require.WithinDuration(t, student1.CreatedAt, student2.CreatedAt, time.Second)
	require.WithinDuration(t, student1.UpdatedAt, student2.UpdatedAt, time.Second)
	require.WithinDuration(t, student1.DeletedAt, student2.DeletedAt, time.Second)
}

func TestGetNotSuspendedStudentByEmail_Suspended(t *testing.T) {
	student1, err := createRandomStudent()
	require.NoError(t, err)

	err = testQueries.SuspendStudentByEmail(context.Background(), student1.Email)
	require.NoError(t, err)

	student2, err := testQueries.GetNotSuspendedStudentByEmail(context.Background(), student1.Email)

	require.Error(t, err)
	require.Empty(t, student2)
	require.ErrorContains(t, err, "no rows")
}
