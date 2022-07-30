package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/DarioRoman01/cqrs/helpers"
	"github.com/DarioRoman01/cqrs/models"
	"github.com/DarioRoman01/cqrs/repository"
	"github.com/golang-jwt/jwt/v4"
)

type loginResponse struct {
	Token string `json:"token"`
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	var req models.SignupLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := repository.Login(r.Context(), req.Name, req.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if ok, _ := helpers.ComparePasswords(user.Password, req.Password); !ok {
		http.Error(w, "invalid password", http.StatusUnauthorized)
		return
	}

	claims := &models.Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 2)),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(cfg.JWTSecret)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse{Token: token})
}

func logoutUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("Authorization")
	if userID == "" {
		http.Error(w, "user id not found", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "logout successful"})
}
