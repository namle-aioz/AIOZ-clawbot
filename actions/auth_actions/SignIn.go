package auth_actions

import (
	"context"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"backend/models"
	"backend/utils/response"
	utils "backend/utils/token"
)

type SignInInput struct {
	Email    string
	Password string
}

type SignInOutput struct {
	AccessToken  string
	RefreshToken string
	User         *models.User
}

type SignInAction struct {
	userRepo    models.UserRepository
	tokenIssuer utils.TokenIssuer
}

func NewSignInAction(userRepo models.UserRepository, tokenIssuer utils.TokenIssuer) SignInAction {
	return SignInAction{userRepo: userRepo, tokenIssuer: tokenIssuer}
}

func (a SignInAction) Exec(ctx context.Context, input SignInInput) (*SignInOutput, error) {
	user, err := a.userRepo.GetUserByEmail(ctx, input.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, response.NewHttpError(nil, response.ErrInvalidUserAndPassword, http.StatusUnauthorized)
		}
		return nil, response.NewInternalError(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, response.NewHttpError(nil, response.ErrInvalidUserAndPassword, http.StatusUnauthorized)
	}

	accessToken, refreshToken, _, err := a.tokenIssuer.CreateCredential(user.Id)
	if err != nil {
		return nil, response.NewInternalError(err)
	}

	return &SignInOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}
