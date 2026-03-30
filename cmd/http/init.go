package main

import (
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	auth_actions "backend/actions/auth_actions"
	"backend/actions/mail_actions"
	user_actions "backend/actions/user_actions"
	"backend/controllers"
	"backend/db"
	"backend/initializers"
	"backend/middlewares"
	"backend/models"
	routes "backend/routers"
	custom_log "backend/utils/log"
	"backend/utils/session"
	"backend/utils/token"

	"github.com/resend/resend-go/v3"
)

var (
	config *initializers.Config

	router *echo.Echo

	tokenIssuer token.TokenIssuer

	userRepo models.UserRepository
	otpRepo  models.OTPCodeRepository

	authMiddleware echo.MiddlewareFunc

	signInAction            auth_actions.SignInAction
	signUpAction            auth_actions.SignUpAction
	sendMailAction          mail_actions.MailSender
	metaMaskChallengeAction auth_actions.MetaMaskChallengeAction
	metaMaskSignInAction    auth_actions.MetaMaskSignInAction
	getMeAction             user_actions.GetMeAction
	sessionStore            *session.Store

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
	router.Validator = &DefaultValidator{validator: validator.New()}

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

	sendMailClient := resend.NewClient(config.ResendMailAPIKey)

	userRepo = db.MustNewUserRepository(initializers.DB, true)
	otpRepo = db.MustNewOTPCodeRepository(initializers.DB, true)

	authMiddleware = middlewares.DeserializeUser(tokenIssuer, userRepo)

	sendMailAction = mail_actions.NewResendMailSenderAction(sendMailClient, config.EmailFrom)
	signInAction = auth_actions.NewSignInAction(userRepo, tokenIssuer)
	signUpAction = auth_actions.NewSignUpAction(userRepo, otpRepo, sendMailAction, tokenIssuer, config.BcryptCost)
	sessionStore = session.NewStore(auth_actions.OTPExpireTTL, config.SessionSecret)
	metaMaskChallengeAction = auth_actions.NewMetaMaskChallengeAction()
	metaMaskSignInAction = auth_actions.NewMetaMaskSignInAction(userRepo, tokenIssuer)
	getMeAction = user_actions.NewGetMeAction()

	authController = controllers.NewAuthController(signInAction, signUpAction, metaMaskChallengeAction, metaMaskSignInAction, sessionStore)
	userController = controllers.NewUserController(getMeAction)

	authRouteController = routes.NewAuthRouteController(authController)
	userRouteController = routes.NewUserRouteController(userController, authMiddleware)
}

type DefaultValidator struct {
	validator *validator.Validate
}

func (v *DefaultValidator) Validate(i any) error {
	if v.validator == nil {
		v.validator = validator.New()
	}

	return v.validator.Struct(i)
}
