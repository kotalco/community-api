package errors

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)


func TestNewBadRequestError(t *testing.T) {
	err:=NewBadRequestError("Bad Request")
	assert.EqualValues(t, err.Message,"Bad Request")
	assert.EqualValues(t, http.StatusBadRequest,err.Status)
}
func TestNewNotFoundError(t *testing.T) {
	err:=NewNotFoundError("Not Found")
	assert.EqualValues(t, err.Message,"Not Found")
	assert.EqualValues(t, http.StatusNotFound,err.Status)
}

func TestNewInternalServerError(t *testing.T) {
	err:=NewInternalServerError("internal server error")
	assert.EqualValues(t, err.Message,"internal server error")
	assert.EqualValues(t,http.StatusInternalServerError, err.Status)
}