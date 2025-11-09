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

const (
	ErrInternalMessage        = "Internal Server Error"
	ErrTimeoutMessage         = "Request Timeout"
	ErrNotFoundMessage        = "Not Found"
	ErrBadRequestMessage      = "Bad Request"
	ErrUnauthorizedMessage    = "Unauthorized"
	ErrForbiddenMessage       = "Forbidden"
	ErrConflictMessage        = "Resource Conflict"
	ErrRateLimitMessage       = "Rate Limit Exceeded"
	ErrUsernameExistedMessage = "Username already exists"
	ErrEmailExistedMessage    = "Email already exists"
	ErrPhoneExistedMessage    = "Phone already exists"
	ErrBusinessMessage        = "Business error"
)

var (
	ErrInternal = &ErrCode{
		Code:       ErrTypeInternal.String(),
		Message:    ErrInternalMessage,
		HttpStatus: http.StatusInternalServerError,
		Type:       ErrTypeInternal,
	}
	ErrTimeout = &ErrCode{
		Code:       ErrTypeTimeout.String(),
		Message:    ErrTimeoutMessage,
		HttpStatus: http.StatusRequestTimeout,
		Type:       ErrTypeTimeout,
	}
	ErrNotFound = &ErrCode{
		Code:       ErrTypeNotFound.String(),
		Message:    ErrNotFoundMessage,
		HttpStatus: http.StatusNotFound,
		Type:       ErrTypeNotFound,
	}
	ErrBadRequest = &ErrCode{
		Code:       ErrTypeBadRequest.String(),
		Message:    ErrBadRequestMessage,
		HttpStatus: http.StatusBadRequest,
		Type:       ErrTypeBadRequest,
	}
	ErrUnauthorized = &ErrCode{
		Code:       ErrTypeUnauthorized.String(),
		Message:    ErrUnauthorizedMessage,
		HttpStatus: http.StatusUnauthorized,
		Type:       ErrTypeUnauthorized,
	}
	ErrForbidden = &ErrCode{
		Code:       ErrTypeForbidden.String(),
		Message:    ErrForbiddenMessage,
		HttpStatus: http.StatusForbidden,
		Type:       ErrTypeForbidden,
	}
	ErrConflict = &ErrCode{
		Code:       ErrTypeConflict.String(),
		Message:    ErrConflictMessage,
		HttpStatus: http.StatusConflict,
		Type:       ErrTypeConflict,
	}
	ErrRateLimit = &ErrCode{
		Code:       ErrTypeRateLimit.String(),
		Message:    ErrRateLimitMessage,
		HttpStatus: http.StatusTooManyRequests,
		Type:       ErrTypeRateLimit,
	}
)

var (
	ErrUsernameExisted = &ErrCode{
		Code:       ErrTypeConflict.String(),
		Message:    ErrUsernameExistedMessage,
		HttpStatus: http.StatusConflict,
		Type:       ErrTypeConflict,
	}
	ErrEmailExisted = &ErrCode{
		Code:       ErrTypeConflict.String(),
		Message:    ErrEmailExistedMessage,
		HttpStatus: http.StatusConflict,
		Type:       ErrTypeConflict,
	}
	ErrPhoneExisted = &ErrCode{
		Code:       ErrTypeConflict.String(),
		Message:    ErrPhoneExistedMessage,
		HttpStatus: http.StatusConflict,
		Type:       ErrTypeConflict,
	}
	BusinessError = &ErrCode{
		Code:       ErrTypeBusiness.String(),
		Message:    ErrBusinessMessage,
		HttpStatus: http.StatusInternalServerError,
		Type:       ErrTypeBusiness,
	}
)

// n 创建一个新的错误实例，根据enableStack参数控制是否收集堆栈信息
// 参数:
//
//	enableStack: 布尔标志，指示是否应该收集堆栈跟踪信息
//	code:        错误码信息，包含错误代码、消息、HTTP状态码和错误类型
//
// 返回值:
//
//	Error: 新创建的错误实例
func n(enableStack bool, code *ErrCode) Error {
	impl := acquireError()
	impl.code = code.Code
	impl.message = code.Message
	impl.httpStatus = code.HttpStatus
	impl.errType = code.Type
	impl.timestamp = time.Now().UTC()

	if enableStack {
		impl.stackTrace = getSimplifiedStackTrace(2, 6)
	}

	return impl
}

