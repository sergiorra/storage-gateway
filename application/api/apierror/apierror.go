package apierror

import "github.com/labstack/echo/v4"

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func Err(c echo.Context, httpCode int, reqErr error) error {
	if err := c.JSON(httpCode, ErrorResponse{Status: "error", Message: reqErr.Error()}); err != nil {
		reqErr = err
	}

	return reqErr
}
