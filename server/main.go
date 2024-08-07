package main

import (
	"log"
	"net/http"

	"github.com/P0u4a/canary/api"
	"github.com/rs/cors"
)

func main() {
	database := api.NewDatabase()
	mux := http.NewServeMux()
	mux.HandleFunc("/signin", api.HandleSignIn(database))
	mux.HandleFunc("/signup", api.HandleSignUp(database))
	mux.HandleFunc("/protected", api.HandleProtected(database))

	middleware := cors.Default().Handler(mux)

	if err := http.ListenAndServe(":3000", middleware); err != nil {
		log.Fatal("Server error:", err)
	}
}
