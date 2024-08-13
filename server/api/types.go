package api

type User struct {
	Username string `json:"username"`
}

type AuthResponse struct {
	AccessToken  string `json:"accesstoken"`
	RefreshToken string `json:"refreshtoken"`
}

type CanarySignInResponse struct {
	Verified int32 `json:"verified"`
}
