package cache

import "sync"

// Could be extended with TTL for example.
type Entry struct {
	Value []byte
}

type Cache struct {
	mu      sync.RWMutex
	entries map[string]*Entry
}

func NewCache() *Cache {
	return &Cache{
		entries: make(map[string]*Entry),
	}
}

// - Input: key.
// - Returns the stored value if present.
func (store *Cache) Get(key string) *Entry {
	store.mu.RLock()
	defer store.mu.RUnlock()
	return store.entries[key]
}

// SET
func (store *Cache) Set(key string, value []byte) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.entries[key] = &Entry{
		Value: value,
	}
}
