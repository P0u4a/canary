package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	FEATURES_ENDPOINT   = "localhost:3001/process-voice"
	TRANSCRIBE_ENDPOINT = "localhost:3001/analyse-voice"
)

func getFeaturesAndPassphrase(voiceData []byte) ([]float64, string, error) {
	payload := map[string][]byte{"data": voiceData}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, "", err
	}

	res, err := http.Post(FEATURES_ENDPOINT, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var featuresAndPassphrase CanarySignUpResponse
	if err := json.NewDecoder(res.Body).Decode(&featuresAndPassphrase); err != nil {
		return nil, "", err
	}

	return featuresAndPassphrase.Features, featuresAndPassphrase.Passphrase, nil

}

func getTranscribedPassAndSimScore(userData User) (string, float64, error) {
	payload := map[string]User{"data": userData}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", 0.0, err
	}

	res, err := http.Post(TRANSCRIBE_ENDPOINT, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", 0.0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", 0.0, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var transcribedPassAndSimScore CanarySignInResponse
	if err := json.NewDecoder(res.Body).Decode(&transcribedPassAndSimScore); err != nil {
		return "", 0.0, err
	}

	return transcribedPassAndSimScore.TranscribedPass, transcribedPassAndSimScore.SimilarityScore, nil
}

func createTokenPair(username string) (string, string, error) {
	accessTokenClaims := jwt.MapClaims{
		"Audience":  username,
		"ExpiresAt": time.Now().Add(time.Hour * 1).Unix(),
	}

	refreshTokenClaims := jwt.MapClaims{
		"Audience":  username,
		"ExpiresAt": time.Now().Add(time.Hour * 24).Unix(),
	}

	secret := os.Getenv("JWT_SECRET")
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)

	accessString, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	refreshString, err := refreshToken.SignedString([]byte(secret))
	if err != nil {
		return "", "", err
	}

	return accessString, refreshString, nil
}

func regenerateJWT(username string) (string, error) {
	// verify user exists first

	accessTokenClaims := jwt.MapClaims{
		"Audience":  username,
		"ExpiresAt": time.Now().Add(time.Hour * 24).Unix(),
	}

	secret := os.Getenv("JWT_SECRET")
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	accessString, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return accessString, nil
}

func validateJWT(token string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})
}
