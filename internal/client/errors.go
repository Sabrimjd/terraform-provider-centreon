package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// APIError represents an error returned by the Centreon API.
type APIError struct {
	StatusCode int
	Message    string
	Code       string
	RawBody    string
}

// ErrorResponse represents the JSON error structure returned by Centreon API.
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("API error: %s (status: %d, code: %s)", e.Message, e.StatusCode, e.Code)
	}
	return fmt.Sprintf("API error: %s (status: %d)", e.Message, e.StatusCode)
}

// IsNotFound reports whether err is a 404 from the Centreon API.
func IsNotFound(err error) bool {
	apiErr, ok := err.(*APIError)
	return ok && apiErr.StatusCode == http.StatusNotFound
}

// HandleAPIError creates an APIError from an HTTP response.
func HandleAPIError(resp *http.Response, body []byte) error {
	ctx := context.Background()
	rawMessage := string(body)

	// Try to parse as JSON error
	var errorResp ErrorResponse
	message := rawMessage

	if err := json.Unmarshal(body, &errorResp); err == nil && errorResp.Message != "" {
		message = errorResp.Message
		tflog.Warn(ctx, "Centreon API error response", map[string]interface{}{
			"status_code": resp.StatusCode,
			"error_code":  errorResp.Code,
			"message":     errorResp.Message,
		})
	} else {
		tflog.Warn(ctx, "Centreon API error response (unparsable body)", map[string]interface{}{
			"status_code": resp.StatusCode,
			"raw_body":    rawMessage,
		})
	}

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return &APIError{StatusCode: resp.StatusCode, Message: "Invalid request parameters: " + message, Code: "BAD_REQUEST", RawBody: rawMessage}
	case http.StatusUnauthorized:
		return &APIError{StatusCode: resp.StatusCode, Message: "Authentication failed", Code: "UNAUTHORIZED", RawBody: rawMessage}
	case http.StatusForbidden:
		return &APIError{StatusCode: resp.StatusCode, Message: "Access forbidden", Code: "FORBIDDEN", RawBody: rawMessage}
	case http.StatusNotFound:
		return &APIError{StatusCode: resp.StatusCode, Message: "Resource not found", Code: "NOT_FOUND", RawBody: rawMessage}
	case http.StatusConflict:
		return &APIError{StatusCode: resp.StatusCode, Message: "Resource conflict: " + message, Code: "CONFLICT", RawBody: rawMessage}
	default:
		return &APIError{StatusCode: resp.StatusCode, Message: "Unexpected error: " + message, Code: "INTERNAL_ERROR", RawBody: rawMessage}
	}
}
