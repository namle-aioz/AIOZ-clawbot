package dto

import "backend/models"

type SignInRequest struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type SignInResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *models.User `json:"user"`
}

type SignUpRequest struct {
	Email string `json:"email"    validate:"required,email"`
}

type ResendOTPRequest struct {
}

type VerifyOTPRequest struct {
	Password string `json:"password" validate:"required,min=6"`
	OTPCode  string `json:"otp_code" validate:"required,len=6"`
}

type VerifyOTPResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         *models.User `json:"user"`
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
