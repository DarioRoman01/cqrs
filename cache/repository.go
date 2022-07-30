package cache

import "context"

// CacheRepository is the interface for the cache repository
type CacheRepository interface {
	// Close closes the connection of the repository
	Close() error
	// Get returns the an item with the given id from the cache or an error if it does not exist
	Get(ctx context.Context, id string) (interface{}, error)
	// Set sets an item with the given id in the cache
	Set(ctx context.Context, id string, value interface{}) error
	// Delete deletes an item with the given id from the cache
	Delete(ctx context.Context, id string) error
}

var cacheRepo CacheRepository

func SetCacheRepository(r CacheRepository) {
	cacheRepo = r
}

func CloseCacheRepo() error {
	return cacheRepo.Close()
}

func Set(ctx context.Context, id string, value interface{}) error {
	return cacheRepo.Set(ctx, id, value)
}

func Get(ctx context.Context, id string) (interface{}, error) {
	return cacheRepo.Get(ctx, id)
}

func Delete(ctx context.Context, id string) error {
	return cacheRepo.Delete(ctx, id)
}
