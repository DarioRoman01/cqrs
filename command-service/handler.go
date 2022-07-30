package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/DarioRoman01/cqrs/events"
	"github.com/DarioRoman01/cqrs/models"
	"github.com/DarioRoman01/cqrs/repository"
	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
)

// createFeedRequest is the request body for the create feed endpoint
type createFeedRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// createFeedHandler is the handler for the create feed endpoint
func createFeedHandler(w http.ResponseWriter, r *http.Request) {
	var req createFeedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := ksuid.NewRandom()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	feed := &models.Feed{
		ID:          id.String(),
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   time.Now().UTC(),
	}

	if err := repository.InsertFeed(r.Context(), feed); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := events.PublishCreatedFeed(r.Context(), feed); err != nil {
		log.Printf("[ERROR] Failed to publish created feed event: %s", err)
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feed)
}

// listFeedHandler is the handler for the list feed endpoint
func deleteFeedHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	feed, err := repository.GetFeed(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := repository.DeleteFeed(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := events.PublishDeletedFeed(r.Context(), feed); err != nil {
		log.Printf("[ERROR] Failed to publish deleted feed event: %s", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "feed deleted"})
}

// listFeedHandler is the handler for the list feed endpoint
func updateFeedHandler(w http.ResponseWriter, r *http.Request) {
	var req models.Feed
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	feed, err := repository.GetFeed(r.Context(), req.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	feed.Title = req.Title
	feed.Description = req.Description

	if err := repository.UpdateFeed(r.Context(), feed); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := events.PublishUpdatedFeed(r.Context(), feed); err != nil {
		log.Printf("[ERROR] Failed to publish updated feed event: %s", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(feed)
}
