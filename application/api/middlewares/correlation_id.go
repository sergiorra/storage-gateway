package middlewares

import (
	"storage-gateway/internal/context-wrapper"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// CorrelationID returns an Echo middleware that generates a correlation ID using UUID and adds it to the request context
func CorrelationID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			correlationID := uuid.New().String()
			c.Request().Header.Set(echo.HeaderXRequestID, correlationID)
			c.SetRequest(c.Request().WithContext(context_wrapper.WithCorrelationID(c.Request().Context(), correlationID)))

			return next(c)
		}
	}
}
