package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	SIGN_UP_ENDPOINT = "http://localhost:3001/process-voice"
	SIGN_IN_ENDPOINT = "http://localhost:3001/analyse-voice"
)

func sendMultipartReq(username string, voiceData multipart.File, endpoint string) (*http.Response, error) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	audioFile, err := writer.CreateFormFile("voicedata", "audio.wav")
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(audioFile, voiceData)
	if err != nil {
		return nil, err
	}

	err = writer.WriteField("username", username)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", endpoint, &buffer)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil

}

func initModel(username string, voiceData multipart.File) error {

	res, err := sendMultipartReq(username, voiceData, SIGN_UP_ENDPOINT)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

	return nil

}

func verifyVoice(username string, voiceData multipart.File) (bool, error) {

	res, err := sendMultipartReq(username, voiceData, SIGN_IN_ENDPOINT)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return false, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	var status CanarySignInResponse
	if err := json.NewDecoder(res.Body).Decode(&status); err != nil {
		return false, err
	}

	return status.Verified, nil
}

func createTokenPair() (string, string, error) {
	accessTokenClaims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	}

	refreshTokenClaims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 24).Unix(),
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

func regenerateJWT() (string, error) {
	accessTokenClaims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	}

	secret := os.Getenv("JWT_SECRET")
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)

	accessString, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return accessString, nil
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if err := validateClaims(claims); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("invalid token claims")
	}

	return token, nil

}

func validateClaims(claims jwt.MapClaims) error {
	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return fmt.Errorf("token has expired")
		}
	}

	return nil
}
