package errors

import (
	"errors"
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
	BusinessError = &ErrCode{
		Code:       "BUSINESS_ERROR",
		Message:    "Business error",
		HttpStatus: http.StatusInternalServerError,
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
	var customErr Error
	if errors.As(err, &customErr) {
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

func Wrapf(format string, err error, code *ErrCode) Error {
	if err == nil {
		return nil
	}
	var customErr Error
	if errors.As(err, &customErr) {
		return customErr
	}

	return &ErrorImpl{
		code:       code.Code,
		message:    fmt.Sprintf(format, code.Message),
		httpStatus: code.HttpStatus,
		errType:    code.Type,
		timestamp:  time.Now().UTC(),
		stackTrace: getSimplifiedStackTrace(2, 6),
		cause:      err,
	}
}

// ==================== 基础错误函数 ====================
//

// ResourceConflictError 创建资源冲突错误
func ResourceConflictError() Error {
	return New(ErrConflict)
}

// ParameterValidationError 创建参数验证错误
func ParameterValidationError() error {
	return New(ErrBadRequest)
}

// ResourceNotFoundError 创建资源未找到错误
func ResourceNotFoundError() Error {
	return New(ErrNotFound)
}

// ForbiddenError 基础禁止访问错误
func ForbiddenError() Error {
	return New(ErrForbidden)
}

// UnAuthorizedError 创建未认证的错误
func UnAuthorizedError() Error {
	return New(ErrUnauthorized)
}

// TimeoutError 创建超时错误
func TimeoutError() Error {
	return New(ErrTimeout)
}

// InternalError 创建内部错误
func InternalError() Error {
	return New(ErrInternal)
}

// ==================== 格式化的错误函数 ====================
//

// UnAuthorizedErrorf 创建未认证错误
func UnAuthorizedErrorf(format string, args ...any) Error {
	return Newf(ErrUnauthorized, format, args...)
}

// ResourceConflictErrorf 创建资源冲突错误
func ResourceConflictErrorf(format string, args ...any) Error {
	return Newf(ErrConflict, format, args...)
}

// ParameterValidationErrorf 创建参数验证错误
func ParameterValidationErrorf(format string, args ...any) Error {
	return Newf(ErrBadRequest, format, args...)
}

// TimeoutErrorf 创建超时的错误
func TimeoutErrorf(format string, args ...any) Error {
	return Newf(ErrTimeout, format, args...)
}

// ForbiddenErrorf 创建禁止访问错误
func ForbiddenErrorf(format string, args ...any) Error {
	return Newf(ErrForbidden, format, args...)
}

// InternalErrorf 创建内部错误
func InternalErrorf(format string, args ...any) Error {
	return Newf(ErrInternal, format, args...)
}

// ==================== 带元数据的错误函数 ====================
//

// UnAuthorizedErrorWithMetadata 创建带元数据的未认证错误
func UnAuthorizedErrorWithMetadata(metadata map[string]any, format string, args ...any) Error {
	return Newf(ErrUnauthorized, format, args...).
		WithMetadataMap(metadata)
}

// ResourceConflictErrorWithMetadata 创建带元数据的资源冲突错误
func ResourceConflictErrorWithMetadata(metadata map[string]any, format string, args ...any) Error {
	return Newf(ErrConflict, format, args...).
		WithMetadataMap(metadata)
}

// ParameterValidationErrorWithMetadata 创建带元数据的参数验证错误
func ParameterValidationErrorWithMetadata(metadata map[string]any, format string, args ...any) Error {
	return Newf(ErrBadRequest, format, args...).
		WithMetadataMap(metadata)
}

// TimeoutErrorWithMetadata 创建带元数据的超时错误
func TimeoutErrorWithMetadata(metadata map[string]any, format string, args ...any) Error {
	return Newf(ErrTimeout, format, args...).
		WithMetadataMap(metadata)
}

// ForbiddenErrorWithMetadata 创建带元数据的禁止访问错误
func ForbiddenErrorWithMetadata(metadata map[string]any, format string, args ...any) Error {
	return Newf(ErrForbidden, format, args...).
		WithMetadataMap(metadata)
}

// InternalErrorWithMetadata 创建带元数据的内部错误
func InternalErrorWithMetadata(metadata map[string]any, format string, args ...any) Error {
	return Newf(ErrInternal, format, args...).
		WithMetadataMap(metadata)
}

// TODO 使用strings.Builder来实现字符串操作的性能优化
