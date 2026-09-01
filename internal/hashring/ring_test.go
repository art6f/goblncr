package hashring

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDistribution(t *testing.T) {
	newVar := []string{"node1", "node2", "node3"}
	ring := NewHashRing(newVar, 10)

	assert.Equal(t, 30, ring.ring)
}
