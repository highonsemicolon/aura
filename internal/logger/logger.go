package logger

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string, errs ...error)
	Error(msg string, errs ...error)
	Fatal(msg string, errs ...error)

	DebugF(format string, args ...any)
	InfoF(format string, args ...any)
	WarnF(format string, args ...any)
	ErrorF(format string, args ...any)
	FatalF(format string, args ...any)

	WithField(key string, value any) Logger
	WithFields(fields map[string]any) Logger
}
