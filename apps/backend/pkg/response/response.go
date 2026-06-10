package response

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type APIResponse struct {
	Success bool            `json:"success"`
	Message string          `json:"message,omitempty"`
	Data    interface{}     `json:"data,omitempty"`
	Meta    *PaginationMeta `json:"meta,omitempty"`
	Error   *ErrorInfo      `json:"error,omitempty"`
}

type PaginationMeta struct {
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
	Total   int64 `json:"total"`
	Pages   int64 `json:"pages"`
}

type ErrorInfo struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

func SuccessResponse(c *gin.Context, status int, message string, data interface{}, meta *PaginationMeta) {
	c.JSON(status, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func ErrorResponse(c *gin.Context, status int, code, message string) {
	c.JSON(status, APIResponse{
		Success: false,
		Error:   &ErrorInfo{Code: code, Message: message},
	})
}

// InternalError sends a 500 response. Optional custom message overrides the default.
func InternalError(c *gin.Context, msg ...string) {
	message := "Internal server error"
	if len(msg) > 0 && msg[0] != "" {
		message = msg[0]
	}

	c.JSON(http.StatusInternalServerError, APIResponse{
		Success: false,
		Error: &ErrorInfo{
			Code:    "INTERNAL_ERROR",
			Message: message,
		},
	})
}

// ValidationError sends a 422 Unprocessable Entity with field-level error details.
// This is used when request body validation fails.
func ValidationError(c *gin.Context, err error) {
	details := FormatValidationErrors(err)
	c.JSON(http.StatusUnprocessableEntity, APIResponse{
		Success: false,
		Message: "Validation failed",
		Error: &ErrorInfo{
			Code:    "VALIDATION_ERROR",
			Message: "One or more fields are invalid",
			Details: details,
		},
	})
}

// FormatValidationErrors converts a validation error into a simple map that can be
// embedded in JSON responses or logs. Expand this to unpack field-level errors
// from your validator as needed.
func FormatValidationErrors(err error) map[string]string {
	out := map[string]string{}
	var verrs validator.ValidationErrors

	// Gin uses validator.v10 under the hood for binding errors.
	if errors.As(err, &verrs) {
		for _, fe := range verrs {
			field := strings.ToLower(fe.Field())
			switch fe.Tag() {
			case "required":
				out[field] = field + " is required"
			case "maxwords":
				out[field] = field + " must be at most " + fe.Param() + " words"
			default:
				out[field] = field + " is invalid"
			}
		}
		return out
	}

	var syntaxError *json.SyntaxError
	var typeError *json.UnmarshalTypeError

	switch {
	case errors.Is(err, io.EOF):
		out["body"] = "request body is required"
	case errors.As(err, &syntaxError):
		out["body"] = "request body contains malformed JSON"
	case errors.As(err, &typeError):
		out[strings.ToLower(typeError.Field)] = typeError.Field + " has an invalid type"
	default:
		out["error"] = "Invalid request"
	}

	// Fallback for non-validation errors.
	return out
}
