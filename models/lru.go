package models

import (
	"log"

	lru "github.com/flyaways/golang-lru"
)

func LruInit() *lru.Cache {
	cache, err := lru.New(128)

	if err != nil {
		log.Panicln(err)
	}

	return cache
}
