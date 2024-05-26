package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/LostArrows27/go-rss-aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) GetAllFeed(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeeds(r.Context())

	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Couldn't get feeds: %v", err))
		return
	}

	RespondWithJSON(w, 200, convertFeedArrayJSON(feeds))

}

func (apiCfg *ApiConfig) HandlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {

	// 1. get body from request
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	// 2. create feed
	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.URL,
		UserID:    user.ID,
	})

	// 3. handle response
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	RespondWithJSON(w, 200, convertFeedJSON(feed))
}
