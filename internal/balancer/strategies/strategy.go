package strategies

import (
	"errors"
	"strings"

	"github.com/art6f/goblncr/internal/k8s"
	"go.yaml.in/yaml/v3"
)

type BalancingStrategy int

const (
	StrategyRandom BalancingStrategy = iota
	StrategyHashring
)

// YAML string to type resolver
func (s *BalancingStrategy) UnmarshalYAML(node *yaml.Node) error {
	switch strings.ToLower(node.Value) {
	case "hashring":
		*s = StrategyHashring
	default:
		*s = StrategyRandom
	}

	return nil
}

// stringyfy the name
func (s BalancingStrategy) String() string {
	switch s {
	case StrategyRandom:
		return "random"
	case StrategyHashring:
		return "hashring"
	default:
		return "unknown"
	}
}

func (s *BalancingStrategy) ResolveStrategy() (*Strategy, error) {
	var strategy Strategy

	switch *s {
	case StrategyRandom:
		strategy = NewRandomStrategy()
	case StrategyHashring:
		strategy = NewHashringStrategy()
	default:
		strategy = nil
	}

	if strategy == nil {
		return nil, errors.New("Not implemented")
	}

	return &strategy, nil
}

type BaseStrategy struct {
	pods k8s.PodsMap
}

type Strategy interface {
	SetPods(pods *k8s.PodsMap)
	Select(servers []string) string
}

func (base *BaseStrategy) SetPods(pods *k8s.PodsMap) {
	base.pods = *pods
}
