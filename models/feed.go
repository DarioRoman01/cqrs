package models

import "time"

// Feed is a feed of articles
type Feed struct {
	// ID is the primary key
	ID string `json:"id"`
	// Title is the title of the feed
	Title string `json:"title"`
	// description is the description of the feed
	Description string `json:"description"`
	// CreatedAt is the date the feed was created
	CreatedAt time.Time `json:"created_at"`
}
