package main

import (
	"fmt"
	"net/http"

	"github.com/DarioRoman01/cqrs/events"
	"github.com/kelseyhightower/envconfig"
)

// Config is the configuration for the application
type Config struct {
	// NatsAddress is the nats address
	NatsAddress string `envconfig:"NATS_ADDRESS"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to process env config: %s", err))
	}

	hub := NewHub()

	eventStore, err := events.NewNatsEventStore(fmt.Sprintf("nats://%s", cfg.NatsAddress))
	if err != nil {
		panic(fmt.Sprintf("failed to create nats event store: %s", err))
	}

	events.SetEventStore(eventStore)
	err = events.OnCreatedFeed(func(msg *events.CreatedFeedMessage) {
		hub.Broadcast(newCreatedFeedMessage(msg.ID, msg.Title, msg.Description, msg.UserID, msg.CreatedAt, Create.String()), nil)
	})

	if err != nil {
		panic(fmt.Sprintf("failed to subscribe to CreatedFeed event: %s", err))
	}

	err = events.OnUpdatedFeed(func(msg *events.UpdatedFeedMessage) {
		hub.Broadcast(newUpdatedFeedMessage(msg.ID, msg.Title, msg.Description, msg.CreatedAt, Update.String()), nil)
	})

	if err != nil {
		panic(fmt.Sprintf("failed to subscribe to UpdatedFeed event: %s", err))
	}

	err = events.OnDeletedFeed(func(msg *events.DeletedFeedMessage) {
		hub.Broadcast(newDeletedFeedMessage(msg.ID, Delete.String()), nil)
	})

	if err != nil {
		panic(fmt.Sprintf("failed to subscribe to DeletedFeed event: %s", err))
	}

	defer events.Close()
	go hub.Run()

	http.HandleFunc("/ws", hub.HandleWebSocket)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(fmt.Sprintf("failed to start server: %s", err))
	}
}
