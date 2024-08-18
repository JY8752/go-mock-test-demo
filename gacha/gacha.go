package gacha

import (
	"errors"
	"fmt"
	"go-mock-test-demo/gacha/repository"
	"go-mock-test-demo/random"
	"go-mock-test-demo/tx"
)

type Gacha struct {
	userRep     repository.User
	itemRep     repository.Item
	userItemRep repository.UserItem
	tx          tx.Transaction
	rnd         random.RandGenerator
}

func NewGacha(
	u repository.User,
	i repository.Item,
	ui repository.UserItem,
	tx tx.Transaction,
	rnd random.RandGenerator,
) *Gacha {
	return &Gacha{userRep: u, itemRep: i, userItemRep: ui, tx: tx, rnd: rnd}
}

const (
	GachaPrice = 10
)

func (g *Gacha) Draw(userId int64) (string, error) {
	// ユーザー情報を取得する
	user, err := g.userRep.FindById(userId)
	if err != nil {
		return "", err
	}

	if user.Coin < GachaPrice {
		return "", errors.New("コインが足りません")
	}

	// ガチャを抽選する
	items, weights, err := g.itemRep.FindItemAndWeights()
	if err != nil {
		return "", err
	}

	i, err := linearSearchLottery(weights, g.rnd)
	if err != nil {
		return "", fmt.Errorf("アイテムの抽選に失敗しました err: %v", err)
	}

	result := items[i]
	itemId := result.ID

	// 取得したアイテムを記録する
	exist, err := g.userItemRep.Exist(userId, itemId)
	if err != nil {
		return "", err
	}

	firstGet := !exist

	if err = g.tx.Begin(); err != nil {
		return "", err
	}

	if firstGet {
		if err = g.userItemRep.CreateWithTx(g.tx, userId, itemId); err != nil {
			return "", err
		}
	} else {
		if err = g.userItemRep.IncrementCountWithTx(g.tx, userId, itemId); err != nil {
			return "", err
		}
	}

	// コインを消費する
	if err = g.userRep.DecreaseCoinsWithTx(g.tx, userId, GachaPrice); err != nil {
		_ = g.tx.Rollback()
		return "", err
	}

	g.tx.Commit()

	return fmt.Sprintf("{\"itemName\": %s, \"rare\": %s}\n", result.Name, result.Rare), nil
}

/*
線形探索で重み付抽選する
@return 当選した要素のインデックス
*/
func linearSearchLottery(weights []int, rg random.RandGenerator) (int, error) {
	//  重みの総和を取得する
	var total int
	for _, weight := range weights {
		total += weight
	}

	// 疑似乱数の取得
	rnd := rg.IntN(total)

	var currentWeight int
	for i, w := range weights {
		// 現在要素までの重みの総和
		currentWeight += w

		if rnd < currentWeight {
			return i, nil
		}
	}

	// たぶんありえない
	return 0, errors.New("the lottery failed")
}
