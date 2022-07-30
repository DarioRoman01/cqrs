package main

import (
	"fmt"
	"net/http"

	"github.com/DarioRoman01/cqrs/database"
	"github.com/DarioRoman01/cqrs/events"
	"github.com/DarioRoman01/cqrs/middlewares"
	"github.com/DarioRoman01/cqrs/repository"
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
)

// Config is the configuration for the application
type Config struct {
	// PostgresUser is the postgres user
	PostgresDB string `envconfig:"POSTGRES_DB"`
	// PostgresUser is the postgres user
	PostgresUser string `envconfig:"POSTGRES_USER"`
	// PostgresPassword is the postgres password
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
	// NatsAddress is the nats address
	NatsAddress string `envconfig:"NATS_ADDRESS"`
	// JWTSecret is the jwt secret
	JWTSecret string `envconfig:"JWT_SECRET"`
}

// newRouter creates a new router
func newRouter(cfg Config) *mux.Router {
	router := mux.NewRouter()
	router.Use(middlewares.CheckAuthMiddleware(cfg.JWTSecret))
	router.HandleFunc("/feeds/create", createFeedHandler).Methods("POST")
	router.HandleFunc("/feeds/update", updateFeedHandler).Methods("PUT")
	router.HandleFunc("/feeds/{id}", deleteFeedHandler).Methods("DELETE")
	return router
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to process env config: %s", err))
	}

	addr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
	postgresRepo, err := database.NewFeedRepository(addr)
	if err != nil {
		panic(fmt.Sprintf("failed to create postgres repository: %s", err))
	}

	repository.SetFeedRepository(postgresRepo)
	eventStore, err := events.NewNatsEventStore(fmt.Sprintf("nats://%s", cfg.NatsAddress))
	if err != nil {
		panic(fmt.Sprintf("failed to create nats event store: %s", err))
	}

	events.SetEventStore(eventStore)
	defer events.Close()
	defer repository.CloseFeedRepo()

	router := newRouter(cfg)
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(fmt.Sprintf("failed to start server: %s", err))
	}
}
