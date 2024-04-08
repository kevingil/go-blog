package store

import (
	"errors"
	"sync"
)

// Store is a key-value store implementation.
// Extends the sync.Map package with a Redis like interface.
// Beware of memory limitations !!!
type Store struct {
	value sync.Map
}

// Store defines the interface for a key-value store.
type StoreInterface interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Delete(key string) error
}

// NewClient initializes a new Store instance.
func NewClient() *Store {
	return &Store{}
}

// Get returns the value associated with the given key or an error.
func (s *Store) Get(key string) ([]byte, error) {
	if value, ok := s.value.Load(key); ok {
		return value.([]byte), nil
	}
	return nil, errors.New("key not found")
}

// Set adds key-value pair and returns an error.
func (s *Store) Set(key string, value []byte) error {
	if len(key) == 0 {
		return errors.New("key can't be empty string")
	}
	s.value.Store(key, value)
	return nil
}

// Delete removes the key-value pair
func (s *Store) Delete(key string) error {
	s.value.Delete(key)
	return nil
}
