package database

import (
	"testing"

	"github.com/kevingil/blog/app/pkg/sider"
)

func TestCache(t *testing.T) {
	// Testing cache in memory database
	cache := sider.NewClient()

	var err error
	err = cache.Set("hello", []byte("world"))
	if err != nil {
		t.Fatalf("Failed to set value: %v", err)
	}

	value, err := cache.Get("hello")
	if err != nil {
		t.Fatalf("Failed to get value: %v", err)
	}

	if string(value) != "world" {
		t.Errorf("Unexpected value retrieved from cache: %s", value)
	}
}
