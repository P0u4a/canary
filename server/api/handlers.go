package api

import (
	"encoding/json"
	"net/http"
)

func HandleSignUp(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		// 10 MB limit on file size
		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		audioFile, _, err := r.FormFile("voicedata")
		if err != nil {
			http.Error(w, "Error retrieving file", http.StatusBadRequest)
			return
		}
		defer audioFile.Close()

		username := r.FormValue("username")

		err = initModel(username, audioFile)
		if err != nil {
			http.Error(w, "Error handling voice data", http.StatusInternalServerError)
			return
		}

		db.Set(username, User{username, "member"})

		accessToken, refreshToken, err := createTokenPair()
		if err != nil {
			http.Error(w, "Failed to generate token pair", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		res := AuthResponse{AccessToken: accessToken, RefreshToken: refreshToken}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			return
		}
	}

}

func HandleSignIn(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		err := r.ParseMultipartForm(10 << 20)
		if err != nil {
			http.Error(w, "Unable to parse form", http.StatusBadRequest)
			return
		}

		audioFile, _, err := r.FormFile("voicedata")
		if err != nil {
			http.Error(w, "Error retrieving file", http.StatusBadRequest)
			return
		}
		defer audioFile.Close()

		// Retrieve the username from the form data
		username := r.FormValue("username")

		_, ok := db.Get(username)
		if !ok {
			http.Error(w, "User does not exist", http.StatusUnauthorized)
			return
		}

		verified, err := verifyVoice(username, audioFile)
		if err != nil {
			http.Error(w, "Error verifying voice data", http.StatusInternalServerError)
			return
		}

		if !verified {
			http.Error(w, "Sign in failed", http.StatusUnauthorized)
			return
		}

		accessToken, refreshToken, err := createTokenPair()
		if err != nil {
			http.Error(w, "Failed to generate token pair", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		res := AuthResponse{AccessToken: accessToken, RefreshToken: refreshToken}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
			return
		}
	}
}

func HandleRefresh(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("x-refresh-token")

		token, err := validateJWT(authHeader)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		if !token.Valid {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		var username map[string]string
		if err := json.NewDecoder(r.Body).Decode(&username); err != nil {
			http.Error(w, "Failed to decode JSON", http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		accessToken, err := regenerateJWT()
		if err != nil {
			http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		}
		if err := json.NewEncoder(w).Encode(accessToken); err != nil {
			http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		}
	}
}

func HandleProtected(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		authHeader := r.Header.Get("x-access-token")

		_, err := validateJWT(authHeader)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
