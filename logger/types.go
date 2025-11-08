package logger

import "time"

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Panic(msg string, fields ...Field)
}

type Field struct {
	Key string
	Val any
}

func StringField(key string, val string) Field {
	return Field{Key: key, Val: val}
}

func IntField(key string, val int) Field {
	return Field{Key: key, Val: val}
}

func Int64Field(key string, val int64) Field {
	return Field{Key: key, Val: val}
}

func BoolField(key string, val bool) Field {
	return Field{Key: key, Val: val}
}

func DurationField(key string, val time.Duration) Field {
	return Field{Key: key, Val: val}
}

func TimeField(key string, val time.Time) Field {
	return Field{Key: key, Val: val}
}

func ErrorField(err error) Field {
	return Field{Key: "error", Val: err}
}
