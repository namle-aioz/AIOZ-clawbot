package auth_actions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"gorm.io/gorm"

	"backend/models"
	"backend/utils/response"
	utils "backend/utils/token"
)

type GoogleSignInInput struct {
	IDToken string
}

type GoogleSignInAction struct {
	userRepo       models.UserRepository
	tokenIssuer    utils.TokenIssuer
	googleClientID string
}

type googleTokenInfo struct {
	Audience      string `json:"aud"`
	Subject       string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Name          string `json:"name"`
}

func NewGoogleSignInAction(userRepo models.UserRepository, tokenIssuer utils.TokenIssuer, googleClientID string) GoogleSignInAction {
	return GoogleSignInAction{userRepo: userRepo, tokenIssuer: tokenIssuer, googleClientID: googleClientID}
}

func (a GoogleSignInAction) Exec(ctx context.Context, input GoogleSignInInput) (*SignInOutput, error) {
	profile, err := a.verifyGoogleIDToken(ctx, input.IDToken)
	if err != nil {
		return nil, err
	}

	user, err := a.findOrCreateGoogleUser(ctx, profile)
	if err != nil {
		return nil, err
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

func (a GoogleSignInAction) verifyGoogleIDToken(ctx context.Context, idToken string) (*googleTokenInfo, error) {
	if a.googleClientID == "" {
		return nil, response.NewInternalError(fmt.Errorf("google client id is not configured"))
	}

	requestURL := "https://oauth2.googleapis.com/tokeninfo?id_token=" + url.QueryEscape(idToken)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, response.NewInternalError(err)
	}

	httpClient := &http.Client{Timeout: 10 * time.Second}
	httpResponse, err := httpClient.Do(request)
	if err != nil {
		return nil, response.NewHttpError(nil, "Invalid Google token.", http.StatusUnauthorized)
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, response.NewHttpError(nil, "Invalid Google token.", http.StatusUnauthorized)
	}

	var profile googleTokenInfo
	if err := json.NewDecoder(httpResponse.Body).Decode(&profile); err != nil {
		return nil, response.NewInternalError(err)
	}

	if profile.Audience != a.googleClientID {
		return nil, response.NewHttpError(nil, "Invalid Google token audience.", http.StatusUnauthorized)
	}

	if !strings.EqualFold(profile.EmailVerified, "true") {
		return nil, response.NewHttpError(nil, "Google email is not verified.", http.StatusUnauthorized)
	}

	if profile.Subject == "" || profile.Email == "" {
		return nil, response.NewHttpError(nil, "Invalid Google profile.", http.StatusUnauthorized)
	}

	return &profile, nil
}

func (a GoogleSignInAction) findOrCreateGoogleUser(ctx context.Context, profile *googleTokenInfo) (*models.User, error) {
	user, err := a.userRepo.GetUserByGoogleSubject(ctx, profile.Subject)
	if err == nil {
		return a.syncGoogleUser(ctx, user, profile)
	}
	if err != gorm.ErrRecordNotFound {
		return nil, response.NewInternalError(err)
	}

	user, err = a.userRepo.GetUserByEmail(ctx, profile.Email)
	if err == nil {
		return a.syncGoogleUser(ctx, user, profile)
	}
	if err != gorm.ErrRecordNotFound {
		return nil, response.NewInternalError(err)
	}

	googleSubject := profile.Subject
	avatarURL := emptyStringAsNil(profile.Picture)
	user = &models.User{
		FirstName:     profile.GivenName,
		LastName:      profile.FamilyName,
		Email:         strings.ToLower(profile.Email),
		DisplayEmail:  profile.Email,
		Password:      "",
		IsVerified:    true,
		AuthProvider:  "google",
		GoogleSubject: &googleSubject,
		AvatarURL:     avatarURL,
	}

	if err := a.userRepo.CreateUser(ctx, user); err != nil {
		return nil, response.NewInternalError(err)
	}

	return user, nil
}

func (a GoogleSignInAction) syncGoogleUser(ctx context.Context, user *models.User, profile *googleTokenInfo) (*models.User, error) {
	googleSubject := profile.Subject
	avatarURL := emptyStringAsNil(profile.Picture)
	user.FirstName = profile.GivenName
	user.LastName = profile.FamilyName
	user.Email = strings.ToLower(profile.Email)
	user.DisplayEmail = profile.Email
	user.IsVerified = true
	user.AuthProvider = "google"
	user.GoogleSubject = &googleSubject
	user.AvatarURL = avatarURL

	if err := a.userRepo.UpdateUser(ctx, user); err != nil {
		return nil, response.NewInternalError(err)
	}

	return user, nil
}

func emptyStringAsNil(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}

	return &trimmed
}
