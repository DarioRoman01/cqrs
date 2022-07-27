package main

import (
	"fmt"
	"net/http"

	"github.com/DarioRoman01/cqrs/database"
	"github.com/DarioRoman01/cqrs/events"
	"github.com/DarioRoman01/cqrs/repository"
	"github.com/DarioRoman01/cqrs/search"
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
	// ElasticSearchAddress is the elastic search address
	ElasticSearchAddress string `envconfig:"ELASTICSEARCH_ADDRESS"`
}

func newRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/feeds", listFeedsHandler).Methods("GET")
	router.HandleFunc("/feeds/{id}", getFeedHandler).Methods("GET")
	router.HandleFunc("/search", searchFeedsHandler).Methods("GET")
	return router
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to process env config: %s", err))
	}

	addr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
	postgresRepo, err := database.NewPostgresRepository(addr)
	if err != nil {
		panic(fmt.Sprintf("failed to create postgres repository: %s", err))
	}

	repository.SetFeedRepository(postgresRepo)
	elasticRepo, err := search.NewElasticRepo(fmt.Sprintf("http://%s", cfg.ElasticSearchAddress))
	if err != nil {
		panic(fmt.Sprintf("failed to create elastic search repository: %s", err))
	}

	search.SetSearchRepository(elasticRepo)

	eventStore, err := events.NewNatsEventStore(fmt.Sprintf("nats://%s", cfg.NatsAddress))
	if err != nil {
		panic(fmt.Sprintf("failed to create nats event store: %s", err))
	}

	events.SetEventStore(eventStore)
	defer events.Close()
	defer repository.Close()
	defer search.Close()

	if err = eventStore.OnCreatedFeed(onCreatedFeed); err != nil {
		panic(fmt.Sprintf("failed to subscribe to event: %s", err))
	}

	if err = eventStore.OnDeletedFeed(onDeletedFeed); err != nil {
		panic(fmt.Sprintf("failed to subscribe to event: %s", err))
	}

	if err = eventStore.OnUpdatedFeed(onUpdatedFeed); err != nil {
		panic(fmt.Sprintf("failed to subscribe to event: %s", err))
	}

	router := newRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(fmt.Sprintf("failed to start server: %s", err))
	}
}
