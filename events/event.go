package events

import (
	"context"

	"github.com/DarioRoman01/cqrs/models"
)

// EventStore is an interface that defines the methods that a repository must implement
type EventStore interface {
	// Close closes the repository
	Close()
	// PublishCreatedFeed publishes a new created feed event
	PublishCreatedFeed(ctx context.Context, feed *models.Feed) error
	// SubscribeCreatedFeed subscribes to the created feed event
	SubscribeCreatedFeed(ctx context.Context) (<-chan *CreatedFeedMessage, error)
	// OnCreatedFeed registers a handler for the created feed event
	OnCreatedFeed(handler func(*CreatedFeedMessage)) error
}

var eventStore EventStore

// SetEventStore sets the event store
func SetEventStore(e EventStore) {
	eventStore = e
}

// Close closes the event store
func Close() {
	eventStore.Close()
}

// PublishCreatedFeed publishes a new created feed event
func PublishCreatedFeed(ctx context.Context, feed *models.Feed) error {
	return eventStore.PublishCreatedFeed(ctx, feed)
}

// SubscribeCreatedFeed subscribes to the created feed event
func SubscribeCreatedFeed(ctx context.Context) (<-chan *CreatedFeedMessage, error) {
	return eventStore.SubscribeCreatedFeed(ctx)
}

// OnCreatedFeed registers a handler for the created feed event
func OnCreatedFeed(handler func(*CreatedFeedMessage)) error {
	return eventStore.OnCreatedFeed(handler)
}
