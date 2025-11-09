// Copyright 2025 TimeWtr
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errors

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"
)

// TestBasicErrorFunctions 测试基础错误函数
func TestBasicErrorFunctions(t *testing.T) {
	tests := []struct {
		name       string
		errFunc    func() Error
		wantCode   string
		wantType   ErrType
		wantStatus int
	}{
		{
			name:       "ResourceConflictError",
			errFunc:    ResourceConflictError,
			wantCode:   ErrTypeConflict.String(),
			wantType:   ErrTypeConflict,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "ResourceNotFoundError",
			errFunc:    ResourceNotFoundError,
			wantCode:   ErrTypeNotFound.String(),
			wantType:   ErrTypeNotFound,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "ForbiddenError",
			errFunc:    ForbiddenError,
			wantCode:   ErrTypeForbidden.String(),
			wantType:   ErrTypeForbidden,
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "UnAuthorizedError",
			errFunc:    UnAuthorizedError,
			wantCode:   ErrTypeUnauthorized.String(),
			wantType:   ErrTypeUnauthorized,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "TimeoutError",
			errFunc:    TimeoutError,
			wantCode:   ErrTypeTimeout.String(),
			wantType:   ErrTypeTimeout,
			wantStatus: http.StatusRequestTimeout,
		},
		{
			name:       "InternalError",
			errFunc:    InternalError,
			wantCode:   ErrTypeInternal.String(),
			wantType:   ErrTypeInternal,
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.errFunc()

			if err.Code() != tt.wantCode {
				t.Errorf("Code() = %v, want %v", err.Code(), tt.wantCode)
			}

			if err.Type() != tt.wantType {
				t.Errorf("Type() = %v, want %v", err.Type(), tt.wantType)
			}

			if err.HttpStatus() != tt.wantStatus {
				t.Errorf("HttpStatus() = %v, want %v", err.HttpStatus(), tt.wantStatus)
			}

			if err.Message() == "" {
				t.Error("Message() should not be empty")
			}

			if err.Timestamp().IsZero() {
				t.Error("Timestamp() should not be zero")
			}

			if err.StackTrace() == "" {
				t.Error("StackTrace() should not be empty")
			}

			// 测试错误接口
			if err.Error() == "" {
				t.Error("Error() should not be empty")
			}

			// 测试 Unwrap 方法
			if err.Unwrap() != nil {
				t.Error("Unwrap() should return nil for basic errors")
			}
		})
	}
}

// TestBasicErrorFunctionsNoStack 测试基础错误函数，不带堆栈信息
func TestBasicErrorFunctionsNoStack(t *testing.T) {
	tests := []struct {
		name       string
		errFunc    func() Error
		wantCode   string
		wantType   ErrType
		wantStatus int
	}{
		{
			name:       "ResourceConflictError",
			errFunc:    ResourceConflictErrorNoStack,
			wantCode:   ErrTypeConflict.String(),
			wantType:   ErrTypeConflict,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "ResourceNotFoundError",
			errFunc:    ResourceNotFoundErrorNoStack,
			wantCode:   ErrTypeNotFound.String(),
			wantType:   ErrTypeNotFound,
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "ForbiddenError",
			errFunc:    ForbiddenErrorNoStack,
			wantCode:   ErrTypeForbidden.String(),
			wantType:   ErrTypeForbidden,
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "UnAuthorizedError",
			errFunc:    UnAuthorizedErrorNoStack,
			wantCode:   ErrTypeUnauthorized.String(),
			wantType:   ErrTypeUnauthorized,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "TimeoutError",
			errFunc:    TimeoutErrorNoStack,
			wantCode:   ErrTypeTimeout.String(),
			wantType:   ErrTypeTimeout,
			wantStatus: http.StatusRequestTimeout,
		},
		{
			name:       "InternalError",
			errFunc:    InternalErrorNoStack,
			wantCode:   ErrTypeInternal.String(),
			wantType:   ErrTypeInternal,
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.errFunc()

			if err.Code() != tt.wantCode {
				t.Errorf("Code() = %v, want %v", err.Code(), tt.wantCode)
			}

			if err.Type() != tt.wantType {
				t.Errorf("Type() = %v, want %v", err.Type(), tt.wantType)
			}

			if err.HttpStatus() != tt.wantStatus {
				t.Errorf("HttpStatus() = %v, want %v", err.HttpStatus(), tt.wantStatus)
			}

			if err.Message() == "" {
				t.Error("Message() should not be empty")
			}

			if err.Timestamp().IsZero() {
				t.Error("Timestamp() should not be zero")
			}

			// 测试错误接口
			if err.Error() == "" {
				t.Error("Error() should not be empty")
			}

			// 测试 Unwrap 方法
			if err.Unwrap() != nil {
				t.Error("Unwrap() should return nil for basic errors")
			}
		})
	}
}

