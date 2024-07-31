package main

import (
	"log"
	"net/http"
)

/**
Routes:
	/signin
		1. Fetch encrypted features and passphrase for given username
		2. Send POST request with encrypted features and passphrase to ML server
		3. Compare response with hash and verify similarity score
		4. If passes, generate and send JWT, else return error

	/signup
		1. Send POST request with audio data for new hashed passphrase and encrypted features
		2. Save response in db for username
		3. Send 2xx and redirect if successfull or error

	/protected (Example route for testing protected routes)
		1. Validate the JWT token in the request
		2. Either serve the route or error

Use hashicorp's memdb in-memory DB to store username, encrypted features, hashed passphrase
Use golang-jwt to generate the JWT token

*/

func signUp(w http.ResponseWriter, r *http.Request) {

}

func signIn(w http.ResponseWriter, r *http.Request) {

}

func handleProtected(w http.ResponseWriter, r *http.Request) {

}

func main() {
	http.HandleFunc("/signin", signIn)
	http.HandleFunc("/signup", signUp)
	http.HandleFunc("/protected", handleProtected)

	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal("Server error:", err)
	}
}
