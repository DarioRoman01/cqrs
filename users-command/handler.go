package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/DarioRoman01/cqrs/helpers"
	"github.com/DarioRoman01/cqrs/models"
	"github.com/DarioRoman01/cqrs/repository"
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