// TestParameterValidationError 测试参数验证错误（返回标准 error 接口）
func TestParameterValidationError(t *testing.T) {
	err := ParameterValidationError()
	if err == nil {
		t.Fatal("ParameterValidationError() should not return nil")
	}

	// 转换为自定义的 Error 类型进行详细检查
	var customErr Error
	ok := errors.As(err, &customErr)
	if !ok {
		t.Fatal("ParameterValidationError() should return Error type")
	}

	if customErr.Code() != ErrTypeBadRequest.String() {
		t.Errorf("Code() = %v, want %v", customErr.Code(), ErrTypeBadRequest.String())
	}

	if customErr.Type() != ErrTypeBadRequest {
		t.Errorf("Type() = %v, want %v", customErr.Type(), ErrTypeBadRequest)
	}

	if customErr.HttpStatus() != http.StatusBadRequest {
		t.Errorf("HttpStatus() = %v, want %v", customErr.HttpStatus(), http.StatusBadRequest)
	}
}

// TestFormattedErrorFunctions 测试格式化错误函数
func TestFormattedErrorFunctions(t *testing.T) {
	tf := time.Now().Format(time.RFC3339)
	tests := []struct {
		name        string
		errFunc     func(format string, args ...any) Error
		format      string
		args        []any
		wantCode    string
		wantType    ErrType
		wantMessage string
	}{
		{
			name:        "UnAuthorizedErrorf",
			errFunc:     UnAuthorizedErrorf,
			format:      "User %s token expired at %v",
			args:        []any{"john_doe", tf},
			wantCode:    ErrTypeUnauthorized.String(),
			wantType:    ErrTypeUnauthorized,
			wantMessage: fmt.Sprintf("User john_doe token expired at %s", tf),
		},
		{
			name:        "ResourceConflictErrorf",
			errFunc:     ResourceConflictErrorf,
			format:      "Resource %s with ID %s already exists",
			args:        []any{"user", "12345"},
			wantCode:    ErrTypeConflict.String(),
			wantType:    ErrTypeConflict,
			wantMessage: "Resource user with ID 12345 already exists",
		},
		{
			name:        "ParameterValidationErrorf",
			errFunc:     ParameterValidationErrorf,
			format:      "Field %s validation failed: %s",
			args:        []any{"email", "invalid format"},
			wantCode:    ErrTypeBadRequest.String(),
			wantType:    ErrTypeBadRequest,
			wantMessage: "Field email validation failed: invalid format",
		},
		{
			name:        "TimeoutErrorf",
			errFunc:     TimeoutErrorf,
			format:      "Operation %s timed out after %v",
			args:        []any{"database_query", 30 * time.Second},
			wantCode:    ErrTypeTimeout.String(),
			wantType:    ErrTypeTimeout,
			wantMessage: "Operation database_query timed out after 30s",
		},
		{
			name:        "ForbiddenErrorf",
			errFunc:     ForbiddenErrorf,
			format:      "User %s cannot access resource %s",
			args:        []any{"user123", "document_456"},
			wantCode:    ErrTypeForbidden.String(),
			wantType:    ErrTypeForbidden,
			wantMessage: "User user123 cannot access resource document_456",
		},
		{
			name:        "InternalErrorf",
			errFunc:     InternalErrorf,
			format:      "Failed to process %s: %v",
			args:        []any{"request", "connection refused"},
			wantCode:    ErrTypeInternal.String(),
			wantType:    ErrTypeInternal,
			wantMessage: "Failed to process request: connection refused",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.errFunc(tt.format, tt.args...)

			if err.Code() != tt.wantCode {
				t.Errorf("Code() = %v, want %v", err.Code(), tt.wantCode)
			}

			if err.Type() != tt.wantType {
				t.Errorf("Type() = %v, want %v", err.Type(), tt.wantType)
			}

			if err.Message() != tt.wantMessage {
				t.Errorf("Message() = %v, want %v", err.Message(), tt.wantMessage)
			}

			// 测试错误消息包含预期内容
			if err.Error() != tt.wantMessage {
				t.Errorf("Error() = %v, want %v", err.Error(), tt.wantMessage)
			}
		})
	}
}

