//go:generate moq -fmt goimports -rm -out item_moq.go -stub . Item

package repository

import (
	"database/sql"
	"fmt"
	"go-mock-test-demo/gacha/domain"
)

type Item interface {
	FindItemAndWeights() (items []*domain.Item, weights []int, err error)
}

type item struct {
	db *sql.DB
}

func NewItem(db *sql.DB) *item {
	return &item{db: db}
}

func (i *item) FindItemAndWeights() (items []*domain.Item, weights []int, err error) {
	rows, err := i.db.Query("SELECT id, name, rare, weight FROM items")
	if err != nil {
		return nil, nil, fmt.Errorf("アイテムの取得に失敗しました err: %v", err)
	}
	defer func() {
		_ = rows.Close()
	}()

	columns, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	items = make([]*domain.Item, 0, len(columns))
	weights = make([]int, 0, len(columns))

	for rows.Next() {
		var item domain.Item
		if err = rows.Scan(&item.ID, &item.Name, &item.Rare, &item.Weight); err != nil {
			return nil, nil, err
		}
		items = append(items, &item)
		weights = append(weights, item.Weight)
	}

	return items, weights, nil
}
