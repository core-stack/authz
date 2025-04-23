package store

import (
	"context"
	"reflect"
	"sync"
	"time"
)

type MemoryStore[T StoreData] struct {
	mu   sync.RWMutex
	data map[string]T
}

func NewMemoryStore[T StoreData]() Store[T] {
	store := &MemoryStore[T]{
		data: make(map[string]T),
	}
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			store.cleanupExpired()
		}
	}()
	return store
}

func (m *MemoryStore[T]) Set(ctx context.Context, s T) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[s.GetID()] = s
	return nil
}

func (m *MemoryStore[T]) Get(ctx context.Context, id string) (*T, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.data[id]
	if !ok || time.Now().After(s.GetExpiresAt()) {
		if ok {
			m.Delete(ctx, id)
		}
		return nil, ErrDataNotFound
	}
	return &s, nil
}

func (s *MemoryStore[T]) GetByFilter(ctx context.Context, filter map[string]any, limit, offset int) ([]*T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []*T
	count := 0

	for _, item := range s.data {
		if time.Now().After(item.GetExpiresAt()) {
			continue
		}
		if matchFilter(item, filter) {
			if count >= offset {
				copy := item
				results = append(results, &copy)
				if limit > 0 && len(results) >= limit {
					break
				}
			}
			count++
		}
	}

	return results, nil
}
func matchFilter[T any](item T, filter map[string]any) bool {
	v := reflect.ValueOf(item)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for key, val := range filter {
		field := v.FieldByName(key)
		if !field.IsValid() {
			return false
		}
		if !reflect.DeepEqual(field.Interface(), val) {
			return false
		}
	}
	return true
}

func (m *MemoryStore[T]) Delete(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, id)
	return nil
}

func (s *MemoryStore[T]) cleanupExpired() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for id, item := range s.data {
		if item.GetExpiresAt().Before(now) {
			delete(s.data, id)
		}
	}
	return nil
}
