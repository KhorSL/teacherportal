package db

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/khorsl/teacherportal/util"
	"github.com/stretchr/testify/require"
)

func createRandomTeacher() (Teacher, error) {
	arg := CreateTeacherParams{
		FullName: util.RandomName(),
		Email:    util.RandomEmail(),
	}

	err := testQueries.CreateTeacher(context.Background(), arg)
	if err != nil {
		log.Fatal("cannot create random teacher:", err)
		return Teacher{}, err
	}

	teacher, err := testQueries.GetTeacherByEmail(context.Background(), arg.Email)

	return teacher, err
}

func TestCreateTeacher(t *testing.T) {
	arg := CreateTeacherParams{
		FullName: util.RandomName(),
		Email:    util.RandomEmail(),
	}

	err := testQueries.CreateTeacher(context.Background(), arg)

	require.NoError(t, err)
}

func TestGetTeacherById(t *testing.T) {
	teacher1, _ := createRandomTeacher()
	teacher2, err := testQueries.GetTeacherById(context.Background(), teacher1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, teacher2)

	require.Equal(t, teacher1.ID, teacher2.ID)
	require.Equal(t, teacher1.FullName, teacher2.FullName)
	require.Equal(t, teacher1.Email, teacher2.Email)
	require.Equal(t, teacher1.IsActive, teacher2.IsActive)
	require.Equal(t, true, teacher2.IsActive)
	require.WithinDuration(t, teacher1.CreatedAt, teacher2.CreatedAt, time.Second)
	require.WithinDuration(t, teacher1.UpdatedAt, teacher2.UpdatedAt, time.Second)
	require.WithinDuration(t, teacher1.DeletedAt, teacher2.DeletedAt, time.Second)
	require.Zero(t, teacher2.DeletedAt)
}

func TestGetTeacherByEmail(t *testing.T) {
	teacher1, _ := createRandomTeacher()
	teacher2, err := testQueries.GetTeacherByEmail(context.Background(), teacher1.Email)

	require.NoError(t, err)
	require.NotEmpty(t, teacher2)

	require.Equal(t, teacher1.ID, teacher2.ID)
	require.Equal(t, teacher1.FullName, teacher2.FullName)
	require.Equal(t, teacher1.Email, teacher2.Email)
	require.Equal(t, teacher1.IsActive, teacher2.IsActive)
	require.WithinDuration(t, teacher1.CreatedAt, teacher2.CreatedAt, time.Second)
	require.WithinDuration(t, teacher1.UpdatedAt, teacher2.UpdatedAt, time.Second)
	require.WithinDuration(t, teacher1.DeletedAt, teacher2.DeletedAt, time.Second)
}
