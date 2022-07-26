package search

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/DarioRoman01/cqrs/models"
	elastic "github.com/elastic/go-elasticsearch/v8"
)

// ElasticSearchRepository is an implementation of the SearchRepository interface that uses ElasticSearch
type ElasticSearchRepository struct {
	// client is the ElasticSearch client
	client *elastic.Client
}

// NewElasticRepo creates a new ElasticSearchRepository
func NewElasticRepo(url string) (*ElasticSearchRepository, error) {
	client, err := elastic.NewClient(elastic.Config{
		Addresses: []string{url},
	})

	if err != nil {
		return nil, err
	}

	return &ElasticSearchRepository{client: client}, nil
}

func (e *ElasticSearchRepository) Close() {
	// TODO implement Close method for ElasticSearchRepository
}

func (e *ElasticSearchRepository) IndexFeed(ctx context.Context, feed *models.Feed) error {
	body, _ := json.Marshal(feed)
	_, err := e.client.Index(
		"feeds",
		bytes.NewReader(body),
		e.client.Index.WithDocumentID(feed.ID),
		e.client.Index.WithContext(ctx),
	)

	return err

}
