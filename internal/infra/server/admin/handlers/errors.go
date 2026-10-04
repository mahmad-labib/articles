package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/mahmad-labib/articles/internal/service"
	"github.com/mahmad-labib/articles/internal/service/validation"
	"gorm.io/gorm"
)

var ErrInvalidId = errors.New("id should be an integer")

type ErrorResponse struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
}

func handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, validation.EmailErr),
		errors.Is(err, validation.Empty),
		errors.Is(err, validation.NameLength),
		errors.Is(err, validation.EmailUnique),
		errors.Is(err, service.PasswordInvalid),
		errors.Is(err, gorm.ErrInvalidValue),
		errors.Is(err, validation.FieldSpaces):
		writeError(w, http.StatusUnprocessableEntity, "Validation Error", []string{err.Error()})

	case errors.Is(err, gorm.ErrRecordNotFound),
		errors.Is(err, gorm.ErrForeignKeyViolated):
		writeError(w, http.StatusNotFound, "Not Found", []string{err.Error()})

	case errors.Is(err, gorm.ErrDuplicatedKey):
		writeError(w, http.StatusConflict, "item already exist", []string{err.Error()})

	default:
		writeError(w, http.StatusInternalServerError, "internal server error", []string{})
	}
}

func handleRequestError(w http.ResponseWriter, err error) {
	writeError(w, http.StatusBadRequest, "Bad Request", []string{err.Error()})
}

func writeError(w http.ResponseWriter, status int, message string, details []string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(ErrorResponse{
		Code:    status,
		Message: message,
		Details: details,
	})
}
