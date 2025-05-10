package util

import (
	"math"
	"sync/atomic"
)

// IDGenerator provides unique sequential IDs with overflow protection
type IDGenerator struct {
	counter atomic.Int64
}

// NewIDGenerator creates a new ID generator with a starting value
func NewIDGenerator(previousID int64) *IDGenerator {
	generator := &IDGenerator{}
	generator.counter.Store(previousID)

	return generator
}

// Generate returns a new unique ID and increments the counter
// If the counter reaches MaxInt64, it will safely reset to 0
func (g *IDGenerator) Generate() int64 {
	for {
		current := g.counter.Load()
		if current == math.MaxInt64 {
			if g.counter.CompareAndSwap(math.MaxInt64, 0) {
				return 1
			}
			continue
		}

		next := current + 1
		if g.counter.CompareAndSwap(current, next) {
			return next
		}
	}
}

// GetCurrent returns the current counter value without incrementing
func (g *IDGenerator) GetCurrent() int64 {
	return g.counter.Load()
}
