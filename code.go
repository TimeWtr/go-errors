package errors

import (
	"fmt"
	"net/http"
	"time"
)

// ErrCode 错误码结构定义
type ErrCode struct {
	Code       string
	Message    string
	HttpStatus int
	Type       ErrType
}

var (
	ErrInternal = &ErrCode{
		Code:       ErrTypeInternal.String(),
		Message:    "Internal Server Error",
		HttpStatus: http.StatusInternalServerError,
		Type:       ErrTypeInternal,
	}
	ErrTimeout = &ErrCode{
		Code:       ErrTypeTimeout.String(),
		Message:    "Request Timeout",
		HttpStatus: http.StatusRequestTimeout,
		Type:       ErrTypeTimeout,
	}
	ErrNotFound = &ErrCode{
		Code:       ErrTypeNotFound.String(),
		Message:    "Not Found",
		HttpStatus: http.StatusNotFound,
		Type:       ErrTypeNotFound,
	}
	ErrBadRequest = &ErrCode{
		Code:       ErrTypeBadRequest.String(),
		Message:    "Bad Request",
		HttpStatus: http.StatusBadRequest,
		Type:       ErrTypeBadRequest,
	}
	ErrUnauthorized = &ErrCode{
		Code:       ErrTypeUnauthorized.String(),
		Message:    "Unauthorized",
		HttpStatus: http.StatusUnauthorized,
		Type:       ErrTypeUnauthorized,
	}
	ErrForbidden = &ErrCode{
		Code:       ErrTypeForbidden.String(),
		Message:    "Forbidden",
		HttpStatus: http.StatusForbidden,
		Type:       ErrTypeForbidden,
	}
	ErrConflict = &ErrCode{
		Code:       ErrTypeConflict.String(),
		Message:    "Resource Conflict",
		HttpStatus: http.StatusConflict,
		Type:       ErrTypeConflict,
	}
	ErrBusiness = &ErrCode{
		Code:       ErrTypeBusiness.String(),
		Message:    "Business rule violation",
		HttpStatus: http.StatusBadRequest,
		Type:       ErrTypeBusiness,
	}
	ErrValidation = &ErrCode{
		Code:       ErrTypeValidation.String(),
		Message:    "Validation error",
		HttpStatus: http.StatusBadRequest,
		Type:       ErrTypeValidation,
	}
	ErrRateLimit = &ErrCode{
		Code:       "RATE_LIMIT",
		Message:    "Rate limit exceeded",
		HttpStatus: http.StatusTooManyRequests,
		Type:       ErrTypeBusiness,
	}
)

var (
	ErrUsernameExisted = &ErrCode{
		Code:       "USERNAME_EXISTED",
		Message:    "Username already exists",
		HttpStatus: http.StatusConflict,
		Type:       ErrTypeBusiness,
	}
	ErrEmailExisted = &ErrCode{
		Code:       "EMAIL_EXISTED",
		Message:    "Email already exists",
		HttpStatus: http.StatusConflict,
		Type:       ErrTypeBusiness,
	}
	ErrPhoneExisted = &ErrCode{
		Code:       "PHONE_EXISTED",
		Message:    "Phone already exists",
		HttpStatus: http.StatusConflict,
		Type:       ErrTypeBusiness,
	}
)

func New(code *ErrCode) Error {
	return &ErrorImpl{
		code:       code.Code,
		message:    code.Message,
		httpStatus: code.HttpStatus,
		errType:    code.Type,
		timestamp:  time.Now().UTC(),
		stackTrace: getSimplifiedStackTrace(2, 6),
		cause:      nil,
		metadata:   nil,
	}
}

func Newf(code *ErrCode, format string, args ...any) Error {
	return &ErrorImpl{
		code:       code.Code,
		message:    fmt.Sprintf(format, args...),
		httpStatus: code.HttpStatus,
		errType:    code.Type,
		timestamp:  time.Now().UTC(),
		stackTrace: getSimplifiedStackTrace(2, 6),
	}
}

func Wrap(err error, code *ErrCode) Error {
	if err == nil {
		return nil
	}
	if customErr, ok := err.(Error); ok {
		return customErr
	}

	return &ErrorImpl{
		code:       code.Code,
		message:    code.Message,
		httpStatus: code.HttpStatus,
		errType:    code.Type,
		timestamp:  time.Now().UTC(),
		stackTrace: getSimplifiedStackTrace(2, 6),
		cause:      err,
	}
}

func UsernameExistedError(username string) Error {
	return Newf(ErrUsernameExisted,
		"Username %s already exists", username).
		WithMetadata("username", username)
}

func PhoneExistedError(phone string) Error {
	return Newf(ErrPhoneExisted,
		"Phone %s already exists", phone).
		WithMetadata("phone", phone)
}

func EmailExistedError(email string) Error {
	return Newf(ErrEmailExisted,
		"Email %s already exists", email).
		WithMetadata("email", email)
}

func ValidationError(field, reason string) error {
	return New(ErrValidation).
		WithMetadata("field", field).
		WithMetadata("reason", reason)
}
