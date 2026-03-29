package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	auth_actions "backend/actions/auth_actions"
	"backend/controllers/dto"
	"backend/utils/response"
)

type AuthController struct {
	signInAction            auth_actions.SignInAction
	googleSignInAction      auth_actions.GoogleSignInAction
	metaMaskChallengeAction auth_actions.MetaMaskChallengeAction
	metaMaskSignInAction    auth_actions.MetaMaskSignInAction
}

func NewAuthController(
	signInAction auth_actions.SignInAction,
	googleSignInAction auth_actions.GoogleSignInAction,
	metaMaskChallengeAction auth_actions.MetaMaskChallengeAction,
	metaMaskSignInAction auth_actions.MetaMaskSignInAction,
) AuthController {
	return AuthController{
		signInAction:            signInAction,
		googleSignInAction:      googleSignInAction,
		metaMaskChallengeAction: metaMaskChallengeAction,
		metaMaskSignInAction:    metaMaskSignInAction,
	}
}

// SignIn godoc
//
//	@Summary	Sign in
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		body	body		dto.SignInRequest	true	"Sign in credentials"
//	@Success	200		{object}	dto.SignInResponse
//	@Failure	400		{object}	response.ResponseMessage
//	@Failure	401		{object}	response.ResponseMessage
//	@Router		/auth/signin [post]
func (c AuthController) SignIn(ctx echo.Context) error {
	var req dto.SignInRequest
	if err := ctx.Bind(&req); err != nil {
		return response.HandleError(ctx, response.NewHttpError(err, "Invalid request body.", http.StatusBadRequest))
	}

	if err := ctx.Validate(&req); err != nil {
		return err
	}

	result, err := c.signInAction.Exec(ctx.Request().Context(), auth_actions.SignInInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return response.HandleError(ctx, err)
	}

	return response.HandleSuccessStatus(ctx, dto.SignInResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		User:         result.User,
	}, http.StatusOK)
}

// GoogleSignIn godoc
//
//	@Summary	Sign in with Google
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		body	body	dto.GoogleSignInRequest	true	"Google ID token"
//	@Success	200		{object}	dto.SignInResponse
//	@Failure	400		{object}	response.ResponseMessage
//	@Failure	401		{object}	response.ResponseMessage
//	@Router		/auth/google [post]
func (c AuthController) GoogleSignIn(ctx echo.Context) error {
	var req dto.GoogleSignInRequest
	if err := ctx.Bind(&req); err != nil {
		return response.HandleError(ctx, response.NewHttpError(err, "Invalid request body.", http.StatusBadRequest))
	}

	if err := ctx.Validate(&req); err != nil {
		return err
	}

	result, err := c.googleSignInAction.Exec(ctx.Request().Context(), auth_actions.GoogleSignInInput{IDToken: req.IDToken})
	if err != nil {
		return response.HandleError(ctx, err)
	}

	return response.HandleSuccessStatus(ctx, dto.SignInResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		User:         result.User,
	}, http.StatusOK)
}

// MetaMaskChallenge godoc
//
//	@Summary	Create a MetaMask sign-in challenge
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		body	body	dto.MetaMaskChallengeRequest	true	"Wallet address"
//	@Success	200		{object}	dto.MetaMaskChallengeResponse
//	@Failure	400		{object}	response.ResponseMessage
//	@Router		/auth/metamask/challenge [post]
func (c AuthController) MetaMaskChallenge(ctx echo.Context) error {
	var req dto.MetaMaskChallengeRequest
	if err := ctx.Bind(&req); err != nil {
		return response.HandleError(ctx, response.NewHttpError(err, "Invalid request body.", http.StatusBadRequest))
	}

	if err := ctx.Validate(&req); err != nil {
		return err
	}

	result, err := c.metaMaskChallengeAction.Exec(ctx.Request().Context(), auth_actions.MetaMaskChallengeInput{WalletAddress: req.WalletAddress})
	if err != nil {
		return response.HandleError(ctx, err)
	}

	return response.HandleSuccessStatus(ctx, dto.MetaMaskChallengeResponse{
		Message:   result.Message,
		ExpiresAt: result.ExpiresAt,
	}, http.StatusOK)
}

// MetaMaskSignIn godoc
//
//	@Summary	Sign in with MetaMask
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		body	body	dto.MetaMaskSignInRequest	true	"Wallet login payload"
//	@Success	200		{object}	dto.SignInResponse
//	@Failure	400		{object}	response.ResponseMessage
//	@Failure	401		{object}	response.ResponseMessage
//	@Router		/auth/metamask [post]
func (c AuthController) MetaMaskSignIn(ctx echo.Context) error {
	var req dto.MetaMaskSignInRequest
	if err := ctx.Bind(&req); err != nil {
		return response.HandleError(ctx, response.NewHttpError(err, "Invalid request body.", http.StatusBadRequest))
	}

	if err := ctx.Validate(&req); err != nil {
		return err
	}

	result, err := c.metaMaskSignInAction.Exec(ctx.Request().Context(), auth_actions.MetaMaskSignInInput{
		WalletAddress: req.WalletAddress,
		Message:       req.Message,
		Signature:     req.Signature,
	})
	if err != nil {
		return response.HandleError(ctx, err)
	}

	return response.HandleSuccessStatus(ctx, dto.SignInResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		User:         result.User,
	}, http.StatusOK)
}
