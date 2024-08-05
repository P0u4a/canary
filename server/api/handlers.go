package api

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
)

func HandleSignUp(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		var req AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		features, passphrase, err := getFeaturesAndPassphrase(req.VoiceData)
		if err != nil {
			http.Error(w, "Error handling voice data", http.StatusInternalServerError)
			return
		}

		db.Set(req.Username, User{features, passphrase})
	}

}

func HandleSignIn(db *DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
			return
		}

		var req AuthRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}

		userData, ok := db.Get(req.Username)
		if !ok {
			http.Error(w, "User does not exist", http.StatusUnauthorized)
			return
		}

		transcribedPass, similarityScore, err := getTranscribedPassAndSimScore(userData.VoiceData, req.VoiceData)
		if err != nil {
			http.Error(w, "Error retrieving voice data", http.StatusInternalServerError)
			return
		}

		similarityThreshold, err := strconv.ParseFloat(os.Getenv("SIMILARITY_THRESHOLD"), 64)
		if err != nil {
			http.Error(w, "Interal Server Error", http.StatusInternalServerError)
			return
		}

		if similarityScore < similarityThreshold || transcribedPass != userData.Passphrase {
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
