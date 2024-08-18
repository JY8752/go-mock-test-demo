//go:generate moq -fmt goimports -rm -out user_item_moq.go -stub . UserItem

package repository

import (
	"database/sql"
	"errors"
	"go-mock-test-demo/tx"
)

type UserItem interface {
	Exist(userId, itemId int64) (bool, error)
	CreateWithTx(tx tx.Transaction, userId, itemId int64) (err error)
	IncrementCountWithTx(tx tx.Transaction, userId, itemId int64) (err error)
}

type userItem struct {
	db *sql.DB
}

func NewUserItem(db *sql.DB) *userItem {
	return &userItem{db: db}
}

func (u *userItem) Exist(userId, itemId int64) (bool, error) {
	if err := u.db.QueryRow(
		"SELECT id, user_id, item_id, count FROM user_items WHERE user_id = ? AND item_id = ?", userId, itemId,
	).Err(); err != nil && errors.Is(err, sql.ErrNoRows) {
		return false, nil
	} else if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return false, err
	} else {
		return true, nil
	}
}

func (u *userItem) CreateWithTx(tx tx.Transaction, userId, itemId int64) (err error) {
	return tx.Exec("INSERT INTO user_items(id, user_id, item_id, count) VALUES(NULL, ?, ?, 1)", userId, itemId)
}

func (u *userItem) IncrementCountWithTx(tx tx.Transaction, userId, itemId int64) (err error) {
	return tx.Exec("UPDATE user_items SET count = count + 1 WHERE user_id = ? AND item_id = ?", userId, itemId)
}
