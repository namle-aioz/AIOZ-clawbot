package routes

import (
	"github.com/labstack/echo/v4"

	"backend/controllers"
)

type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) AuthRouteController {
	return AuthRouteController{authController: authController}
}

func (rc *AuthRouteController) RegisterRoute(g *echo.Group) {
	auth := g.Group("/auth")
	auth.POST("/signin", rc.authController.SignIn)
	auth.POST("/signup", rc.authController.SignUp)
	auth.POST("/verify-otp", rc.authController.VerifyOTP)
	auth.POST("/resend-otp", rc.authController.ResendOTPcode)
	auth.POST("/metamask/challenge", rc.authController.MetaMaskChallenge)
	auth.POST("/metamask", rc.authController.MetaMaskSignIn)
}
