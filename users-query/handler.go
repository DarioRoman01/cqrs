package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/DarioRoman01/cqrs/cache"
	"github.com/DarioRoman01/cqrs/repository"
	"github.com/gorilla/mux"
)

func listUsers(w http.ResponseWriter, r *http.Request) {
	cachedUsers, err := cache.Get("users")
	if err != nil {
		fmt.Printf("error getting users from cache: %v", err)
	}

	if cachedUsers != nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cachedUsers)
		return
	}

	users, err := repository.ListUsers(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := cache.Set("users", users); err != nil {
		log.Printf("error setting users in cache: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		http.Error(w, "user id not found", http.StatusBadRequest)
		return
	}

	cachedUser, err := cache.Get("user_" + userID)
	if err != nil {
		log.Printf("error getting user from cache: %v", err)
	}

	if cachedUser != nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cachedUser)
		return
	}

	user, err := repository.GetUser(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := cache.Set("user_"+userID, user); err != nil {
		log.Printf("error setting user in cache: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
