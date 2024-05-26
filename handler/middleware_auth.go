package handler

import (
	"fmt"
	"net/http"

	"github.com/LostArrows27/go-rss-aggregator/internal/auth"

	"github.com/LostArrows27/go-rss-aggregator/internal/database"
)

// 1. authHandler = function need to be called after auth
func (cfg *ApiConfig) MiddlewareAuth(authHandler AuthedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)

		if err != nil {
			RespondWithError(w, 403, fmt.Sprintf("Auth error: %v", err))
			return
		}

		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)

		if err != nil {
			RespondWithError(w, 400, fmt.Sprintf("Error getting user: %v", err))
			return
		}

		authHandler(w, r, user)
	}
}

type AuthedHandler func(http.ResponseWriter, *http.Request, database.User)
