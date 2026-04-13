package httpservice

import (
	"cdek/internal/adapter/in/dto"
	"cdek/internal/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func mapErrorCode(err error) (int, dto.ErrorResponse) {
	switch {
	case errors.Is(err, model.ErrInvalidRequest):
		return http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "invalid request",
		}
	case errors.Is(err, model.ErrUnauthorized):
		return http.StatusUnauthorized, dto.ErrorResponse{
			Error:   "unauthorized",
			Message: "authentication required",
		}
	case errors.Is(err, model.ErrForbidden):
		return http.StatusForbidden, dto.ErrorResponse{
			Error:   "forbidden",
			Message: "access denied",
		}
	case errors.Is(err, model.ErrNotFound):
		return http.StatusNotFound, dto.ErrorResponse{
			Error:   "not_found",
			Message: "resource not found",
		}
	case errors.Is(err, model.ErrInternalError):
		return http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "internal server error",
		}
	default:
		return http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "internal server error",
		}
	}
}

func writeError(c *gin.Context, err error) {
	code, errResp := mapErrorCode(err)
	c.JSON(code, errResp)
}
