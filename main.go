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

	conn, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal("Can't connect to database", err)
	}

	// 2. config server
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

	// 3. config router

	apiConfig := handler.ApiConfig{
		DB: database.New(conn),
	}

	v1Router := chi.NewRouter()
	v1Router.HandleFunc("/", handler.HandlerReadiness)
	v1Router.HandleFunc("/users", apiConfig.HandlerCreateUser)

	router.HandleFunc("/", handler.HandlerReadiness)
	router.Mount("/v1", v1Router)

	// 4. set up server to run
	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	fmt.Printf("Server starting on port %v", portString)

	// server run non-stop from here
	serverError := srv.ListenAndServe()

	if serverError != nil {
		log.Fatal(serverError)
	}

}
