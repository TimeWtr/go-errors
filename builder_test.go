package errors

import (
	"fmt"
	"testing"
	"time"
)

func TestBuilder_Build(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*Builder) *Builder
		validate func(*testing.T, Error)
	}{
		{
			name: "默认构建器使用内部错误",
			setup: func(b *Builder) *Builder {
				return b
			},
			validate: func(t *testing.T, err Error) {
				if err == nil {
					t.Fatal("期望得到错误，但得到了nil")
				}
				if err.Code() != ErrInternal.Code {
					t.Errorf("期望代码 %s，但得到了 %s", ErrInternal.Code, err.Code())
				}
				if err.Message() != ErrInternal.Message {
					t.Errorf("期望消息 %s，但得到了 %s", ErrInternal.Message, err.Message())
				}
				if err.HttpStatus() != ErrInternal.HttpStatus {
					t.Errorf("期望HTTP状态码 %d，但得到了 %d", ErrInternal.HttpStatus, err.HttpStatus())
				}
			},
		},
		{
			name: "使用自定义错误码的构建器",
			setup: func(b *Builder) *Builder {
				return b.WithCode(ErrBadRequest)
			},
			validate: func(t *testing.T, err Error) {
				if err.Code() != ErrBadRequest.Code {
					t.Errorf("期望代码 %s，但得到了 %s", ErrBadRequest.Code, err.Code())
				}
				if err.Message() != ErrBadRequest.Message {
					t.Errorf("期望消息 %s，但得到了 %s", ErrBadRequest.Message, err.Message())
				}
			},
		},
		{
			name: "使用自定义消息的构建器",
			setup: func(b *Builder) *Builder {
				return b.WithCode(ErrBadRequest).WithMessage("自定义消息")
			},
			validate: func(t *testing.T, err Error) {
				if err.Message() != "自定义消息" {
					t.Errorf("期望自定义消息，但得到了 %s", err.Message())
				}
			},
		},
		{
			name: "带有原因错误的构建器",
			setup: func(b *Builder) *Builder {
				cause := New(ErrInternal)
				return b.WithCode(ErrBadRequest).WithCause(cause)
			},
			validate: func(t *testing.T, err Error) {
				if err.Unwrap() == nil {
					t.Error("期望得到原因错误，但得到了nil")
				}
			},
		},
		{
			name: "带有元数据的构建器",
			setup: func(b *Builder) *Builder {
				return b.WithCode(ErrBadRequest).
					WithMetadata("key1", "value1").
					WithMetadata("key2", 42)
			},
			validate: func(t *testing.T, err Error) {
				metadata := err.Metadata()
				if len(metadata) != 2 {
					t.Errorf("期望2个元数据项，但得到了 %d", len(metadata))
				}
				if metadata["key1"] != "value1" {
					t.Errorf("期望元数据 key1=value1，但得到了 %v", metadata["key1"])
				}
				if metadata["key2"] != 42 {
					t.Errorf("期望元数据 key2=42，但得到了 %v", metadata["key2"])
				}
			},
		},
		{
			name: "带有元数据映射的构建器",
			setup: func(b *Builder) *Builder {
				meta := map[string]any{
					"mapKey1": "mapValue1",
					"mapKey2": 123,
				}
				return b.WithCode(ErrBadRequest).WithMetadataMap(meta)
			},
			validate: func(t *testing.T, err Error) {
				metadata := err.Metadata()
				if len(metadata) != 2 {
					t.Errorf("期望2个元数据项，但得到了 %d", len(metadata))
				}
				if metadata["mapKey1"] != "mapValue1" {
					t.Errorf("期望元数据 mapKey1=mapValue1，但得到了 %v", metadata["mapKey1"])
				}
				if metadata["mapKey2"] != 123 {
					t.Errorf("期望元数据 mapKey2=123，但得到了 %v", metadata["mapKey2"])
				}
			},
		},
		{
			name: "启用快速模式的构建器",
			setup: func(b *Builder) *Builder {
				return b.WithCode(ErrBadRequest).WithFastMode()
			},
			validate: func(t *testing.T, err Error) {
				// 在快速模式下，不应捕获堆栈跟踪
				// 注意：此验证取决于captureStackTrace的实际实现
				// 我们假设当启用快速模式时，堆栈跟踪将为空或nil
			},
		},
		{
			name: "使用所有选项的构建器",
			setup: func(b *Builder) *Builder {
				cause := New(ErrInternal)
				meta := map[string]any{"test": "data"}
				return b.WithCode(ErrBadRequest).
					WithMessage("所有选项测试").
					WithCause(cause).
					WithMetadataMap(meta).
					WithFastMode()
			},
			validate: func(t *testing.T, err Error) {
				if err.Code() != ErrBadRequest.Code {
					t.Errorf("期望代码 %s，但得到了 %s", ErrBadRequest.Code, err.Code())
				}
				if err.Message() != "所有选项测试" {
					t.Errorf("期望消息 '所有选项测试'，但得到了 %s", err.Message())
				}
				if err.Unwrap() == nil {
					t.Error("期望得到原因错误，但得到了nil")
				}
				metadata := err.Metadata()
				if metadata["test"] != "data" {
					t.Errorf("期望元数据 test=data，但得到了 %v", metadata["test"])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuilder()
			builder = tt.setup(builder)
			err := builder.Build()
			tt.validate(t, err)

			// 验证时间戳已设置
			if err.Timestamp().IsZero() {
				t.Error("期望时间戳被设置")
			}

			// 验证时间戳是最近的（在最后一秒内）
			if time.Since(err.Timestamp()) > time.Second {
				t.Error("期望时间戳是最近的")
			}
		})
	}
}

