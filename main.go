package main

import (
	"legitt-backend/model"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	model.InitRedis()

	model.InitDb()

	// Register routes
	router := RegisterRoutes()

	print("Server running on port 8080")

	s := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	err := s.ListenAndServe()
	if err != nil {
		os.Exit(1)
	}

}
