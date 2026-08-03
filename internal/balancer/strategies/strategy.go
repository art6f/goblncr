package strategies

type BalancingStrategies int

const (
	StrategyRandom BalancingStrategies = iota
	StrategyHashring
)

// Strategy interface defines the contract for load balancing strategies
type Strategy interface {
	Select(servers []string) string
}