// TestErrorWithMetadataFunctions 测试带元数据的错误函数
func TestErrorWithMetadataFunctions(t *testing.T) {
	tests := []struct {
		name         string
		errFunc      func(metadata map[string]any, format string, args ...any) Error
		metadata     map[string]any
		format       string
		args         []any
		wantCode     string
		wantType     ErrType
		wantMessage  string
		wantMetadata map[string]any
	}{
		{
			name:    "UnAuthorizedErrorWithMeta",
			errFunc: UnAuthorizedErrorWithMeta,
			metadata: map[string]any{
				"user_id":    "user123",
				"token_type": "bearer",
				"expired_at": time.Now().Add(-1 * time.Hour),
			},
			format:      "Token for user %s expired",
			args:        []any{"user123"},
			wantCode:    ErrTypeUnauthorized.String(),
			wantType:    ErrTypeUnauthorized,
			wantMessage: "Token for user user123 expired",
			wantMetadata: map[string]any{
				"user_id":    "user123",
				"token_type": "bearer",
				"expired_at": time.Time{},
			},
		},
		{
			name:    "ResourceConflictErrorWithMeta",
			errFunc: ResourceConflictErrorWithMeta,
			metadata: map[string]any{
				"resource_type": "user",
				"resource_id":   "12345",
				"conflict_with": "email@example.com",
			},
			format:      "Resource %s conflict",
			args:        []any{"user"},
			wantCode:    ErrTypeConflict.String(),
			wantType:    ErrTypeConflict,
			wantMessage: "Resource user conflict",
			wantMetadata: map[string]any{
				"resource_type": "user",
				"resource_id":   "12345",
				"conflict_with": "email@example.com",
			},
		},
		{
			name:    "ParameterValidationErrorWithMeta",
			errFunc: ParameterValidationErrorWithMeta,
			metadata: map[string]any{
				"field":  "email",
				"value":  "invalid-email",
				"reason": "must be valid email format",
			},
			format:      "Validation failed for field %s",
			args:        []any{"email"},
			wantCode:    ErrTypeBadRequest.String(),
			wantType:    ErrTypeBadRequest,
			wantMessage: "Validation failed for field email",
			wantMetadata: map[string]any{
				"field":  "email",
				"value":  "invalid-email",
				"reason": "must be valid email format",
			},
		},
		{
			name:    "TimeoutErrorWithMeta",
			errFunc: TimeoutErrorWithMeta,
			metadata: map[string]any{
				"operation": "database_query",
				"timeout":   30 * time.Second,
				"attempts":  3,
			},
			format:      "Operation %s timed out",
			args:        []any{"database_query"},
			wantCode:    ErrTypeTimeout.String(),
			wantType:    ErrTypeTimeout,
			wantMessage: "Operation database_query timed out",
			wantMetadata: map[string]any{
				"operation": "database_query",
				"timeout":   30 * time.Second,
				"attempts":  3,
			},
		},
		{
			name:    "ForbiddenErrorWithMeta",
			errFunc: ForbiddenErrorWithMeta,
			metadata: map[string]any{
				"user_id":    "user123",
				"resource":   "document_456",
				"permission": "write",
				"user_roles": []string{"viewer"},
			},
			format:      "Access denied to %s",
			args:        []any{"document_456"},
			wantCode:    ErrTypeForbidden.String(),
			wantType:    ErrTypeForbidden,
			wantMessage: "Access denied to document_456",
			wantMetadata: map[string]any{
				"user_id":    "user123",
				"resource":   "document_456",
				"permission": "write",
				"user_roles": []string{"viewer"},
			},
		},
		{
			name:    "InternalErrorWithMeta",
			errFunc: InternalErrorWithMeta,
			metadata: map[string]any{
				"component":   "database",
				"operation":   "create_user",
				"error_code":  "23505",
				"stack_trace": "some stack trace...",
			},
			format:      "Internal error in %s",
			args:        []any{"database"},
			wantCode:    ErrTypeInternal.String(),
			wantType:    ErrTypeInternal,
			wantMessage: "Internal error in database",
			wantMetadata: map[string]any{
				"component":   "database",
				"operation":   "create_user",
				"error_code":  "23505",
				"stack_trace": "some stack trace...",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.errFunc(tt.metadata, tt.format, tt.args...)

			// 验证基本属性
			if err.Code() != tt.wantCode {
				t.Errorf("Code() = %v, want %v", err.Code(), tt.wantCode)
			}

			if err.Type() != tt.wantType {
				t.Errorf("Type() = %v, want %v", err.Type(), tt.wantType)
			}

			if err.Message() != tt.wantMessage {
				t.Errorf("Message() = %v, want %v", err.Message(), tt.wantMessage)
			}

			// 验证元数据
			metadata := err.Metadata()
			if len(metadata) != len(tt.wantMetadata) {
				t.Errorf("Metadata length = %v, want %v", len(metadata), len(tt.wantMetadata))
			}

			for key, wantValue := range tt.wantMetadata {
				actualValue, exists := metadata[key]
				if !exists {
					t.Errorf("Metadata key %s not found", key)
					continue
				}

				// 特殊处理 time.Time 类型，因为时间比较需要特殊处理
				if key == "expired_at" {
					if _, isTime := actualValue.(time.Time); !isTime {
						t.Errorf("Metadata key %s should be time.Time, got %T", key, actualValue)
					}
					continue
				}

				// 比较其他类型的值
				if fmt.Sprintf("%v", actualValue) != fmt.Sprintf("%v", wantValue) {
					t.Errorf("Metadata[%s] = %v, want %v", key, actualValue, wantValue)
				}
			}
		})
	}
}

