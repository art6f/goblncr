package strategies

type HashringStrategy struct {
	BaseStrategy
}

var _ Strategy = (*HashringStrategy)(nil)

func NewHashringStrategy() *HashringStrategy {
	return &HashringStrategy{}
}

func (h *HashringStrategy) Select(servers []string) (string, error) {
	panic("unimplemented")
}
