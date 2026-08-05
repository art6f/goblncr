package strategies

import "math/rand"

type RandomStrategy struct {
	BaseStrategy
}

var _ Strategy = (*RandomStrategy)(nil)

func NewRandomStrategy() *RandomStrategy {
	return &RandomStrategy{}
}

func (r *RandomStrategy) Select(pods []string) string {
	if len(pods) == 0 {
		return ""
	}

	if len(pods) == 1 {
		return pods[0]
	}

	return pods[rand.Intn(len(pods))]
}
