package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/LostArrows27/go-rss-aggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *ApiConfig) GetAllFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feed_follows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)

	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Couldn't get feed follow list: %v", err))
		return
	}

	RespondWithJSON(w, 201, convertFeedFollowArrayJSON(feed_follows))

}

func (apiCfg *ApiConfig) DeleteFeedFollowByID(w http.ResponseWriter, r *http.Request, user database.User) {
	feed_follow_str := chi.URLParam(r, "feed_follow_id")

	feedFollowID, err := uuid.Parse(feed_follow_str)

	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Error parsing UUID: %v", err))
		return
	}

	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})

	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Couldn't delete feed follow: %v", err))
		return
	}

	RespondWithJSON(w, 200, map[string]string{"message": "Feed follow deleted"})

}

func (apiCfg *ApiConfig) HandlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {

	// 1. get body from request
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	// 2. create feed follow
	feed_follow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})

	// 3. handle response
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	RespondWithJSON(w, 200, convertFeedFollowJSON(feed_follow))
}
