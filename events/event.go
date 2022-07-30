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
	SubscribeCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error)
	// OnCreatedFeed registers a handler for the created feed event
	OnCreatedFeed(handler func(*CreatedFeedMessage)) error
	// PublishUpdatedFeed publishes a new updated feed event
	PublishUpdatedFeed(ctx context.Context, feed *models.Feed) error
	// SubscribeUpdatedFeed subscribes to the updated feed event
	SubscribeUpdatedFeed(ctx context.Context) (<-chan UpdatedFeedMessage, error)
	// OnUpdatedFeed registers a handler for the updated feed event
	OnUpdatedFeed(handler func(*UpdatedFeedMessage)) error
	// PublishDeletedFeed publishes a new deleted feed event
	PublishDeletedFeed(ctx context.Context, feed *models.Feed) error
	// SubscribeDeletedFeed subscribes to the deleted feed event
	SubscribeDeletedFeed(ctx context.Context) (<-chan DeletedFeedMessage, error)
	// OnDeletedFeed registers a handler for the deleted feed event
	OnDeletedFeed(handler func(*DeletedFeedMessage)) error
	// PublishDeletedUser publishes a new deleted user event
	PublishDeletedUser(ctx context.Context, user *models.User) error
	// SubscribeDeletedUser subscribes to the deleted user event
	SubscribeDeletedUser(ctx context.Context) (<-chan DeletedUserMessage, error)
	// OnDeletedUser registers a handler for the deleted user event
	OnDeletedUser(handler func(*DeletedUserMessage)) error
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
func SubscribeCreatedFeed(ctx context.Context) (<-chan CreatedFeedMessage, error) {
	return eventStore.SubscribeCreatedFeed(ctx)
}

// OnCreatedFeed registers a handler for the created feed event
func OnCreatedFeed(handler func(*CreatedFeedMessage)) error {
	return eventStore.OnCreatedFeed(handler)
}

//PublishUpdatedFeed publishes a new updated feed event
func PublishUpdatedFeed(ctx context.Context, feed *models.Feed) error {
	return eventStore.PublishUpdatedFeed(ctx, feed)
}

//SubscribeUpdatedFeed subscribes to the updated feed event
func SubscribeUpdatedFeed(ctx context.Context) (<-chan UpdatedFeedMessage, error) {
	return eventStore.SubscribeUpdatedFeed(ctx)
}

// OnUpdatedFeed registers a handler for the updated feed event
func OnUpdatedFeed(handler func(*UpdatedFeedMessage)) error {
	return eventStore.OnUpdatedFeed(handler)
}

// PublishDeletedFeed publishes a new deleted feed event
func PublishDeletedFeed(ctx context.Context, feed *models.Feed) error {
	return eventStore.PublishDeletedFeed(ctx, feed)
}

//SubscribeDeletedFeed subscribes to the deleted feed event
func SubscribeDeletedFeed(ctx context.Context) (<-chan DeletedFeedMessage, error) {
	return eventStore.SubscribeDeletedFeed(ctx)
}

// OnDeletedFeed registers a handler for the deleted feed event
func OnDeletedFeed(handler func(*DeletedFeedMessage)) error {
	return eventStore.OnDeletedFeed(handler)
}

func PublishDeletedUser(ctx context.Context, user *models.User) error {
	return eventStore.PublishDeletedUser(ctx, user)
}

func SubscribeDeletedUser(ctx context.Context) (<-chan DeletedUserMessage, error) {
	return eventStore.SubscribeDeletedUser(ctx)
}

func OnDeletedUser(handler func(*DeletedUserMessage)) error {
	return eventStore.OnDeletedUser(handler)
}
