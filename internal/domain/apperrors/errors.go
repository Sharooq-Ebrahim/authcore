package apperrors

import "fmt"

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NewAppError(code, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

var (
	ErrAccountAlreadyExists = NewAppError("ACCOUNT_EXISTS", "Account already exists")
	ErrInvalidCredentials  = NewAppError("INVALID_CREDENTIALS", "Invalid email or password")
	ErrUserNotFound        = NewAppError("USER_NOT_FOUND", "User not found")
	ErrInternalServer      = NewAppError("INTERNAL_SERVER_ERROR", "An unexpected error occurred")
	ErrInvalidToken        = NewAppError("INVALID_TOKEN", "The provided token is invalid or expired")
	ErrUnauthorized        = NewAppError("UNAUTHORIZED", "You are not authorized to perform this action")
	ErrInvalidInput        = NewAppError("INVALID_INPUT", "The request input is invalid")
	ErrTokenInvalid        = NewAppError("TOKEN_INVALID", "Invalid token")
	ErrTokenRevoked        = NewAppError("TOKEN_REVOKED", "Token has been revoked")
	ErrForbidden           = NewAppError("FORBIDDEN", "Access forbidden")
	ErrNotFound            = NewAppError("NOT_FOUND", "Resource not found")
	ErrConflict            = NewAppError("CONFLICT", "Resource conflict")
	ErrInvalidRole         = NewAppError("INVALID_ROLE", "Invalid role provided")
	ErrBadRequest          = NewAppError("BAD_REQUEST", "Bad request")
	ErrInvalidCredential   = NewAppError("INVALID_CREDENTIAL", "Invalid credential")
)
