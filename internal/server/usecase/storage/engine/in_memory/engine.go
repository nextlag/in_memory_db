package in_memory

import (
	"context"
	"hash/fnv"
	"log/slog"

	"github.com/nextlag/in_memory_db/internal/common"
	"github.com/nextlag/in_memory_db/internal/errs"
)

const (
	defaultPartitionCount = 1
	defaultPartitionIndex = 0
)

type Engine struct {
	partitions []*HashTable
	log        *slog.Logger
}

func NewEngine(log *slog.Logger, options ...EngineOptions) (*Engine, error) {
	if log == nil {
		return nil, errs.ErrLoggerIsNil
	}

	e := &Engine{
		log: log,
	}

	for _, opt := range options {
		opt(e)
	}

	if len(e.partitions) == 0 {
		e.partitions = make([]*HashTable, defaultPartitionCount)
		e.partitions[0] = NewHashTable()
	}

	return e, nil
}

func (e *Engine) Set(ctx context.Context, key, value string) {
	partitionIdx := defaultPartitionIndex

	if len(e.partitions) > 1 {
		partitionIdx = e.partitionIdx(key)
	}

	partition := e.partitions[partitionIdx]
	partition.Set(key, value)
	txID := common.GetTxIDFromContext(ctx)

	e.log.Debug("successfully set query", "txID", txID)
}

func (e *Engine) Get(ctx context.Context, key string) (value string, ok bool) {
	partitionIdx := defaultPartitionIndex

	if len(e.partitions) > 1 {
		partitionIdx = e.partitionIdx(key)
	}

	partition := e.partitions[partitionIdx]
	value, ok = partition.Get(key)
	txID := common.GetTxIDFromContext(ctx)

	e.log.Debug("successfully get query", "txID", txID)

	return
}

func (e *Engine) Del(ctx context.Context, key string) bool {
	partitionIdx := defaultPartitionIndex
	if len(e.partitions) > 1 {
		partitionIdx = e.partitionIdx(key)
	}

	partition := e.partitions[partitionIdx]
	deleted := partition.Del(key)

	txID := common.GetTxIDFromContext(ctx)
	e.log.Debug("delete query processed", "txID", txID, "key", key, "deleted", deleted)

	return deleted
}

func (e *Engine) partitionIdx(key string) int {
	hash := fnv.New32a()
	_, _ = hash.Write([]byte(key))
	return int(hash.Sum32()) % len(e.partitions)
}
