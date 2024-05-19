package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	// 1. configure port
	fmt.Println("Hello, World!")

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error in loading .env")
	}

	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("Not found port")
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

	v1Router := chi.NewRouter()
	v1Router.HandleFunc("/", handlerErr)
	router.HandleFunc("/", handlerReadiness)
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
		log.Fatal(err)
	}

}
