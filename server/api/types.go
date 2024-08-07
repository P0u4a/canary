package api

type User struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

type AuthResponse struct {
	AccessToken  string `json:"accesstoken"`
	RefreshToken string `json:"refreshtoken"`
}

type CanarySignInResponse struct {
	Verified bool `json:"verified"`
}
