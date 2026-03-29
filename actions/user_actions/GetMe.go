package user_actions

import (
	"github.com/labstack/echo/v4"

	"backend/models"
)

type GetMeAction struct{}

func NewGetMeAction() GetMeAction {
	return GetMeAction{}
}

func (a GetMeAction) Exec(c echo.Context) (*models.User, error) {
	authInfo := c.Get("authInfo").(models.AuthenticationInfo)
	return authInfo.User, nil
}
