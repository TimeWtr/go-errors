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
}

func NewBuilder() *Builder {
	return &Builder{
		metadata: map[string]any{},
	}
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

	err := &ErrorImpl{
		code:       b.code.Code,
		message:    b.message,
		httpStatus: b.code.HttpStatus,
		errType:    b.code.Type,
		timestamp:  time.Now().UTC(),
		stackTrace: captureStackTrace(2),
	}
	
	return err
}
