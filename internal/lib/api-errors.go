package lib

import (
	"fmt"
	"net/http"
)

type APIError struct {
	StatusCode int `json:"statusCode"`
	Msg        any `json:"msg"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("%v", e.Msg)
}

func NewAPIError(statusCode int, err error) APIError {
	return APIError{
		StatusCode: statusCode,
		Msg:        err.Error(),
	}
}

func InvalidJSON() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid JSON request"))
}

func TaskNotFound(id int) APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("task not found with ID: %d", id))
}

func NothigToUpdate() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("nothing to update"))
}

func InvalidID() APIError {
	return NewAPIError(http.StatusBadRequest, fmt.Errorf("invalid ID"))
}

func InternalServerError() APIError {
	return NewAPIError(http.StatusInternalServerError, fmt.Errorf("internal server error"))
}

func InvalidRequestData(errors map[string]string) APIError {
	return APIError{
		StatusCode: http.StatusUnprocessableEntity,
		Msg:        errors,
	}
}
