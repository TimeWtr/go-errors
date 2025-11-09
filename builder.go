package errors

import "time"

// Builder 错误构造器
type Builder struct {
	// 错误码
	code *ErrCode
	// 错误详情
	message string
	// 原始错误
	cause error
	// 元数据信息
	metadata map[string]any
	// 是否启用快速模式，快速模式不会捕获堆栈信息
	fastMode bool
}

func NewBuilder() *Builder {
	return &Builder{
		metadata: map[string]any{},
	}
}

func (b *Builder) WithFastMode() *Builder {
	b.fastMode = true
	return b
}

func (b *Builder) WithCode(code *ErrCode) *Builder {
	b.code = code
	return b
}

func (b *Builder) WithMessage(message string) *Builder {
	b.message = message
	return b
}

func (b *Builder) WithCause(cause error) *Builder {
	b.cause = cause
	return b
}

func (b *Builder) WithMetadata(key string, val any) *Builder {
	b.metadata[key] = val
	return b
}

func (b *Builder) WithMetadataMap(metadata map[string]any) *Builder {
	for k, v := range metadata {
		b.metadata[k] = v
	}
	return b
}

func (b *Builder) Build() Error {
	if b.code == nil {
		b.code = ErrInternal
	}

	if b.message == "" {
		b.message = b.code.Message
	}

	impl := acquireError()
	impl.code = b.code.Code
	impl.message = b.message
	impl.httpStatus = b.code.HttpStatus
	impl.errType = b.code.Type
	impl.timestamp = time.Now().UTC()
	impl.cause = b.cause
	impl.metadata = b.metadata

	// 不开启快速模式，则记录堆栈信息
	if !b.fastMode {
		impl.stackTrace = captureStackTrace(2)
	}

	return impl
}
