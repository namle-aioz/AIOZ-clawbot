package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	user_actions "backend/actions/user_actions"
	"backend/utils/response"
)

type UserController struct {
	getMeAction user_actions.GetMeAction
}

func NewUserController(getMeAction user_actions.GetMeAction) UserController {
	return UserController{getMeAction: getMeAction}
}

// GetMe godoc
//
//	@Summary	Get current user
//	@Tags		user
//	@Security	jwt
//	@Produce	json
//	@Param		Authorization	header		string	false	"Bearer Token"
//	@Success	200				{object}	models.User
//	@Failure	401				{object}	response.ResponseMessage
//	@Failure	403				{object}	response.ResponseMessage
//	@Router		/user/me [get]
func (c UserController) GetMe(ctx echo.Context) error {
	user, err := c.getMeAction.Exec(ctx)
	if err != nil {
		return response.HandleError(ctx, err)
	}

	return response.HandleSuccessStatus(ctx, *user, http.StatusOK)
}
