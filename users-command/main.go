package main

import (
	"fmt"
	"net/http"

	"github.com/DarioRoman01/cqrs/cache"
	"github.com/DarioRoman01/cqrs/database"
	"github.com/DarioRoman01/cqrs/events"
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
	// MemcacheAddress is the memcache address
	MemCacheAddress string `envconfig:"MEMCACHE_ADDRESS"`
}

func newRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/users/register", register).Methods("POST")
	router.HandleFunc("users/delete/{id}", delete).Methods("DELETE")
	return router
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to process env config: %s", err))
	}

	addr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
	userRepo, err := database.NewUserRepository(addr)
	if err != nil {
		panic(fmt.Sprintf("failed to create postgres repository: %s", err))
	}

	repository.SetUserRepository(userRepo)
	memcacheRepo, err := cache.NewCache(cfg.MemCacheAddress)
	if err != nil {
		panic(fmt.Sprintf("failed to create memcache repository: %s", err))
	}

	cache.SetCacheRepository(memcacheRepo)

	natsRepo, err := events.NewNatsEventStore(cfg.NatsAddress)
	if err != nil {
		panic(fmt.Sprintf("failed to create nats repository: %s", err))
	}

	events.SetEventStore(natsRepo)

	defer memcacheRepo.Close()
	defer userRepo.Close()
	defer events.Close()

	router := newRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(fmt.Sprintf("failed to start server: %s", err))
	}
}
