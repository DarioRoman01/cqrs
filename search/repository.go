package search

import (
	"context"

	"github.com/DarioRoman01/cqrs/models"
)

// SearchRepository is an interface that defines the methods that a repository must implement
type SearchRepository interface {
	// Close closes the repository
	Close()
	// IndexFeed indexes a feed
	IndexFeed(ctx context.Context, feed *models.Feed) error
	// SearchFeeds searches for feeds
	SearchFeed(ctx context.Context, query string) ([]*models.Feed, error)
	// UpdateIndex updates the index
	UpdateIndex(ctx context.Context, feed *models.Feed) error
	// UnindexFeed unindexes a feed
	UnIndexFeed(ctx context.Context, feed *models.Feed) error
}

var repo SearchRepository

// SetSearchRepository sets the search repository
func SetSearchRepository(r SearchRepository) {
	repo = r
}

// Close closes the search repository
func Close() {
	repo.Close()
}

// IndexFeed indexes a feed
func IndexFeed(ctx context.Context, feed *models.Feed) error {
	return repo.IndexFeed(ctx, feed)
}

// SearchFeeds searches for feeds
func SearchFeed(ctx context.Context, query string) ([]*models.Feed, error) {
	return repo.SearchFeed(ctx, query)
}

// UnindexFeed unindexes a feed
func UnIndexFeed(ctx context.Context, feed *models.Feed) error {
	return repo.UnIndexFeed(ctx, feed)
}

// UpdateIndex updates the index
func UpdateIndex(ctx context.Context, feed *models.Feed) error {
	return repo.UpdateIndex(ctx, feed)
}
