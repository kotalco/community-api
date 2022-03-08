package errors

import (
	"net/http"
)

type RestErr struct {
	Message     string            `json:"message"`
	Status      int               `json:"status"`
	ErrorMsg      string            `json:"error"`
	Validations map[string]string `json:"validations,omitempty"`
}

//Error used to mimic build in error pkg so this can be replaceable  for go error pkg
func (err RestErr)Error() string  {
	return err.Message
}

func NewValidationError(validations map[string]string) *RestErr {
	return &RestErr{
		Message:     "Invalid Body Request",
		Status:      http.StatusBadRequest,
		ErrorMsg:       "Bad Request",
		Validations: validations,
	}
}
func NewBadRequestError(message string) *RestErr {
	return &RestErr{
		Message:     message,
		Status:      http.StatusBadRequest,
		ErrorMsg:       "Bad Request",
	}
}
func NewNotFoundError(message string) *RestErr {
	return &RestErr{
		Message:     message,
		Status:      http.StatusNotFound,
		ErrorMsg:       "Not Found",
	}
}
func NewInternalServerError(message string) *RestErr {
	return &RestErr{
		Message:     message,
		Status:      http.StatusInternalServerError,
		ErrorMsg:       "Internal Server Error",
	}
}

func NewUnAuthorizedError(message string) *RestErr {
	return &RestErr{
		Message:     message,
		Status:      http.StatusUnauthorized,
		ErrorMsg:       "UnAuthorized",
	}
}

func NewForbiddenError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusForbidden,
		ErrorMsg:   "Forbidden",
	}
}

func NewTooManyRequestsError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusTooManyRequests,
		ErrorMsg:   "Too Many Requests",
	}
}



func NewConflictError(message string) *RestErr {
	return &RestErr{
		Message: message,
		Status:  http.StatusConflict,
		ErrorMsg:   "Conflict",
	}
}
