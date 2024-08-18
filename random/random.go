//go:generate moq -fmt goimports -rm -out random_moq.go -stub . RandGenerator

package random

import "math/rand/v2"

type RandGenerator interface {
	IntN(n int) int
}

type randGenerator struct {
}

func NewRandGenerator() *randGenerator {
	return &randGenerator{}
}

func (r *randGenerator) IntN(n int) int {
	return rand.IntN(n)
}