// New 创建一个带堆栈信息的错误实例
// 参数:
//
//	code: 错误码信息，包含完整的错误定义
//
// 返回值:
//
//	Error: 包含堆栈信息的新错误实例
func New(code *ErrCode) Error {
	return n(true, code)
}

// FastNew 创建一个不带堆栈信息的错误实例，适用于性能敏感场景
// 参数:
//
//	code: 错误码信息，包含完整的错误定义
//
// 返回值:
//
//	Error: 不包含堆栈信息的新错误实例
func FastNew(code *ErrCode) Error {
	return n(false, code)
}

func nf(enableStack bool, code *ErrCode, format string, args ...any) Error {
	impl := acquireError()
	impl.code = code.Code
	impl.message = fmt.Sprintf(format, args...)
	impl.httpStatus = code.HttpStatus
	impl.errType = code.Type
	impl.timestamp = time.Now().UTC()

	if enableStack {
		impl.stackTrace = getSimplifiedStackTrace(2, 6)
	}

	return impl
}

// Newf 创建一个错误，带堆栈信息格式化的错误
func Newf(code *ErrCode, format string, args ...any) Error {
	return nf(true, code, format, args...)
}

func FastNewf(code *ErrCode, format string, args ...any) Error {
	return nf(false, code, format, args...)
}

// wrap 将一个标准error包装成自定义的Error类型，可选择是否包含堆栈信息
// 参数:
//
//	enableStack: 是否启用堆栈跟踪信息收集
//	err:         需要被包装的原始错误，如果为nil则返回nil
//	code:        错误码信息，包含错误代码、消息、HTTP状态码和错误类型
//
// 返回值:
//
//	Error: 包装后的自定义错误类型，如果原始错误已经是自定义类型则直接返回
func wrap(enableStack bool, err error, code *ErrCode) Error {
	if err == nil {
		return nil
	}

	// 检查错误是否已经是自定义的Error类型，如果是则直接返回
	var customErr Error
	if errors.As(err, &customErr) {
		return customErr
	}

	impl := &ErrorImpl{
		code:       code.Code,
		message:    code.Message,
		httpStatus: code.HttpStatus,
		errType:    code.Type,
		timestamp:  time.Now().UTC(),
		cause:      err,
	}

	// 根据enableStack参数决定是否记录简化版的调用堆栈
	if enableStack {
		impl.stackTrace = getSimplifiedStackTrace(2, 6)
	}

	return impl
}

// Wrap 将一个标准error包装成自定义的Error类型
// 参数:
//
//	err:  需要被包装的原始错误，如果为nil则返回nil
//	code: 错误码信息，包含错误代码、消息、HTTP状态码和错误类型
//
// 返回值:
//
//	Error: 包装后的自定义错误类型，如果原始错误已经是自定义类型则直接返回
//
// 该函数会携带调用堆栈信息
func Wrap(err error, code *ErrCode) Error {
	return wrap(true, err, code)
}

// FastWrap 将一个标准error包装成自定义的Error类型
// 参数:
//
//	err:  需要被包装的原始错误，如果为nil则返回nil
//	code: 错误码信息，包含错误代码、消息、HTTP状态码和错误类型
//
// 返回值:
//
//	Error: 包装后的自定义错误类型，如果原始错误已经是自定义类型则直接返回
//
// 该函数不会携带调用堆栈信息，适用于性能敏感场景
func FastWrap(err error, code *ErrCode) Error {
	return wrap(false, err, code)
}

