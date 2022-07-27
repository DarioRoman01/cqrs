package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/DarioRoman01/cqrs/events"
	"github.com/DarioRoman01/cqrs/models"
	"github.com/DarioRoman01/cqrs/repository"
	"github.com/DarioRoman01/cqrs/search"
	"github.com/gorilla/mux"
)

// onCreatedFeed is the handler that is called when a feed is created
func onCreatedFeed(msg *events.CreatedFeedMessage) {
	feed := models.Feed{
		ID:          msg.ID,
		Title:       msg.Title,
		Description: msg.Description,
		CreatedAt:   msg.CreatedAt,
	}

	if err := search.IndexFeed(context.Background(), &feed); err != nil {
		log.Printf("error indexing feed: %v", err)
	}
}

// onDeletedFeed is the handler that is called when a feed is deleted
func onDeletedFeed(msg *events.DeletedFeedMessage) {
	feed := models.Feed{
		ID: msg.ID,
	}

	if err := search.UnIndexFeed(context.Background(), &feed); err != nil {
		log.Printf("error deleting feed: %v", err)
	}
}

// onUpdatedFeed is the handler that is called when a feed is updated
func onUpdatedFeed(msg *events.UpdatedFeedMessage) {
	feed := models.Feed{
		ID:          msg.ID,
		Title:       msg.Title,
		Description: msg.Description,
		CreatedAt:   msg.CreatedAt,
	}

	if err := search.UpdateIndex(context.Background(), &feed); err != nil {
		log.Printf("error updating feed: %v", err)
	}
}

// listFeedHandler manage the request to list all the feeds
func listFeedsHandler(w http.ResponseWriter, r *http.Request) {
	feeds, err := repository.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feeds)
}

// searchFeedHandler manage the request to search for feeds
func searchFeedsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "query parameter is required", http.StatusBadRequest)
		return
	}

	feeds, err := search.SearchFeed(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feeds)
}

// getFeedHandler manage the request to get a feed
func getFeedHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		http.Error(w, "id parameter is required", http.StatusBadRequest)
		return
	}

	feed, err := repository.Get(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feed)
}
