package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/DarioRoman01/cqrs/events"
	"github.com/DarioRoman01/cqrs/helpers"
	"github.com/DarioRoman01/cqrs/models"
	"github.com/DarioRoman01/cqrs/repository"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
)

func register(w http.ResponseWriter, r *http.Request) {
	var req models.SignupLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user := &models.User{
		ID:        ksuid.New().String(),
		Name:      req.Name,
		Password:  req.Password,
		CreatedAt: time.Now().UTC(),
	}

	cfg := helpers.GetPasswordConfig()
	hashPwd, err := helpers.GeneratePassword(cfg, user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.Password = hashPwd
	err = repository.InsertUser(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	user, err := repository.GetUser(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	requestUserId := r.Context().Value("user_id").(string)
	if requestUserId != user.ID {
		http.Error(w, "you are not allowed to delete this user", http.StatusForbidden)
		return
	}

	err = repository.DeleteUser(r.Context(), user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := events.PublishDeletedUser(r.Context(), user); err != nil {
		log.Printf("error publishing deleted user event: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
