//go:generate moq -fmt goimports -rm -out user_moq.go -stub . User

package repository

import (
	"database/sql"
	"fmt"
	"go-mock-test-demo/gacha/domain"
	"go-mock-test-demo/tx"
)

type User interface {
	FindById(id int64) (user *domain.User, err error)
	DecreaseCoinsWithTx(tx tx.Transaction, userId int64, amount int) (err error)
}

type user struct {
	db *sql.DB
}

func NewUser(db *sql.DB) *user {
	return &user{db: db}
}

func (u *user) FindById(id int64) (*domain.User, error) {
	var user domain.User
	if err := u.db.QueryRow("SELECT id, name, coin FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Coin); err != nil {
		return nil, fmt.Errorf("ユーザーの取得に失敗しました err: %v", err)
	}
	return &user, nil
}

func (u *user) DecreaseCoinsWithTx(tx tx.Transaction, userId int64, amount int) (err error) {
	return tx.Exec("UPDATE users SET coin = coin - ? WHERE id = ?", amount, userId)
}