// TestErrorChaining 测试错误链功能
func TestErrorChaining(t *testing.T) {
	// 创建原始错误
	originalErr := fmt.Errorf("database connection failed: connection refused")

	// 包装错误
	wrappedErr := Wrap(originalErr, ErrInternal)

	// 验证包装错误
	if wrappedErr.Code() != ErrTypeInternal.String() {
		t.Errorf("Wrapped error Code() = %v, want %v", wrappedErr.Code(), ErrTypeInternal.String())
	}

	// 验证错误链
	unwrapped := wrappedErr.Unwrap()
	if unwrapped == nil {
		t.Fatal("Unwrap() should return the original error")
	}

	if unwrapped.Error() != originalErr.Error() {
		t.Errorf("Unwrap().Error() = %v, want %v", unwrapped.Error(), originalErr.Error())
	}

	// 测试 errors.Is 兼容性
	if !errors.Is(wrappedErr, originalErr) {
		t.Error("errors.Is should return true for wrapped error and original error")
	}
}

// TestErrorChainingWithFormat 测试带格式化的错误链
func TestErrorChainingWithFormat(t *testing.T) {
	originalErr := fmt.Errorf("file not found: /path/to/file.txt")

	wrappedErr := Wrapf("additional context: %s", originalErr, ErrNotFound)

	if wrappedErr.Code() != ErrTypeNotFound.String() {
		t.Errorf("Code() = %v, want %v", wrappedErr.Code(), ErrTypeNotFound.String())
	}

	// 注意：Wrapf 的格式化是针对错误消息的
	if wrappedErr.Message() != fmt.Sprintf("additional context: %s", ErrNotFound.Message) {
		t.Errorf("Message() = %v, want %v", wrappedErr.Message(),
			fmt.Sprintf("additional context: %s", ErrNotFound.Message))
	}

	if wrappedErr.Unwrap().Error() != originalErr.Error() {
		t.Errorf("Unwrap().Error() = %v, want %v", wrappedErr.Unwrap().Error(), originalErr.Error())
	}
}