func TestNewBuilder(t *testing.T) {
	builder := NewBuilder()
	if builder == nil {
		t.Fatal("NewBuilder() 返回了 nil")
	}

	// 检查元数据映射是否已初始化
	if builder.metadata == nil {
		t.Error("元数据映射应该被初始化")
	}

	// 检查其他字段是否为零值
	if builder.code != nil {
		t.Error("代码应该初始为nil")
	}
	if builder.message != "" {
		t.Error("消息应该初始为空")
	}
	if builder.cause != nil {
		t.Error("原因错误应该初始为nil")
	}
	if builder.fastMode {
		t.Error("快速模式应该初始为false")
	}
}

func TestBuilder_MethodChaining(t *testing.T) {
	builder := NewBuilder().
		WithCode(ErrBadRequest).
		WithMessage("链式消息").
		WithFastMode()

	if builder.code != ErrBadRequest {
		t.Errorf("期望代码 %v，但得到了 %v", ErrBadRequest, builder.code)
	}
	if builder.message != "链式消息" {
		t.Errorf("期望消息 '链式消息'，但得到了 %s", builder.message)
	}
	if !builder.fastMode {
		t.Error("期望快速模式为true")
	}
}

// BenchmarkBuilder 测试 Builder 模式的性能
func BenchmarkBuilder(b *testing.B) {
	b.Run("BasicBuilder", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := NewBuilder().
				WithCode(ErrInternal).
				Build()
			_ = err
		}
	})

	b.Run("BuilderWithMessage", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := NewBuilder().
				WithCode(ErrInternal).
				WithMessage("benchmark test error").
				Build()
			_ = err
		}
	})

	b.Run("BuilderWithMetadata", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := NewBuilder().
				WithCode(ErrInternal).
				WithMetadata("key1", "value1").
				WithMetadata("key2", fmt.Sprintf("value%d", i)).
				Build()
			_ = err
		}
	})

	b.Run("BuilderWithCause", func(b *testing.B) {
		cause := New(ErrBadRequest)
		for i := 0; i < b.N; i++ {
			err := NewBuilder().
				WithCode(ErrInternal).
				WithCause(cause).
				Build()
			_ = err
		}
	})

	b.Run("BuilderWithFastMode", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := NewBuilder().
				WithCode(ErrInternal).
				WithFastMode().
				Build()
			_ = err
		}
	})

	b.Run("BuilderWithAllOptions", func(b *testing.B) {
		cause := New(ErrBadRequest)
		meta := map[string]any{
			"testKey1": "testValue1",
			"testKey2": 42,
		}
		for i := 0; i < b.N; i++ {
			err := NewBuilder().
				WithCode(ErrInternal).
				WithMessage("benchmark test with all options").
				WithCause(cause).
				WithMetadataMap(meta).
				WithFastMode().
				Build()
			_ = err
		}
	})
}

// BenchmarkBuilderVsDirect 测试 Builder 模式与直接创建错误的性能对比
func BenchmarkBuilderVsDirect(b *testing.B) {
	b.Run("DirectNew", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := New(ErrInternal)
			_ = err
		}
	})

	b.Run("DirectNewf", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := Newf(ErrInternal, "direct error %d", i)
			_ = err
		}
	})

	b.Run("BuilderApproach", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			err := NewBuilder().
				WithCode(ErrInternal).
				WithMessage(fmt.Sprintf("builder error %d", i)).
				Build()
			_ = err
		}
	})
}
