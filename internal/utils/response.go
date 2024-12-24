package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/wavinamayola/user-management/internal/models"
)

func RespondWithError(w http.ResponseWriter, statusCode int, err error) {
	errMsg, errData := handleError(err)
	response := models.Response{
		Status:  "error",
		Message: errMsg,
		Data:    errData,
	}
	RespondWithJSON(w, statusCode, response)
}

func RespondWithJSON(w http.ResponseWriter, statusCode int, payload models.Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if payload.Status == "" {
		payload.Status = "success"
	}
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func handleError(err error) (string, interface{}) {
	errMsg := err.Error()

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var validationErrs []models.ValidationErrorResponse
		for _, err := range validationErrors {
			fmt.Println(err.Param())
			validateErrMsg := ""
			switch err.Tag() {
			case "required":
				validateErrMsg = fmt.Sprintf("%s is required", err.StructNamespace())
			case "email":
				validateErrMsg = fmt.Sprintf("%s should be email format", err.StructNamespace())
			case "min":
				validateErrMsg = fmt.Sprintf("%s must have at least %s characters", err.StructNamespace(), err.Param())
			default:
				validateErrMsg = fmt.Sprintf("%s is incorrect", err.StructNamespace())
			}
			validationErrs = append(validationErrs, models.ValidationErrorResponse{
				FailedField: err.StructNamespace(),
				Tag:         err.Tag(),
				Message:     validateErrMsg,
			})
		}

		return "validation errors", validationErrs
	}

	return errMsg, nil
}
