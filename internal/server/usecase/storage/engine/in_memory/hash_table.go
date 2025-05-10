package in_memory

import "sync"

type HashTable struct {
	mtx  sync.Mutex
	data map[string]string
}

func NewHashTable() *HashTable {
	return &HashTable{data: make(map[string]string)}
}

func (h *HashTable) Set(key, value string) {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	h.data[key] = value
}

func (h *HashTable) Get(key string) (string, bool) {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	value, ok := h.data[key]

	return value, ok
}

func (h *HashTable) Del(key string) bool {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	if _, ok := h.data[key]; !ok {
		return false
	}

	delete(h.data, key)

	return true
}
