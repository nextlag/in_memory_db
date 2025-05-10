package repository

import (
	"context"
	"log/slog"

	"github.com/nextlag/in_memory_db/internal/common"
	"github.com/nextlag/in_memory_db/internal/errs"
)

// Repository implements DataStore interface using the in-memory storage engine
type Repository struct {
	engine Engine
	log    *slog.Logger
}

// New creates a new repository instance with the given engine
func New(engine Engine, logger *slog.Logger) (*Repository, error) {
	if engine == nil {
		return nil, errs.ErrEngineIsNil
	}

	if logger == nil {
		return nil, errs.ErrLoggerIsNil
	}

	return &Repository{
		engine: engine,
		log:    logger,
	}, nil
}

// Set stores a key-value pair
func (r *Repository) Set(ctx context.Context, key, value string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	ctx = common.ContextWithTxID(ctx, 0)
	r.engine.Set(ctx, key, value)

	return nil
}

// Get retrieves a value by key
func (r *Repository) Get(ctx context.Context, key string) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	ctx = common.ContextWithTxID(ctx, 0)
	value, found := r.engine.Get(ctx, key)
	if !found {
		return "", errs.ErrNotFound
	}

	return value, nil
}

// Del removes a key-value pair
func (r *Repository) Del(ctx context.Context, key string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	ctx = common.ContextWithTxID(ctx, 0)
	if ok := r.engine.Del(ctx, key); !ok {
		return errs.ErrNotFound
	}

	return nil
}