// TestNilErrorHandling 测试对 nil 错误的处理
func TestNilErrorHandling(t *testing.T) {
	// 测试 Wrap 和 Wrapf 对 nil 的处理
	if Wrap(nil, ErrInternal) != nil {
		t.Error("Wrap(nil, errCode) should return nil")
	}

	if Wrapf("format", nil, ErrInternal) != nil {
		t.Error("Wrapf(nil, errCode) should return nil")
	}
}

// TestMetadataManipulation 测试元数据操作
func TestMetadataManipulation(t *testing.T) {
	err := InternalError()

	// 测试 WithMetadata
	errWithMeta := err.WithMetadata("key1", "value1")
	if errWithMeta.Metadata()["key1"] != "value1" {
		t.Error("WithMetadata failed to add metadata")
	}

	// 测试 WithMetadataMap
	errWithMetaMap := err.WithMetadataMap(map[string]any{
		"key2": "value2",
		"key3": 123,
	})

	metadata := errWithMetaMap.Metadata()
	if metadata["key2"] != "value2" {
		t.Error("WithMetadataMap failed to add key2")
	}
	if metadata["key3"] != 123 {
		t.Error("WithMetadataMap failed to add key3")
	}

	// 测试链式调用
	chainErr := err.
		WithMetadata("chain1", "value1").
		WithMetadata("chain2", "value2").
		WithMetadataMap(map[string]any{"chain3": "value3"})

	chainMetadata := chainErr.Metadata()
	if chainMetadata["chain1"] != "value1" || chainMetadata["chain2"] != "value2" || chainMetadata["chain3"] != "value3" {
		t.Error("Chain metadata manipulation failed")
	}
}

// TestPredefinedBusinessErrors 测试预定义的业务错误
func TestPredefinedBusinessErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        *ErrCode
		wantCode   string
		wantType   ErrType
		wantStatus int
	}{
		{
			name:       "ErrUsernameExisted",
			err:        ErrUsernameExisted,
			wantCode:   "USERNAME_EXISTED",
			wantType:   ErrTypeBusiness,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "ErrEmailExisted",
			err:        ErrEmailExisted,
			wantCode:   "EMAIL_EXISTED",
			wantType:   ErrTypeBusiness,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "ErrPhoneExisted",
			err:        ErrPhoneExisted,
			wantCode:   "PHONE_EXISTED",
			wantType:   ErrTypeBusiness,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "BusinessError",
			err:        BusinessError,
			wantCode:   "BUSINESS_ERROR",
			wantType:   ErrTypeBusiness,
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "ErrRateLimit",
			err:        ErrRateLimit,
			wantCode:   "RATE_LIMIT",
			wantType:   ErrTypeBusiness,
			wantStatus: http.StatusTooManyRequests,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.err)

			if err.Code() != tt.wantCode {
				t.Errorf("Code() = %v, want %v", err.Code(), tt.wantCode)
			}

			if err.Type() != tt.wantType {
				t.Errorf("Type() = %v, want %v", err.Type(), tt.wantType)
			}

			if err.HttpStatus() != tt.wantStatus {
				t.Errorf("HttpStatus() = %v, want %v", err.HttpStatus(), tt.wantStatus)
			}
		})
	}
}