// wrapf 将给定的错误包装成自定义错误类型，可选择是否包含堆栈信息
// enableStack: 是否启用堆栈跟踪信息
// err: 原始错误，如果为nil则返回nil
// format: 格式化字符串，用于格式化错误码中的消息
// code: 错误码信息，包含错误代码、HTTP状态码、错误类型等信息
// 返回值: 包装后的自定义错误类型，如果输入错误为nil则返回nil
func wrapf(enableStack bool, err error, format string, code *ErrCode) Error {
	// 如果原始错误为nil，直接返回nil
	if err == nil {
		return nil
	}

	// 检查原始错误是否已经是自定义错误类型，如果是则直接返回
	var customErr Error
	if errors.As(err, &customErr) {
		return customErr
	}

	// 创建新的错误实现，包装原始错误并添加错误码信息
	impl := &ErrorImpl{
		code:       code.Code,
		message:    fmt.Sprintf(format, code.Message),
		httpStatus: code.HttpStatus,
		errType:    code.Type,
		timestamp:  time.Now().UTC(),
		cause:      err,
	}

	if enableStack {
		impl.stackTrace = getSimplifiedStackTrace(2, 6)
	}

	return impl
}

// Wrapf 将给定的错误包装成带有堆栈信息的自定义错误类型，如果原始错误已经是自定义
// 错误类型则直接返回
// format: 格式化字符串，用于格式化错误码中的消息
// err: 原始错误，如果为nil则返回nil
// code: 错误码信息，包含错误代码、HTTP状态码、错误类型等信息
// 返回值: 包装后的自定义错误类型，如果输入错误为nil则返回nil
func Wrapf(format string, err error, code *ErrCode) Error {
	return wrapf(true, err, format, code)
}

// FastWrapf 将给定的错误包装成自定义错误类型，不包含堆栈信息，性能更高
// format: 格式化字符串，用于格式化错误码中的消息
// err: 原始错误，如果为nil则返回nil
// code: 错误码信息，包含错误代码、HTTP状态码、错误类型等信息
// 返回值: 包装后的自定义错误类型，如果输入错误为nil则返回nil
func FastWrapf(format string, err error, code *ErrCode) Error {
	return wrapf(false, err, format, code)
}

// ==================== 基础错误函数，带堆栈信息 ====================
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

// ==================== 基础错误函数，不带堆栈信息 ====================
//

// ResourceConflictErrorNoStack 创建资源冲突错误
func ResourceConflictErrorNoStack() Error {
	return FastNew(ErrConflict)
}

// ParameterValidationErrorNoStack 创建参数验证错误
func ParameterValidationErrorNoStack() error {
	return FastNew(ErrBadRequest)
}

// ResourceNotFoundErrorNoStack 创建资源未找到错误
func ResourceNotFoundErrorNoStack() Error {
	return FastNew(ErrNotFound)
}

// ForbiddenErrorNoStack 基础禁止访问错误
func ForbiddenErrorNoStack() Error {
	return FastNew(ErrForbidden)
}

// UnAuthorizedErrorNoStack 创建未认证的错误
func UnAuthorizedErrorNoStack() Error {
	return FastNew(ErrUnauthorized)
}

// TimeoutErrorNoStack 创建超时错误
func TimeoutErrorNoStack() Error {
	return FastNew(ErrTimeout)
}

// InternalErrorNoStack 创建内部错误
func InternalErrorNoStack() Error {
	return New(ErrInternal)
}

// ==================== 格式化的错误函数，带堆栈信息 ====================
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

// ==================== 格式化的错误函数，不带堆栈信息 ====================
//

// UnAuthorizedErrorfNoStack 创建未认证错误
func UnAuthorizedErrorfNoStack(format string, args ...any) Error {
	return FastNewf(ErrUnauthorized, format, args...)
}

// ResourceConflictErrorfNoStack 创建资源冲突错误
func ResourceConflictErrorfNoStack(format string, args ...any) Error {
	return FastNewf(ErrConflict, format, args...)
}

// ParameterValidationErrorfNoStack 创建参数验证错误
func ParameterValidationErrorfNoStack(format string, args ...any) Error {
	return FastNewf(ErrBadRequest, format, args...)
}

