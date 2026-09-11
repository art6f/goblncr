package hashring

import (
	"fmt"
	"sort"

	"github.com/cespare/xxhash"
)

type Ring struct {
	nodesMap map[uint64]string
	ring     []uint64
	vnodes   uint64
}

func NewEmptyHashRing(vnodes uint64) *Ring {
	return &Ring{
		vnodes:   vnodes,
		nodesMap: map[uint64]string{},
		ring:     []uint64{},
	}
}

func NewHashRing(servers []string, vnodes uint64) *Ring {
	ring := NewEmptyHashRing(vnodes)

	for _, node := range servers {
		ring.AddNode(node)
	}

	return ring
}

func (ring *Ring) AddNode(newNode string) {
	for vnodeId := range ring.vnodes {
		nodeKey := fmt.Sprintf("node-%s-%d", newNode, vnodeId)
		nodeHash := ring.hashKey(nodeKey)

		ring.nodesMap[nodeHash] = newNode
		ring.ring = append(ring.ring, nodeHash)
	}

	sort.Slice(ring.ring, func(i, j int) bool {
		return ring.ring[i] < ring.ring[j]
	})
}

func (ring *Ring) RemoveNode(targetNode string) {
	newMap := make(map[uint64]string, 0)
	newRing := make([]uint64, 0)

	for _, nodeHash := range ring.ring {
		currNode := ring.nodesMap[nodeHash]

		if currNode != targetNode {
			newRing = append(newRing, nodeHash)
			newMap[nodeHash] = currNode
		}
	}

	ring.nodesMap, ring.ring = newMap, newRing
}

func (ring *Ring) GetNode(key string) string {
	return ""
}

func (ring *Ring) hashKey(key string) uint64 {
	return xxhash.Sum64([]byte(key))
}
