package response

import (
	"log/slog"
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

const (
	SuccessStatus = "success"
	FailStatus    = "fail"
	ErrorStatus   = "error"
)

type ResponseMessage struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type ResponseData struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func HandleErrorStatus(c echo.Context, err error, method string) error {
	slog.Error("internal error", slog.String("method", method), slog.Any("err", err))
	return c.JSON(http.StatusInternalServerError, ResponseMessage{
		Status:  ErrorStatus,
		Message: "Internal server error.",
	})
}

func HandleFailStatus(c echo.Context, message string, code int) error {
	slog.Info("fail", slog.String("message", message), slog.Int("code", code))
	return c.JSON(code, ResponseMessage{Status: FailStatus, Message: message})
}

func HandleSuccessStatus(c echo.Context, object any, code int) error {
	switch reflect.TypeOf(object).Kind() {
	case reflect.String:
		return c.JSON(code, ResponseMessage{Status: SuccessStatus, Message: object.(string)})
	case reflect.Struct, reflect.Slice:
		return c.JSON(code, ResponseData{Status: SuccessStatus, Data: object})
	default:
		return c.JSON(code, map[string]string{"status": SuccessStatus})
	}
}
