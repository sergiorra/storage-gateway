package put_object

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"os"

	"storage-gateway/application/api/apierror"
	"storage-gateway/domain/models"
	"storage-gateway/domain/services"

	"github.com/labstack/echo/v4"
)

type PutObjectHandler struct {
	putObjectService *services.PutObjectService
}

func NewPutObjectHandler(putObjectService *services.PutObjectService) *PutObjectHandler {
	return &PutObjectHandler{
		putObjectService: putObjectService,
	}
}

func (h *PutObjectHandler) PutObject(c echo.Context) error {
	id := c.Param("objectID")

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, c.Request().Body); err != nil || c.Request().Body == nil {
		return apierror.Err(c, http.StatusBadRequest, err)
	}

	obj := &models.Object{
		ID:          models.ObjectID(id),
		Content:     &buf,
		ContentType: c.Request().Header.Get("Content-Type"),
		Size:        c.Request().ContentLength,
	}

	err := h.putObjectService.PutObject(c.Request().Context(), obj)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrObjectIDNotValid):
			return apierror.Err(c, http.StatusBadRequest, err)
		case errors.Is(err, models.ErrObjectStorageNotAvailable):
			return apierror.Err(c, http.StatusServiceUnavailable, err)
		case os.IsTimeout(err):
			return apierror.Err(c, http.StatusBadGateway, err)
		default:
			return apierror.Err(c, http.StatusInternalServerError, err)
		}
	}

	return c.JSON(http.StatusOK, "")
}
