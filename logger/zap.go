package logger

import (
	"go.uber.org/zap"
)

type ZapLogger struct {
	l *zap.Logger
}

func NewZapLogger(l *zap.Logger) Logger {
	return &ZapLogger{l: l}
}

func (z *ZapLogger) Debug(msg string, fields ...Field) {
	z.l.Debug(msg, z.withField(fields...)...)
}

func (z *ZapLogger) Info(msg string, fields ...Field) {
	z.l.Info(msg, z.withField(fields...)...)
}

func (z *ZapLogger) Warn(msg string, fields ...Field) {
	z.l.Warn(msg, z.withField(fields...)...)
}

func (z *ZapLogger) Error(msg string, fields ...Field) {
	z.l.Error(msg, z.withField(fields...)...)
}

func (z *ZapLogger) Fatal(msg string, fields ...Field) {
	z.l.Fatal(msg, z.withField(fields...)...)
}

func (z *ZapLogger) Panic(msg string, fields ...Field) {
	z.l.Panic(msg, z.withField(fields...)...)
}

func (z *ZapLogger) withField(fields ...Field) []zap.Field {
	zapField := make([]zap.Field, len(fields))
	for i := range fields {
		zapField[i] = zap.Any(fields[i].Key, fields[i].Val)
	}

	return zapField
}
