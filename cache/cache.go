package cache

import (
	"context"
	"time"
)

// Cache
// this interface give us an abstraction over internal cache we are using
// currently the cache repo contains in-memory cache and redis cache which
// is a de-facto
type Cache interface {
	// Ping
	// ICMP ping for checking connection
	// but be careful in some cases infrastructure filter ICMP packets for prevention of this protocol attacks
	Ping(ctx context.Context) error
	// GetKey
	// to read specific key from cache
	// method :: used for metrics
	// key :: the item which you wish to find in cache
	GetKey(ctx context.Context, method string, key string) (string, error)
	// Set
	// to store specified item with specified key and value in cache
	// method :: used for metrics
	// key, val :: specified key and corresponding value in cache
	Set(ctx context.Context, method string, key string, val string, expiration time.Duration) error

	//RemoveKey
	// to remove specific key from cache
	// method :: used for metrics
	// key :: the item which you wish to remove from the cache
	RemoveKey(ctx context.Context, method string, key string) error
}
