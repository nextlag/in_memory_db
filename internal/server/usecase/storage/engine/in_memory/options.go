package in_memory

type EngineOptions func(e *Engine)

func WithPartitions(partition uint) EngineOptions {
	return func(e *Engine) {
		e.partitions = make([]*HashTable, partition)
		for i := range e.partitions {
			e.partitions[i] = NewHashTable()
		}
	}
}
