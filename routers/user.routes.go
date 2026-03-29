package routes

import (
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"

	"backend/controllers"
	"backend/middlewares"
)

type UserRouteController struct {
	userController controllers.UserController
	authMiddleware echo.MiddlewareFunc
}

func NewUserRouteController(
	userController controllers.UserController,
	authMiddleware echo.MiddlewareFunc,
) UserRouteController {
	return UserRouteController{
		userController: userController,
		authMiddleware: authMiddleware,
	}
}

func (rc *UserRouteController) RegisterRoute(g *echo.Group) {
	user := g.Group("/user",
		middlewares.NewRateLimiter("/user", rate.NewLimiter(rate.Every(time.Second), 5), time.Hour),
		rc.authMiddleware,
	)
	user.GET("/me", rc.userController.GetMe)
}
