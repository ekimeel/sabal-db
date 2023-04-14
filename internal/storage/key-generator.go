package storage

import "sync/atomic"

type KeyGenerator[K comparable] interface {
	NextKey() K
}

type SerialKeyGenerator struct {
	lastKey uint32
}

func NewSerialKeyGenerator(startKey uint32) *SerialKeyGenerator {
	return &SerialKeyGenerator{lastKey: startKey}
}

func (g *SerialKeyGenerator) NextKey() uint32 {
	return atomic.AddUint32(&g.lastKey, 1)
}
