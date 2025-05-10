package repository

import (
	"context"
)

// DataStore defines the interface for all data storage operations
type DataStore interface {
	Set(ctx context.Context, key, value string) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
}

// IDGeneration defines the interface for generating unique IDs
type IDGeneration interface {
	Generate() int64
	GetCurrent() int64
}

// Engine defines the interface for low-level storage operations
// This interface is implemented by concrete storage engines like in_memory
type Engine interface {
	Set(context.Context, string, string)
	Get(context.Context, string) (string, bool)
	Del(context.Context, string) bool
}
