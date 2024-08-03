package api

type AuthRequest struct {
	Username  string `json:"username"`
	VoiceData []byte `json:"voicedata"`
}

type AuthResponse struct {
	AccessToken  string `json:"accesstoken"`
	RefreshToken string `json:"refreshtoken"`
}

type CanarySignUpResponse struct {
	Features   []float64 `json:"features"`
	Passphrase string    `json:"passphrase"`
}

type CanarySignInResponse struct {
	SimilarityScore float64 `json:"similarityScore"`
	TranscribedPass string  `json:"transcribedPass"`
}
