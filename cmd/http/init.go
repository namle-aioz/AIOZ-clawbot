package main

import (
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"

	auth_actions "backend/actions/auth_actions"
	user_actions "backend/actions/user_actions"
	"backend/controllers"
	"backend/db"
	"backend/initializers"
	"backend/middlewares"
	"backend/models"
	routes "backend/routers"
	custom_log "backend/utils/log"
	"backend/utils/token"
)

var (
	config *initializers.Config

	router *echo.Echo

	tokenIssuer token.TokenIssuer

	userRepo models.UserRepository

	authMiddleware echo.MiddlewareFunc

	signInAction            auth_actions.SignInAction
	googleSignInAction      auth_actions.GoogleSignInAction
	metaMaskChallengeAction auth_actions.MetaMaskChallengeAction
	metaMaskSignInAction    auth_actions.MetaMaskSignInAction
	getMeAction             user_actions.GetMeAction

	authController controllers.AuthController
	userController controllers.UserController

	authRouteController routes.AuthRouteController
	userRouteController routes.UserRouteController
)

func init() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "app"
	}

	config = initializers.MustLoadConfig(".", env)
	initializers.ConnectDB(config)

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	router = echo.New()
	router.HideBanner = true

	slog.SetDefault(slog.New(
		custom_log.NewHandler(
			&slog.HandlerOptions{},
		),
	))

	tokenIssuer = token.MustNewTokenIssuer(
		config.AccessTokenPrivateKey,
		config.AccessTokenPublicKey,
		config.RefreshTokenPrivateKey,
		config.RefreshTokenPublicKey,
		int(config.AccessTokenExpiresIn.Seconds()),
		int(config.RefreshTokenExpiresIn.Seconds()),
	)

	userRepo = db.MustNewUserRepository(initializers.DB, true)

	authMiddleware = middlewares.DeserializeUser(tokenIssuer, userRepo)

	signInAction = auth_actions.NewSignInAction(userRepo, tokenIssuer)
	googleSignInAction = auth_actions.NewGoogleSignInAction(userRepo, tokenIssuer, config.GoogleClientID)
	metaMaskChallengeAction = auth_actions.NewMetaMaskChallengeAction()
	metaMaskSignInAction = auth_actions.NewMetaMaskSignInAction(userRepo, tokenIssuer)
	getMeAction = user_actions.NewGetMeAction()

	authController = controllers.NewAuthController(signInAction, googleSignInAction, metaMaskChallengeAction, metaMaskSignInAction)
	userController = controllers.NewUserController(getMeAction)

	authRouteController = routes.NewAuthRouteController(authController)
	userRouteController = routes.NewUserRouteController(userController, authMiddleware)
}