// TimeoutErrorfNoStack 创建超时的错误
func TimeoutErrorfNoStack(format string, args ...any) Error {
	return FastNewf(ErrTimeout, format, args...)
}

// ForbiddenErrorfNoStack 创建禁止访问错误
func ForbiddenErrorfNoStack(format string, args ...any) Error {
	return FastNewf(ErrForbidden, format, args...)
}

// InternalErrorfNoStack 创建内部错误
func InternalErrorfNoStack(format string, args ...any) Error {
	return FastNewf(ErrInternal, format, args...)
}

// ==================== 带元数据的错误函数，带堆栈信息 ====================
//

// UnAuthorizedErrorWithMeta 创建带元数据的未认证错误
func UnAuthorizedErrorWithMeta(metadata map[string]any, format string, args ...any) Error {
	return Newf(ErrUnauthorized, format, args...).
		WithMetadataMap(metadata)
}

// ResourceConflictErrorWithMeta 创建带元数据的资源冲突错误
func ResourceConflictErrorWithMeta(metadata map[string]any, format string, args ...any) Error {
	return Newf(ErrConflict, format, args...).
		WithMetadataMap(metadata)
}

// ParameterValidationErrorWithMeta 创建带元数据的参数验证错误
func ParameterValidationErrorWithMeta(metadata map[string]any, format string, args ...any) Error {
	return Newf(ErrBadRequest, format, args...).
		WithMetadataMap(metadata)
}

// TimeoutErrorWithMeta 创建带元数据的超时错误
func TimeoutErrorWithMeta(metadata map[string]any, format string, args ...any) Error {
	return Newf(ErrTimeout, format, args...).
		WithMetadataMap(metadata)
}

// ForbiddenErrorWithMeta 创建带元数据的禁止访问错误
func ForbiddenErrorWithMeta(metadata map[string]any, format string, args ...any) Error {
	return Newf(ErrForbidden, format, args...).
		WithMetadataMap(metadata)
}

// InternalErrorWithMeta 创建带元数据的内部错误
func InternalErrorWithMeta(metadata map[string]any, format string, args ...any) Error {
	return Newf(ErrInternal, format, args...).
		WithMetadataMap(metadata)
}

// ==================== 带元数据的错误函数，不带堆栈信息 ====================
//

// UnAuthorizedErrorWithMetaNoStack 创建带元数据的未认证错误
func UnAuthorizedErrorWithMetaNoStack(metadata map[string]any, format string, args ...any) Error {
	return FastNewf(ErrUnauthorized, format, args...).
		WithMetadataMap(metadata)
}

// ResourceConflictErrorWithMetaNoStack 创建带元数据的资源冲突错误
func ResourceConflictErrorWithMetaNoStack(metadata map[string]any, format string, args ...any) Error {
	return FastNewf(ErrConflict, format, args...).
		WithMetadataMap(metadata)
}

// ParameterValidationErrorWithMetaNoStack 创建带元数据的参数验证错误
func ParameterValidationErrorWithMetaNoStack(metadata map[string]any, format string, args ...any) Error {
	return FastNewf(ErrBadRequest, format, args...).
		WithMetadataMap(metadata)
}

// TimeoutErrorWithMetaNoStack 创建带元数据的超时错误
func TimeoutErrorWithMetaNoStack(metadata map[string]any, format string, args ...any) Error {
	return FastNewf(ErrTimeout, format, args...).
		WithMetadataMap(metadata)
}

// ForbiddenErrorWithMetaNoStack 创建带元数据的禁止访问错误
func ForbiddenErrorWithMetaNoStack(metadata map[string]any, format string, args ...any) Error {
	return FastNewf(ErrForbidden, format, args...).
		WithMetadataMap(metadata)
}

// InternalErrorWithMetaNoStack 创建带元数据的内部错误
func InternalErrorWithMetaNoStack(metadata map[string]any, format string, args ...any) Error {
	return FastNewf(ErrInternal, format, args...).
		WithMetadataMap(metadata)
}
