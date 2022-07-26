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

func Close() error {
	return repository.Close()
}

func Insert(ctx context.Context, feed *models.Feed) error {
	return repository.Insert(ctx, feed)
}

func List(ctx context.Context) ([]*models.Feed, error) {
	return repository.List(ctx)
}

func Get(ctx context.Context, id string) (*models.Feed, error) {
	return repository.Get(ctx, id)
}

func Delete(ctx context.Context, id string) error {
	return repository.Delete(ctx, id)
}

func Update(ctx context.Context, feed *models.Feed) error {
	return repository.Update(ctx, feed)
}
