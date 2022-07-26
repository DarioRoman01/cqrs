package search

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

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

// IndexFeed indexes the given feed
func (e *ElasticSearchRepository) IndexFeed(ctx context.Context, feed *models.Feed) error {
	body, _ := json.Marshal(feed)
	_, err := e.client.Index(
		"feeds",
		bytes.NewReader(body),
		e.client.Index.WithDocumentID(feed.ID),
		e.client.Index.WithContext(ctx),
		e.client.Index.WithRefresh("wait_for"),
	)

	return err

}

// UnIndexFeed unindexes the given feed
func (e *ElasticSearchRepository) UnIndexFeed(ctx context.Context, feed *models.Feed) error {
	_, err := e.client.Delete("feeds", feed.ID, e.client.Delete.WithContext(ctx), e.client.Delete.WithRefresh("wait_for"))
	return err
}

// UpdateIndex updates the index
func (e *ElasticSearchRepository) UpdateIndex(ctx context.Context, feed *models.Feed) error {
	body, err := json.Marshal(feed)
	if err != nil {
		return err
	}

	_, err = e.client.Update(
		"feeds",
		feed.ID,
		bytes.NewReader(body),
		e.client.Update.WithContext(ctx),
		e.client.Update.WithRefresh("wait_for"),
	)

	return err
}

// SearchFeeds searches for feeds with the given query
func (e *ElasticSearchRepository) SearchFeed(ctx context.Context, query string) ([]*models.Feed, error) {
	buff := new(bytes.Buffer)
	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":            query,
				"fields":           []string{"title", "description"},
				"fuzziness":        3,
				"cutoff_frequency": 0.0001,
			},
		},
	}

	if err := json.NewEncoder(buff).Encode(searchQuery); err != nil {
		return nil, err
	}

	res, err := e.client.Search(
		e.client.Search.WithContext(ctx),
		e.client.Search.WithIndex("feeds"),
		e.client.Search.WithBody(buff),
		e.client.Search.WithTrackTotalHits(true),
	)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.IsError() {
		return nil, errors.New(res.String())
	}

	var eRes map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&eRes); err != nil {
		return nil, err
	}

	var feeds []*models.Feed
	for _, hit := range eRes["hits"].(map[string]interface{})["hits"].([]interface{}) {
		var feed models.Feed
		if err := json.Unmarshal([]byte(hit.(map[string]interface{})["_source"].(string)), &feed); err != nil {
			return nil, err
		}

		// var feed models.Feed
		// source := hit.(map[string]interface{})["_source"]
		// marshal, err := json.Marshal(source)
		// if err != nil {
		// 	return nil, err
		// }

		// if err := json.Unmarshal(marshal, &feed); err != nil {
		// 	return nil, err
		// }

		// feeds = append(feeds, &feed)
		feeds = append(feeds, &feed)
	}

	return feeds, nil
}
