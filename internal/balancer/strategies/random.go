package strategies

import "math/rand"

// RandomStrategy implements the Strategy interface for random distribution
type RandomStrategy struct {
}

// NewRandomStrategy creates a new RandomStrategy instance
func NewRandomStrategy() *RandomStrategy {
	return &RandomStrategy{}
}

// Select selects a server randomly from the available servers
func (r *RandomStrategy) Select(pods []string) string {
	if len(pods) == 0 {
		return ""
	}

	if len(pods) == 1 {
		return pods[0]
	}

	return pods[rand.Intn(len(pods))]
}
