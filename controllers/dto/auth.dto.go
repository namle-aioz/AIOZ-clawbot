package dto

import "backend/models"

type SignInRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type GoogleSignInRequest struct {
	IDToken string `json:"id_token" validate:"required"`
}

type MetaMaskChallengeRequest struct {
	WalletAddress string `json:"wallet_address" validate:"required"`
}

type MetaMaskChallengeResponse struct {
	Message   string `json:"message"`
	ExpiresAt string `json:"expires_at"`
}

type MetaMaskSignInRequest struct {
	WalletAddress string `json:"wallet_address" validate:"required"`
	Message       string `json:"message" validate:"required"`
	Signature     string `json:"signature" validate:"required"`
}

type SignInResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *models.User `json:"user"`
}