// TestNewfWithPredefinedErrors 测试 Newf 与预定义错误码的结合使用
func TestNewfWithPredefinedErrors(t *testing.T) {
	username := "testuser"
	err := Newf(ErrUsernameExisted, "Username '%s' is already taken", username)

	if err.Code() != "USERNAME_EXISTED" {
		t.Errorf("Code() = %v, want %v", err.Code(), "USERNAME_EXISTED")
	}

	if err.Type() != ErrTypeBusiness {
		t.Errorf("Type() = %v, want %v", err.Type(), ErrTypeBusiness)
	}

	if err.HttpStatus() != http.StatusConflict {
		t.Errorf("HttpStatus() = %v, want %v", err.HttpStatus(), http.StatusConflict)
	}

	expectedMessage := fmt.Sprintf("Username '%s' is already taken", username)
	if err.Message() != expectedMessage {
		t.Errorf("Message() = %v, want %v", err.Message(), expectedMessage)
	}
}

// BenchmarkErrorCreation 性能测试
func BenchmarkErrorCreation(b *testing.B) {
	b.Run("New", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = New(ErrInternal)
		}
	})

	b.Run("Newf", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = Newf(ErrInternal, "error with value: %d", i)
		}
	})

	b.Run("Wrap", func(b *testing.B) {
		originalErr := fmt.Errorf("original error")
		for i := 0; i < b.N; i++ {
			_ = Wrap(originalErr, ErrInternal)
		}
	})
}

// TestConcurrentErrorCreation 并发安全性测试
func TestConcurrentErrorCreation(t *testing.T) {
	concurrency := 100
	done := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(index int) {
			err := InternalErrorf("concurrent error %d", index)

			// 验证错误基本属性
			if err.Code() != ErrTypeInternal.String() {
				t.Errorf("Concurrent error Code() = %v, want %v", err.Code(), ErrTypeInternal.String())
			}

			if err.Type() != ErrTypeInternal {
				t.Errorf("Concurrent error Type() = %v, want %v", err.Type(), ErrTypeInternal)
			}

			done <- true
		}(i)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < concurrency; i++ {
		<-done
	}
}

// TestConcurrentErrorCreationNoStack 并发安全性测试，不带堆栈信息
func TestConcurrentErrorCreationNoStack(t *testing.T) {
	concurrency := 100
	done := make(chan bool, concurrency)

	for i := 0; i < concurrency; i++ {
		go func(index int) {
			err := InternalErrorfNoStack("concurrent error %d", index)

			// 验证错误基本属性
			if err.Code() != ErrTypeInternal.String() {
				t.Errorf("Concurrent error Code() = %v, want %v", err.Code(), ErrTypeInternal.String())
			}

			if err.Type() != ErrTypeInternal {
				t.Errorf("Concurrent error Type() = %v, want %v", err.Type(), ErrTypeInternal)
			}

			done <- true
		}(i)
	}

	// 等待所有 goroutine 完成
	for i := 0; i < concurrency; i++ {
		<-done
	}
}

func TestNew(t *testing.T) {
	err := New(&ErrCode{
		Code:       "TEST_CODE",
		Message:    "This is a test error",
		HttpStatus: http.StatusBadRequest,
		Type:       ErrTypeBusiness,
	})
	t.Logf("Error: %+v", err)
	t.Logf("Error Code: %s", err.Code())
	t.Logf("Error Message: %s", err.Message())
	t.Logf("Error HTTP Status: %d", err.HttpStatus())
	t.Logf("Error Type: %s", err.Type())
}

func TestNewf(t *testing.T) {
	err := Newf(&ErrCode{
		Code:       "TEST_CODE",
		Message:    "This is a test error",
		HttpStatus: http.StatusBadRequest,
		Type:       ErrTypeBusiness,
	}, "This is a formatted error: %s", "test")
	t.Logf("Error: %+v", err)
	t.Logf("Error Code: %s", err.Code())
	t.Logf("Error Message: %s", err.Message())
	t.Logf("Error HTTP Status: %d", err.HttpStatus())
	t.Logf("Error Type: %s", err.Type())
}
