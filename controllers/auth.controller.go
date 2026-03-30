package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	auth_actions "backend/actions/auth_actions"
	"backend/controllers/dto"
	"backend/utils/response"
	"backend/utils/session"
)

type AuthController struct {
	signInAction            auth_actions.SignInAction
	signUpAction            auth_actions.SignUpAction
	metaMaskChallengeAction auth_actions.MetaMaskChallengeAction
	metaMaskSignInAction    auth_actions.MetaMaskSignInAction
	signupSessionStore      *session.Store
}

func NewAuthController(
	signInAction auth_actions.SignInAction,
	signUpAction auth_actions.SignUpAction,
	metaMaskChallengeAction auth_actions.MetaMaskChallengeAction,
	metaMaskSignInAction auth_actions.MetaMaskSignInAction,
	signupSessionStore *session.Store,
) AuthController {
	return AuthController{
		signInAction:            signInAction,
		signUpAction:            signUpAction,
		metaMaskChallengeAction: metaMaskChallengeAction,
		metaMaskSignInAction:    metaMaskSignInAction,
		signupSessionStore:      signupSessionStore,
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
		return response.HandleError(ctx, response.NewHttpError(err, response.ErrInvalidRequestBody, http.StatusBadRequest))
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
		return response.HandleError(ctx, response.NewHttpError(err, response.ErrInvalidRequestBody, http.StatusBadRequest))
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
		return response.HandleError(ctx, response.NewHttpError(err, response.ErrInvalidRequestBody, http.StatusBadRequest))
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

// SignUp godoc
//
//	@Summary	Sign up a new user
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		body	body	dto.SignUpRequest	true	"User sign-up details"
//	@Success	200		{object}	response.ResponseMessage
//	@Failure	400		{object}	response.ResponseMessage
//	@Failure	401		{object}	response.ResponseMessage
//	@Router		/auth/signup [post]
func (c AuthController) SignUp(ctx echo.Context) error {
	var req dto.SignUpRequest
	if err := ctx.Bind(&req); err != nil {
		return response.HandleError(ctx, response.NewHttpError(err, response.ErrInvalidRequestBody, http.StatusBadRequest))
	}

	if err := ctx.Validate(&req); err != nil {
		return err
	}

	err := c.signUpAction.SignUp(ctx.Request().Context(), auth_actions.SignUpInput{
		Email: req.Email,
	})
	if err != nil {
		return response.HandleError(ctx, err)
	}

	if err := c.signupSessionStore.Set(ctx, session.EmailKey, req.Email); err != nil {
		return response.HandleError(ctx, err)
	}

	return response.HandleSuccessStatus(ctx, "Send OTP code successfully.", http.StatusOK)
}

// ResendOTPcode godoc
//
//	@Summary	Resend OTP code for user sign-up
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		body	body		dto.ResendOTPRequest	true	"User sign-up details"
//	@Success	200		{object}	response.ResponseMessage
//	@Failure	400		{object}	response.ResponseMessage
//	@Failure	401		{object}	response.ResponseMessage
//	@Router		/auth/resend-otp [post]
func (c AuthController) ResendOTPcode(ctx echo.Context) error {
	var req dto.ResendOTPRequest
	if err := ctx.Bind(&req); err != nil {
		return response.HandleError(ctx, response.NewHttpError(err, response.ErrInvalidRequestBody, http.StatusBadRequest))
	}

	if err := ctx.Validate(&req); err != nil {
		return err
	}

	email, err := c.signupSessionStore.Get(ctx, session.EmailKey)
	if err != nil {
		return response.HandleError(ctx, err)
	}

	err = c.signUpAction.ResendOTPcode(ctx.Request().Context(), auth_actions.SignUpInput{
		Email: email,
	})
	if err != nil {
		return response.HandleError(ctx, err)
	}

	if err := c.signupSessionStore.Set(ctx, session.EmailKey, email); err != nil {
		return response.HandleError(ctx, err)
	}

	return response.HandleSuccessStatus(ctx, "OTP code sent successfully.", http.StatusOK)
}

// VerifyOTP godoc
//
//	@Summary	Verify OTP for user sign-up
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		body	body		dto.VerifyOTPRequest	true	"User sign-up details"
//	@Success	200		{object}	dto.VerifyOTPResponse
//	@Failure	400		{object}	response.ResponseMessage
//	@Failure	401		{object}	response.ResponseMessage
//	@Router		/auth/verify-otp [post]
func (c AuthController) VerifyOTP(ctx echo.Context) error {
	var req dto.VerifyOTPRequest
	if err := ctx.Bind(&req); err != nil {
		return response.HandleError(ctx, response.NewHttpError(err, response.ErrInvalidRequestBody, http.StatusBadRequest))
	}

	if err := ctx.Validate(&req); err != nil {
		return err
	}

	email, err := c.signupSessionStore.Get(ctx, session.EmailKey)
	if err != nil {
		return response.HandleError(ctx, err)
	}

	result, err := c.signUpAction.VerifyOTP(ctx.Request().Context(), auth_actions.VerifyOTPInput{
		Email:    email,
		Password: req.Password,
		OTPCode:  req.OTPCode,
	})
	if err != nil {
		return response.HandleError(ctx, err)
	}

	c.signupSessionStore.Clear(ctx)

	return response.HandleSuccessStatus(ctx, dto.SignInResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		User:         result.User,
	}, http.StatusOK)
}
