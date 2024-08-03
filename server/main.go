package main

import (
	"log"
	"net/http"

	"github.com/P0u4a/canary/api"
)

func main() {
	database := api.NewDatabase()

	http.HandleFunc("/signin", api.HandleSignIn(database))
	http.HandleFunc("/signup", api.HandleSignUp(database))
	http.HandleFunc("/protected", api.HandleProtected(database))

	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
