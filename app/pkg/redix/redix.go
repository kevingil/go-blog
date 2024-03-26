package redix

import (
	"errors"
	"sync"
)

// Redix is an in memory key value store service
// can hold multiple tables
type Redix struct {
	tables map[string]*Table
	mutex  *sync.RWMutex
}

// Table holds the key value store
type Table struct {
	Name  string
	data  map[string][]byte
	mutex *sync.RWMutex
}

// New initializes a new Redix servie.
func New() *Redix {
	return &Redix{
		tables: make(map[string]*Table),
		mutex:  &sync.RWMutex{},
	}
}

// CreateTable creates a new table and returns it
func (kv *Redix) CreateTable(name string) (*Table, error) {
	kv.mutex.Lock()
	defer kv.mutex.Unlock()

	if _, exists := kv.tables[name]; exists {
		return nil, errors.New("table already exists")
	}

	newTable := &Table{
		Name:  name, // Set the name of the table
		data:  make(map[string][]byte),
		mutex: &sync.RWMutex{},
	}

	kv.tables[name] = newTable
	return newTable, nil
}

// GetTable retrieves a table by name.
func (kv *Redix) GetTable(name string) (*Table, error) {
	kv.mutex.RLock()
	defer kv.mutex.RUnlock()

	table, found := kv.tables[name]
	if !found {
		return nil, errors.New("table not found")
	}

	return table, nil
}

// Table operations

// Get retrieves a value by key from the table.
func (t *Table) Get(key string) ([]byte, error) {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	value, found := t.data[key]
	if !found {
		return nil, errors.New("key not found")
	}

	return value, nil
}

// Set stores the given value under the specified key in the table.
func (t *Table) Set(key string, value []byte) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.data[key] = value
	return nil
}

// Delete removes the key-value pair associated with the given key from the table.
func (t *Table) Delete(key string) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	delete(t.data, key)
	return nil
}
