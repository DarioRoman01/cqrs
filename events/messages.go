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
