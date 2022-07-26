package repository

import (
	"context"

	"github.com/DarioRoman01/cqrs/models"
)

// Repository is the interface that wraps the basic CRUD operations
type Repository interface {
	// Close closes the repository
	Close() error
	// InsertFeed inserts a new feed
	InsertFeed(ctx context.Context, feed *models.Feed) error
	// ListFeeds lists all feeds in the repository
	ListFeeds(ctx context.Context) ([]*models.Feed, error)
	// GetFeed gets a feed by ID
	GetFeed(ctx context.Context, id string) (*models.Feed, error)
	// DeleteFeed deletes a feed by ID
	DeleteFeed(ctx context.Context, id string) error
	// UpdateFeed updates a feed
	UpdateFeed(ctx context.Context, feed *models.Feed) error
}

var repository Repository

func SetRepository(r Repository) {
	repository = r
}

func Close() error {
	return repository.Close()
}

func InsertFeed(ctx context.Context, feed *models.Feed) error {
	return repository.InsertFeed(ctx, feed)
}

func ListFeeds(ctx context.Context) ([]*models.Feed, error) {
	return repository.ListFeeds(ctx)
}

func GetFeed(ctx context.Context, id string) (*models.Feed, error) {
	return repository.GetFeed(ctx, id)
}

func DeleteFeed(ctx context.Context, id string) error {
	return repository.DeleteFeed(ctx, id)
}

func UpdateFeed(ctx context.Context, feed *models.Feed) error {
	return repository.UpdateFeed(ctx, feed)
}
