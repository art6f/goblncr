package hashring

import (
	"fmt"

	"github.com/cespare/xxhash"
)

type Ring struct {
	nodesMap map[uint64]string
	ring     []string
	vnodes   uint64
}

func NewEmptyHashRing(vnodes uint64) *Ring {
	return &Ring{
		vnodes:   vnodes,
		nodesMap: map[uint64]string{},
		ring:     []string{},
	}
}

func NewHashRing(servers []string, vnodes uint64) *Ring {
	ring := NewEmptyHashRing(vnodes)

	for _, node := range servers {
		ring.AddNode(node)
	}

	return ring
}

func (ring *Ring) AddNode(node string) {
	for vnodeId := range ring.vnodes {
		vnodeKey := fmt.Sprintf("node-%s-%d", node, vnodeId)
		vnodeHash := ring.hashKey(vnodeKey)

		ring.nodesMap[vnodeHash] = node
	}
}

func (ring *Ring) RemoveNode(node string) {

}

func (ring *Ring) GetNode(key string) string {
	return ""
}

func (ring *Ring) hashKey(key string) uint64 {
	return xxhash.Sum64([]byte(key))
}
