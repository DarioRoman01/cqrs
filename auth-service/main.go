package main

import (
	"fmt"
	"net/http"

	"github.com/DarioRoman01/cqrs/database"
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
	// MemcacheAddress is the memcache address
	MemCacheAddress string `envconfig:"MEMCACHE_ADDRESS"`
	// JWT_SECRET is the jwt secret
	JWTSecret string `envconfig:"JWT_SECRET"`
}

func newRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("auth/login", loginUser).Methods("POST")
	router.HandleFunc("auth/logout", logoutUser).Methods("GET")
	return router
}

var cfg Config

func init() {
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to process env config: %s", err))
	}
}

func main() {
	addr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
	userRepo, err := database.NewUserRepository(addr)
	if err != nil {
		panic(fmt.Sprintf("failed to create postgres repository: %s", err))
	}

	repository.SetUserRepository(userRepo)

	defer userRepo.Close()

	router := newRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(fmt.Sprintf("failed to start server: %s", err))
	}
}
