package repository

import (
	"context"

	"github.com/DarioRoman01/cqrs/models"
)

type FeedRepository interface {
	Repository[*models.Feed]
}

var repository FeedRepository

func SetFeedRepository(r Repository[*models.Feed]) {
	repository = r
}

func CloseFeedRepo() error {
	return repository.Close()
}

func InsertFeed(ctx context.Context, feed *models.Feed) error {
	return repository.Insert(ctx, feed)
}

func ListFeeds(ctx context.Context) ([]*models.Feed, error) {
	return repository.List(ctx)
}

func GetFeed(ctx context.Context, id string) (*models.Feed, error) {
	return repository.Get(ctx, id)
}

func DeleteFeed(ctx context.Context, id string) error {
	return repository.Delete(ctx, id)
}

func UpdateFeed(ctx context.Context, feed *models.Feed) error {
	return repository.Update(ctx, feed)
}
