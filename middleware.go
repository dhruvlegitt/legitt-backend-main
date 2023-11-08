package main

import (
	"log"
	"net/http"
	"time"
)

type MiddleWare func(http.Handler) http.Handler

// Wraps the http handlers using function closures
func WrapMiddlewares(h http.Handler, middlewares ...MiddleWare) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func RequestTimeMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		h.ServeHTTP(w, r)
		log.Printf("time take : %s", time.Since(start))
	})
}

func SetHeadersMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	})
}
