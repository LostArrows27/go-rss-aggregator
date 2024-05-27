package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/LostArrows27/go-rss-aggregator/env"
	"github.com/LostArrows27/go-rss-aggregator/handler"
	"github.com/LostArrows27/go-rss-aggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
)

func main() {
	// 1. connect to database

	env.LoadEnv(".env")

	portString := env.GetEnv("PORT")

	dbURL := env.GetEnv("DB_URL")

	connection, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("Can't connect to database", err)
	}

	// 2. config server HTTPS
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS",
		},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// 3. config database

	apiConfig := handler.ApiConfig{
		DB: database.New(connection),
	}

	// 4. config router handler

	v1Router := chi.NewRouter()

	v1Router.Get("/", handler.HandlerReadiness)

	v1Router.Route("/users", func(r chi.Router) {
		r.Get("/", apiConfig.MiddlewareAuth(apiConfig.HandlerGetUserByAPIKey))
		r.Post("/create", apiConfig.HandlerCreateUser)
	})

	v1Router.Route("/feeds", func(r chi.Router) {
		r.Get("/all", apiConfig.GetAllFeed)
		r.Post("/create", apiConfig.MiddlewareAuth(apiConfig.HandlerCreateFeed))
	})

	v1Router.Route("/feed_follow", func(r chi.Router) {
		r.Get("/all", apiConfig.MiddlewareAuth(apiConfig.GetAllFeedFollows))
		r.Post("/create", apiConfig.MiddlewareAuth(apiConfig.HandlerCreateFeedFollow))
		r.Delete("/delete/{feed_follow_id}", apiConfig.MiddlewareAuth(apiConfig.DeleteFeedFollowByID))
	})

	v1Router.Post("/feed_follow/create", apiConfig.MiddlewareAuth(apiConfig.HandlerCreateFeedFollow))

	router.Get("/", handler.HandlerReadiness)
	router.Mount("/v1", v1Router)

	// 5. set up server to run
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	fmt.Printf("Server starting on port %v", portString)

	// NOTE: server run non-stop from here
	serverError := srv.ListenAndServe()

	if serverError != nil {
		log.Fatal(serverError)
	}

}
