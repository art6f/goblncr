package strategies

import (
	"errors"
	"math/rand"
)

type RandomStrategy struct {
	BaseStrategy
}

var _ Strategy = (*RandomStrategy)(nil)

func NewRandomStrategy() *RandomStrategy {
	return &RandomStrategy{}
}

func (r *RandomStrategy) Select(pods []string) (string, error) {
	if len(pods) == 0 {
		return "", errors.New("no pods to select from")
	}

	if len(pods) == 1 {
		return pods[0], nil
	}

	return pods[rand.Intn(len(pods))], nil
}
