package response

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

var (
	InvalidToken = fmt.Errorf("invalid token")
)

type HttpError struct {
	cause   error
	Message string
	Code    int
}

func (e HttpError) Error() string {
	if e.cause != nil {
		return e.cause.Error()
	}
	return e.Message
}

func NewHttpError(cause error, message string, code int) HttpError {
	return HttpError{cause: cause, Message: message, Code: code}
}

func NewHttpErrorWithNoMsg(cause error, code int) HttpError {
	return HttpError{cause: cause, Message: cause.Error(), Code: code}
}

func NewInternalError(cause error) HttpError {
	return NewHttpError(cause, "Internal server error.", http.StatusInternalServerError)
}

func NewNotFoundError(message string) HttpError {
	return NewHttpError(nil, message, http.StatusNotFound)
}

func HandleError(c echo.Context, err error) error {
	httpErr, ok := err.(HttpError)
	if !ok {
		return HandleErrorStatus(c, err, "unhandled error")
	}

	if httpErr.cause != nil {
		slog.Info("http error",
			slog.String("err", httpErr.cause.Error()),
			slog.String("msg", httpErr.Message),
			slog.Int("code", httpErr.Code),
		)
	}

	status := FailStatus
	if httpErr.Code == http.StatusInternalServerError {
		status = ErrorStatus
	}

	return c.JSON(httpErr.Code, map[string]string{"status": status, "message": httpErr.Message})
}
