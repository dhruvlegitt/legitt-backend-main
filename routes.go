package main

import (
	"legitt-backend/auth"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func registerAuthRoutes(router *mux.Router) {
	router.HandleFunc("/login", auth.LoginView).Methods(http.MethodPost)
	// router.HandleFunc("/register", RegisterHandler).Methods(http.MethodPost)
	router.HandleFunc("/checkSession", CheckSessionHandler).Methods()
}

func registerProRoutes(router *mux.Router) {
	router.HandleFunc("/registerApp", RegisterAppHandler)
}

func RegisterRoutes() http.Handler {
	router := mux.NewRouter()

	authRouter := router.PathPrefix("/auth").Subrouter()
	registerAuthRoutes(authRouter)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"},
	})

	return c.Handler(router)
}
