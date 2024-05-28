package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/LostArrows27/go-rss-aggregator/internal/database"
	"github.com/google/uuid"
)

type ApiConfig struct {
	DB *database.Queries
}

func (apiCfg *ApiConfig) HandlerCreateUser(w http.ResponseWriter, r *http.Request) {

	// 1. get body from request
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)

	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Error parsing JSON: %v", err))
		return
	}

	// 2. create user
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	// 3. handle response
	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	RespondWithJSON(w, 200, convertUserJSON(user))
}

func (apiCfg *ApiConfig) HandlerGetUserByAPIKey(w http.ResponseWriter, r *http.Request, user database.User) {
	RespondWithJSON(w, 201, convertUserJSON(user))

}

func (apiCfg *ApiConfig) HanlderGetPostForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  10,
	})

	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Couldn't get posts: %v", err))
		return
	}

	postsJSON := convertDatabasePostArrayToPostArray(posts)

	RespondWithJSON(w, 200, postsJSON)

}
