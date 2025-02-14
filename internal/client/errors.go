package client

import (
	"fmt"
	"net/http"
)

// APIError represents an error returned by the Centreon API.
type APIError struct {
	StatusCode int
	Message    string
	Code       string
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("API error: %s (status: %d, code: %s)", e.Message, e.StatusCode, e.Code)
	}
	return fmt.Sprintf("API error: %s (status: %d)", e.Message, e.StatusCode)
}

// HandleAPIError creates an APIError from an HTTP response.
func HandleAPIError(resp *http.Response, body []byte) error {
	message := string(body)
	switch resp.StatusCode {
	case http.StatusBadRequest:
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    "Invalid request parameters: " + message,
			Code:       "BAD_REQUEST",
		}
	case http.StatusUnauthorized:
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    "Authentication failed",
			Code:       "UNAUTHORIZED",
		}
	case http.StatusForbidden:
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    "Access forbidden",
			Code:       "FORBIDDEN",
		}
	case http.StatusNotFound:
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    "Resource not found",
			Code:       "NOT_FOUND",
		}
	case http.StatusConflict:
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    "Resource conflict: " + message,
			Code:       "CONFLICT",
		}
	default:
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    "Unexpected error: " + message,
			Code:       "INTERNAL_ERROR",
		}
	}
}
