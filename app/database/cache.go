package database

import (
	"log"

	"github.com/kevingil/blog/app/pkg/sider"
)

var cache *sider.Sider

func Cache() {
	// Testing cache in memory database
	cache = sider.NewClient()

	var err error
	err = cache.Set("hello", []byte("world"))
	if err != nil {
		log.Fatalf("Failed to set value: %v", err)
	}

	value, err := cache.Get("hello")
	if err != nil {
		log.Fatalf("Failed to get value: %v", err)
	}

	log.Printf("value test, hello: %s\n", value)
}
