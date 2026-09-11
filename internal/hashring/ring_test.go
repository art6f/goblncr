package hashring

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDistribution(t *testing.T) {
	testNodes := []string{"node1", "node2", "node3"}
	ring := NewHashRing(testNodes, 10)
	assert.Len(t, ring.ring, 30)

	ring.RemoveNode("node1")
	assert.Len(t, ring.ring, 20)
}
