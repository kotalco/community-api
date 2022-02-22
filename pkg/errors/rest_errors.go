package errors

import (
	"net/http"
)

type RestErr struct {
	Message     string            `json:"message"`
	Status      int               `json:"status"`
	Error       string            `json:"error"`
	Validations map[string]string `json:"validations,omitempty"`
}

func NewValidationError(validations map[string]string) *RestErr {
	return &RestErr{
		Message:     "Invalid Body Request",
		Status:      http.StatusBadRequest,
		Error:       "Bad Request",
		Validations: validations,
	}
}
func NewBadRequestError(message string) *RestErr {
	return &RestErr{
		Message:     message,
		Status:      http.StatusBadRequest,
		Error:       "Bad Request",
		Validations: make(map[string]string),
	}
}
func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		Message:     message,
		Status:      http.StatusNotFound,
		Error:       "Not Found",
		Validations: make(map[string]string),
	}
}
func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Message:     message,
		Status:      http.StatusInternalServerError,
		Error:       "Internal Server Error",
		Validations: make(map[string]string),
	}
}

func NewUnAuthorizedError(message string) *RestErr {
	return &RestErr{
		Message:     message,
		Status:      http.StatusUnauthorized,
		Error:       "UnAuthorized",
		Validations: make(map[string]string),
	}
}
