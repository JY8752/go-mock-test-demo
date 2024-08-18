package gacha_test

import (
	"fmt"
	"go-mock-test-demo/gacha"
	"go-mock-test-demo/gacha/domain"
	"go-mock-test-demo/gacha/repository"
	"go-mock-test-demo/random"
	"go-mock-test-demo/tx"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDraw(t *testing.T) {
	// Arrange
	var (
		items = []*domain.Item{
			{ID: 1, Name: "item1", Rare: "N"},
			{ID: 2, Name: "item2", Rare: "N"},
			{ID: 3, Name: "item3", Rare: "N"},
			{ID: 4, Name: "item4", Rare: "N"},
			{ID: 5, Name: "item5", Rare: "N"},
			{ID: 6, Name: "item6", Rare: "R"},
			{ID: 7, Name: "item7", Rare: "R"},
			{ID: 8, Name: "item8", Rare: "R"},
			{ID: 9, Name: "item9", Rare: "R"},
			{ID: 10, Name: "item10", Rare: "SR"},
		}

		weights = []int{
			15,
			15,
			15,
			15,
			15,
			6,
			6,
			6,
			6,
			1,
		}

		userId int64 = 1
	)

	var (
		userRep = &repository.UserMock{
			FindByIdFunc: func(id int64) (*domain.User, error) {
				return &domain.User{Coin: 100}, nil
			},
		}
		itemRep = &repository.ItemMock{
			FindItemAndWeightsFunc: func() ([]*domain.Item, []int, error) {
				return items, weights, nil
			},
		}
		userItemRep = &repository.UserItemMock{}
		tx          = &tx.TransactionMock{}
		rnd         = &random.RandGeneratorMock{
			IntNFunc: func(n int) int {
				return 99
			},
		}
	)

	var (
		expectedItemName         = "item10"
		expectedItemRare         = "SR"
		expectedItemId     int64 = 10
		expectedGachaPrice       = 10
	)

	sut := gacha.NewGacha(userRep, itemRep, userItemRep, tx, rnd)

	// Act
	result, err := sut.Draw(userId)

	// Assertion
	require.Nil(t, err)
	assertResult(t, expectedItemName, expectedItemRare, result)

	if assert.Len(t, userItemRep.CreateWithTxCalls(), 1) {
		assert.Equal(t, userId, userItemRep.CreateWithTxCalls()[0].UserId)
		assert.Equal(t, expectedItemId, userItemRep.CreateWithTxCalls()[0].ItemId)
	}

	if assert.Len(t, userRep.DecreaseCoinsWithTxCalls(), 1) {
		assert.Equal(t, userId, userRep.DecreaseCoinsWithTxCalls()[0].UserId)
		assert.Equal(t, expectedGachaPrice, userRep.DecreaseCoinsWithTxCalls()[0].Amount)
	}

	assert.Len(t, userItemRep.IncrementCountWithTxCalls(), 0)

	assert.Len(t, tx.BeginCalls(), 1)
	assert.Len(t, tx.CommitCalls(), 1)
}

func assertResult(t *testing.T, itemName, rare, act string) {
	t.Helper()
	assert.Equal(t, fmt.Sprintf("{\"itemName\": %s, \"rare\": %s}\n", itemName, rare), act)
}
