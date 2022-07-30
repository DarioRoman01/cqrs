package main

import (
	"encoding/json"
	"net/http"

	"github.com/DarioRoman01/cqrs/cache"
	"github.com/DarioRoman01/cqrs/helpers"
	"github.com/DarioRoman01/cqrs/models"
	"github.com/DarioRoman01/cqrs/repository"
)

func loginUser(w http.ResponseWriter, r *http.Request) {
	var req models.SignupLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := repository.Login(r.Context(), req.Name, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if ok, _ := helpers.ComparePasswords(user.Password, req.Password); !ok {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		return
	}

	if err := cache.Set(user.ID, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func logoutUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Authorization")
	if userID == "" {
		http.Error(w, "user id not found", http.StatusUnauthorized)
		return
	}

	if err := cache.Delete(userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "logout successful"})
}