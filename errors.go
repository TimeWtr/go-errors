package errors

import (
	"fmt"
	"time"
)

// Error 定义Error的通用接口
type Error interface {
	error
	Code() string
	Message() string
	HttpStatus() int
	Type() ErrType
	Timestamp() time.Time
	StackTrace() string
	Metadata() map[string]any
	WithMetadata(key string, val any) Error
	WithMetadataMap(metadata map[string]any) Error
	Unwarp() error
}

type ErrType string

const (
	ErrTypeInternal     ErrType = "INTERNAL"
	ErrTypeBadRequest   ErrType = "BAD_REQUEST"
	ErrTypeUnauthorized ErrType = "UNAUTHORIZED"
	ErrTypeForbidden    ErrType = "FORBIDDEN"
	ErrTypeNotFound     ErrType = "NOT_FOUND"
	ErrTypeConflict     ErrType = "CONFLICT"
	ErrTypeValidation   ErrType = "VALIDATION"
	ErrTypeBusiness     ErrType = "BUSINESS"
	ErrTypeTimeout      ErrType = "TIMEOUT"
	ErrTypeRateLimit    ErrType = "RATE_LIMIT"
	ErrTypeExternal     ErrType = "EXTERNAL"
)

func (e ErrType) String() string {
	return string(e)
}

type ErrorImpl struct {
	// 错误码
	code string
	// 详细信息
	message string
	// http状态码
	httpStatus int
	// 错误类型
	errType ErrType
	// 时间戳
	timestamp time.Time
	// 堆栈信息
	stackTrace string
	// 原始的错误，error类型
	cause error
	// 其它的元数据
	metadata map[string]any
}

func (e *ErrorImpl) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s:%s", e.message, e.cause.Error())
	}

	return e.message
}

func (e *ErrorImpl) Code() string {
	return e.code
}

func (e *ErrorImpl) Message() string {
	return e.message
}

func (e *ErrorImpl) HttpStatus() int {
	return e.httpStatus
}

func (e *ErrorImpl) Type() ErrType {
	return e.errType
}

func (e *ErrorImpl) Timestamp() time.Time {
	return e.timestamp
}

func (e *ErrorImpl) StackTrace() string {
	return e.stackTrace
}

func (e *ErrorImpl) Unwarp() error {
	return e.cause
}

func (e *ErrorImpl) WithMetadata(key string, val any) Error {
	if e.metadata == nil {
		e.metadata = make(map[string]any)
	}

	e.metadata[key] = val
	return e
}

func (e *ErrorImpl) WithMetadataMap(metadata map[string]any) Error {
	if e.metadata == nil {
		e.metadata = make(map[string]any)
	}
	for k, v := range metadata {
		e.metadata[k] = v
	}

	return e
}

func (e *ErrorImpl) Metadata() map[string]any {
	return e.metadata
}
