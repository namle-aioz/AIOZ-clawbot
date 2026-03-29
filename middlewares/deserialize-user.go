package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"backend/models"
	"backend/utils/response"
	"backend/utils/token"
)

func DeserializeUser(
	tokenIssuer token.TokenIssuer,
	userRepo models.UserRepository,
) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userId, err := getUserIdFromToken(c, tokenIssuer)
			if err != nil {
				return response.HandleError(c, err)
			}

			user, err := userRepo.GetUserById(c.Request().Context(), userId)
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					return response.HandleFailStatus(c, "Record not found.", http.StatusUnauthorized)
				}
				return response.HandleErrorStatus(c, err, "DeserializeUser")
			}

			if !user.IsVerified {
				return response.HandleFailStatus(c, "Your account is not verified.", http.StatusForbidden)
			}

			sessionId, err := uuid.Parse(c.Request().Header.Get("Sync-Session-Id"))
			if err != nil {
				sessionId = uuid.New()
			}

			c.Set("currentUser", *user)
			c.Set("sessionId", sessionId)
			c.Set("authInfo", *models.NewAuthenticationInfo(user, sessionId))
			return next(c)
		}
	}
}

func getUserIdFromToken(c echo.Context, t token.TokenIssuer) (uuid.UUID, error) {
	authHeader := c.Request().Header.Get("Authorization")
	var accessToken string
	if fields := strings.Fields(authHeader); len(fields) == 2 && fields[0] == "Bearer" {
		accessToken = fields[1]
	}

	if accessToken == "" {
		return uuid.Nil, response.NewHttpErrorWithNoMsg(
			fmt.Errorf("you are not logged in"),
			http.StatusUnauthorized,
		)
	}

	sub, err := t.ValidateAccessToken(accessToken)
	if err != nil {
		return uuid.Nil, response.NewHttpErrorWithNoMsg(response.InvalidToken, http.StatusUnauthorized)
	}

	userId, err := uuid.Parse(fmt.Sprint(sub["user_id"]))
	if err != nil {
		return uuid.Nil, response.NewHttpErrorWithNoMsg(response.InvalidToken, http.StatusUnauthorized)
	}

	return userId, nil
}
