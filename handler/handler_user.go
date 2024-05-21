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

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		RespondWithError(w, 400, fmt.Sprintf("Couldn't create user: %v", err))
		return
	}

	RespondWithJSON(w, 200, convertUserJSON(user))
}
