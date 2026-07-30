package strategies

// Strategy interface defines the contract for load balancing strategies
type Strategy interface {
	Select(servers []string) string
}