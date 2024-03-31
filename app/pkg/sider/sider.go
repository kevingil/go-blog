package sider

import (
	"errors"
	"sync"
)

// Store defines the interface for a key-value store.
type Store interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
	Delete(key string) error
}

// Sider is a key-value store implementation.
// Extends the sync.Map package with a Redis like interface.
type Sider struct {
	value sync.Map
}

// NewClient initializes a new Sider instance.
func NewClient() *Sider {
	return &Sider{}
}

// Get returns the value associated with the given key or an error.
func (s *Sider) Get(key string) ([]byte, error) {
	if value, ok := s.value.Load(key); ok {
		return value.([]byte), nil
	}
	return nil, errors.New("key not found")
}

// Set adds key-value pair and returns an error.
func (s *Sider) Set(key string, value []byte) error {
	if len(key) == 0 {
		return errors.New("key can't be empty string")
	}
	s.value.Store(key, value)
	return nil
}

// Delete removes the key-value pair
func (s *Sider) Delete(key string) error {
	s.value.Delete(key)
	return nil
}
