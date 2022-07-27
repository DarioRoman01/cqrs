package main

import "time"

type MessgaType int

const (
	Create MessgaType = iota
	Update
	Delete
)

func (m MessgaType) String() string {
	switch m {
	case Create:
		return "create"
	case Update:
		return "update"
	case Delete:
		return "delete"
	default:
		return "unknown"
	}
}

// CreatedFeedMessage is the message that is published when a feed is created
type CreatedFeedMessage struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Type        string    `json:"type"`
}

// UpdatedFeedMessage is the message that is published when a feed is updated
type UpdatedFeedMessage struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Type        string    `json:"type"`
}

// DeletedFeedMessage is the message that is published when a feed is deleted
type DeletedFeedMessage struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

// newCreatedFeedMessage creates a new CreatedFeedMessage
func newCreatedFeedMessage(id, title, description string, createdAt time.Time, t string) *CreatedFeedMessage {
	return &CreatedFeedMessage{
		ID:          id,
		Title:       title,
		Description: description,
		CreatedAt:   createdAt,
		Type:        t,
	}
}

func newUpdatedFeedMessage(id, title, description string, createdAt time.Time, t string) *UpdatedFeedMessage {
	return &UpdatedFeedMessage{
		ID:          id,
		Title:       title,
		Description: description,
		CreatedAt:   createdAt,
		Type:        t,
	}
}

func newDeletedFeedMessage(id string, t string) *DeletedFeedMessage {
	return &DeletedFeedMessage{
		ID:   id,
		Type: t,
	}
}
