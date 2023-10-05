package get_object

import (
	"errors"
	"io"
	"net/http"
	"os"

	"storage-gateway/application/api/apierror"
	"storage-gateway/domain/models"
	"storage-gateway/domain/services"

	"github.com/labstack/echo/v4"
)

type GetObjectHandler struct {
	getObjectService *services.GetObjectService
}

func NewGetObjectHandler(getObjectService *services.GetObjectService) *GetObjectHandler {
	return &GetObjectHandler{
		getObjectService: getObjectService,
	}
}

func (h *GetObjectHandler) GetObject(c echo.Context) error {
	id := c.Param("objectID")

	obj, err := h.getObjectService.GetObject(c.Request().Context(), models.ObjectID(id))
	if err != nil {
		switch {
		case errors.Is(err, models.ErrObjectIDNotValid):
			return apierror.Err(c, http.StatusBadRequest, err)
		case errors.Is(err, models.ErrObjectNotFound):
			return apierror.Err(c, http.StatusNotFound, err)
		case errors.Is(err, models.ErrObjectStorageNotAvailable):
			return apierror.Err(c, http.StatusServiceUnavailable, err)
		case os.IsTimeout(err):
			return apierror.Err(c, http.StatusBadGateway, err)
		default:
			return apierror.Err(c, http.StatusInternalServerError, err)
		}
	}

	_, err = io.Copy(c.Response().Writer, obj.Content)
	if err != nil {
		return apierror.Err(c, http.StatusInternalServerError, err)
	}

	c.Response().Header().Set(echo.HeaderContentType, obj.ContentType)
	c.Response().WriteHeader(http.StatusOK)

	return nil
}
