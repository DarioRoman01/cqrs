package events

import "time"

// Message is the message that is published in the queue
type Message interface {
	// Type returns the type of the message
	Type() string
}

// CreatedFeedMessage is the message that is published when a feed is created
type CreatedFeedMessage struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// Type returns the type of the message
func (CreatedFeedMessage) Type() string {
	return "created_feed"
}

// UpdatedFeedMessage is the message that is published when a feed is updated
type UpdatedFeedMessage struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

// Type returns the type of the message
func (UpdatedFeedMessage) Type() string {
	return "updated_feed"
}

// DeletedFeedMessage is the message that is published when a feed is deleted
type DeletedFeedMessage struct {
	ID        string    `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}

// Type returns the type of the message
func (DeletedFeedMessage) Type() string {
	return "deleted_feed"
}
