package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

type Store interface {
	Querier
	RegisterTx(ctx context.Context, arg RegisterTxParams) error
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

func (store *SQLStore) execTx(ctx context.Context, callbackFn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	queries := New(tx)
	err = callbackFn(queries)
	if err != nil {
		rollbackErr := tx.Rollback()
		if rollbackErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rollbackErr)
		}
		return err
	}

	return tx.Commit()
}

type RegisterTxParams struct {
	Teacher  string   `json:"teacher"`
	Students []string `json:"students"`
}

func (store *SQLStore) RegisterTx(ctx context.Context, arg RegisterTxParams) error {
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		teacher, err := q.GetTeacherByEmail(ctx, arg.Teacher)
		if err != nil {
			if strings.Contains(err.Error(), "no rows") {
				return errors.New("teacher does not exist")
			}
			return fmt.Errorf("could not retrieve teacher: %s", err.Error())
		}

		for i := 0; i < len(arg.Students); i++ {
			student, err := q.GetStudentByEmail(ctx, arg.Students[i])
			if err != nil {
				if strings.Contains(err.Error(), "no rows") {
					return fmt.Errorf("student %s does not exist", arg.Students[i])
				}
				return fmt.Errorf("could not retrieve student %s: %s", arg.Students[i], err.Error())
			}

			registerArg := CreateRegisterParams{
				TeacherID: teacher.ID,
				StudentID: student.ID,
			}

			err = q.CreateRegister(ctx, registerArg)
			if err != nil {
				if strings.Contains(err.Error(), "Duplicate") {
					return fmt.Errorf("student id %d previously registered", student.ID)
				}
				return fmt.Errorf("could not register student id %d: %s", student.ID, err.Error())
			}
		}
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
