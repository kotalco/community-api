package errors

import (
	"net/http"
)

type IRestErr interface {
	Error() string
	StatusCode() int
}

type RestErr struct {
	Message     string            `json:"message"`
	Status      int               `json:"status"`
	Name        string            `json:"name"`
	Validations map[string]string `json:"validations,omitempty"`
}

func NewRestErr() IRestErr {
	return &RestErr{}
}

// Error used to mimic build in error pkg so this can be replaceable  for go error pkg
func (err RestErr) Error() string {
	return err.Message
}
func (err RestErr) StatusCode() int {
	return err.Status
}

func NewValidationError(validations map[string]string) IRestErr {
	return RestErr{
		Message:     "Invalid Body Request",
		Status:      http.StatusBadRequest,
		Name:        "Bad Request",
		Validations: validations,
	}
}
func NewBadRequestError(message string) IRestErr {
	return RestErr{
		Message: message,
		Status:  http.StatusBadRequest,
		Name:    "Bad Request",
	}
}
func NewNotFoundError(message string) IRestErr {
	return RestErr{
		Message: message,
		Status:  http.StatusNotFound,
		Name:    "Not Found",
	}
}
func NewInternalServerError(message string) IRestErr {
	return RestErr{
		Message: message,
		Status:  http.StatusInternalServerError,
		Name:    "Internal Server Error",
	}
}

func NewUnAuthorizedError(message string) IRestErr {
	return RestErr{
		Message: message,
		Status:  http.StatusUnauthorized,
		Name:    "UnAuthorized",
	}
}

func NewForbiddenError(message string) IRestErr {
	return RestErr{
		Message: message,
		Status:  http.StatusForbidden,
		Name:    "Forbidden",
	}
}

func NewTooManyRequestsError(message string) IRestErr {
	return RestErr{
		Message: message,
		Status:  http.StatusTooManyRequests,
		Name:    "Too Many Requests",
	}
}

func NewConflictError(message string) IRestErr {
	return RestErr{
		Message: message,
		Status:  http.StatusConflict,
		Name:    "Conflict",
	}
}
