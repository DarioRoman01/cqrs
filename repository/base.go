package repository

import "context"

// Repository is the generic interface for a repository
type Repository[T any] interface {
	// Close closes the repository
	Close() error
	// Insert inserts a new item into the repository
	Insert(ctx context.Context, entity T) error
	// List returns all items in the repository
	List(ctx context.Context) ([]T, error)
	// Get returns an item from the repository
	Get(ctx context.Context, id string) (T, error)
	// Delete deletes an item from the repository
	Delete(ctx context.Context, id string) error
	// Update updates an item in the repository
	Update(ctx context.Context, entity T) error
}
