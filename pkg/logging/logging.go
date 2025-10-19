package logging

import "context"

type Logger interface {
	Debug(msg ...string)
	Info(msg string, args ...interface{})
	Warn(msg string, errs ...error)
	Error(msg string, errs ...error)
	Fatal(msg string, errs ...interface{})

	DebugF(format string, args ...any)
	InfoF(format string, args ...any)
	WarnF(format string, args ...any)
	ErrorF(format string, args ...any)
	FatalF(format string, args ...any)

	WithField(key string, value any) Logger
	WithFields(fields map[string]any) Logger

	DebugCtx(ctx context.Context, msg string)
	InfoCtx(ctx context.Context, msg string)
	WarnCtx(ctx context.Context, msg string, errs ...error)
	ErrorCtx(ctx context.Context, msg string, errs ...error)
	FatalCtx(ctx context.Context, msg string, errs ...error)

	DebugFCtx(ctx context.Context, format string, args ...any)
	InfoFCtx(ctx context.Context, format string, args ...any)
	WarnFCtx(ctx context.Context, format string, args ...any)
	ErrorFCtx(ctx context.Context, format string, args ...any)
	FatalFCtx(ctx context.Context, format string, args ...any)
}
