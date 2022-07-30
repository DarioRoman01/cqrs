package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/DarioRoman01/cqrs/cache"
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

	if err := cache.Delete("feeds"); err != nil {
		log.Printf("error deleting cache: %v", err)
	}

	if err := cache.Set(feed.ID, feed); err != nil {
		log.Printf("error setting cache: %v", err)
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

	if err := cache.Delete("feeds"); err != nil {
		log.Printf("error deleting cache: %v", err)
	}

	if err := cache.Delete(feed.ID); err != nil {
		log.Printf("error deleting cache: %v", err)
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

	if err := cache.Delete("feeds"); err != nil {
		log.Printf("error deleting cache: %v", err)
	}

	if err := cache.Delete(feed.ID); err != nil {
		log.Printf("error deleting cache: %v", err)
	}

	if err := cache.Set(feed.ID, feed); err != nil {
		log.Printf("error setting cache: %v", err)
	}
}

// listFeedHandler manage the request to list all the feeds
func listFeedsHandler(w http.ResponseWriter, r *http.Request) {
	cachedFeeds, err := cache.Get("feeds")
	if err != nil {
		log.Printf("error getting cache: %v", err)
	}

	if cachedFeeds != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cachedFeeds)
		return
	}

	feeds, err := repository.ListFeeds(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := cache.Set("feeds", feeds); err != nil {
		log.Printf("error setting cache: %v", err)
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

	cachedFeed, err := cache.Get(id)
	if err != nil {
		log.Printf("error getting cache: %v", err)
	}

	if cachedFeed != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(cachedFeed)
		return
	}

	feed, err := repository.GetFeed(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := cache.Set(id, feed); err != nil {
		log.Printf("error setting cache: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feed)
}
