//go:generate moq -fmt goimports -rm -out tx_moq.go -stub . Transaction

package tx

import (
	"database/sql"
	"errors"
)

const (
	NotYetCompletedErr = "not yet commit or rollback"
	NotBeginErr        = "not begin transaction"
)

type Transaction interface {
	Begin() error
	Commit() error
	Rollback() error
	Exec(query string, args ...any) error
}

type transaction struct {
	db *sql.DB
	tx *sql.Tx
}

func NewTransaction(db *sql.DB) *transaction {
	return &transaction{db: db}
}

func (t *transaction) Begin() error {
	if t.tx != nil {
		return errors.New(NotYetCompletedErr)
	}

	tx, err := t.db.Begin()
	if err != nil {
		return err
	}
	t.tx = tx

	return nil
}

func (t *transaction) Commit() error {
	if t.tx == nil {
		return errors.New(NotBeginErr)
	}

	if err := t.tx.Commit(); err != nil {
		return err
	}
	t.tx = nil

	return nil
}

func (t *transaction) Rollback() error {
	if t.tx == nil {
		return errors.New(NotBeginErr)
	}

	if err := t.tx.Rollback(); err != nil {
		return err
	}
	t.tx = nil

	return nil
}

func (t *transaction) Exec(query string, args ...any) error {
	if t.tx == nil {
		return errors.New(NotBeginErr)
	}

	_, err := t.tx.Exec(query, args...)
	return err
}
