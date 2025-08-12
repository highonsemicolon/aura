package logger

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
	otelCodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type zerologAdapter struct {
	logger *zerolog.Logger
}

func NewZerologAdapter(format, level string) Logger {
	writer := os.Stdout
	format = strings.ToLower(format)
	level = strings.ToLower(level)

	var logWriter io.Writer
	if format == "json" {
		logWriter = writer
	} else {
		logWriter = zerolog.ConsoleWriter{Out: writer}
	}

	parsedLevel, err := zerolog.ParseLevel(level)
	if err != nil {
		parsedLevel = zerolog.InfoLevel
	}

	logger := zerolog.New(logWriter).
		Level(parsedLevel).
		With().
		Timestamp().
		Logger()

	return &zerologAdapter{
		logger: &logger,
	}
}

func UnaryServerZerologInterceptor(base Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		log := base.WithField("method", info.FullMethod)

		ctx = contextWithLogger(ctx, log)
		span := trace.SpanFromContext(ctx)

		resp, err = handler(ctx, req)

		if err != nil {
			span.RecordError(err)
			span.SetStatus(otelCodes.Error, err.Error())
			st := status.Convert(err)
			log.WithField("status", st.Code()).
				WithField("error", st.Message()).
				ErrorCtx(ctx, "request failed")
		} else {
			log.WithField("status", grpcCodes.OK).
				InfoCtx(ctx, "request completed")
		}

		return resp, err
	}
}

type loggerContextKey struct{}

var defaultLogger Logger = NewZerologAdapter("console", "info")

func contextWithLogger(ctx context.Context, log Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey{}, log)
}

func FromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value(loggerContextKey{}).(Logger); ok {
		return wrapWithTrace(ctx, l)
	}
	return wrapWithTrace(ctx, defaultLogger)
}

func (z *zerologAdapter) Debug(msg ...string) {
	z.logger.Debug().Msg(strings.Join(msg, " "))
}

func (z *zerologAdapter) Info(msg ...string) {
	z.logger.Info().Msg(strings.Join(msg, " "))
}

func (z *zerologAdapter) Warn(msg string, errs ...error) {
	events := z.logger.Warn()
	if len(errs) > 0 {
		events = events.Errs("errors", errs)
	}
	events.Msg(msg)
}

func (z *zerologAdapter) Error(msg string, errs ...error) {
	events := z.logger.Error()
	if len(errs) > 0 {
		events = events.Errs("errors", errs)
	}
	events.Msg(msg)
}

func (z *zerologAdapter) Fatal(msg string, errs ...error) {
	event := z.logger.Fatal()

	if len(errs) > 0 {
		event = event.Errs("errors", errs)
	}

	event.Msg(msg)
}

func (z *zerologAdapter) DebugF(format string, args ...any) {
	z.logger.Debug().Msgf(format, args...)
}
func (z *zerologAdapter) InfoF(format string, args ...any) {
	z.logger.Info().Msgf(format, args...)
}
func (z *zerologAdapter) WarnF(format string, args ...any) {
	z.logger.Warn().Msgf(format, args...)
}
func (z *zerologAdapter) ErrorF(format string, args ...any) {
	z.logger.Error().Msgf(format, args...)
}
func (z *zerologAdapter) FatalF(format string, args ...any) {
	z.logger.Fatal().Msgf(format, args...)
}

func (z *zerologAdapter) WithField(key string, value any) Logger {
	newLogger := z.logger.With().Interface(key, value).Logger()
	return &zerologAdapter{logger: &newLogger}
}

func (z *zerologAdapter) WithFields(fields map[string]any) Logger {
	ctx := z.logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	newLogger := ctx.Logger()
	return &zerologAdapter{logger: &newLogger}
}

func (z *zerologAdapter) withTrace(ctx context.Context) *zerolog.Logger {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return z.logger
	}
	l := z.logger.With().
		Str("trace_id", span.SpanContext().TraceID().String()).
		Str("span_id", span.SpanContext().SpanID().String()).
		Logger()
	return &l
}

func wrapWithTrace(ctx context.Context, l Logger) Logger {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return l
	}
	return l.WithFields(map[string]any{
		"trace_id": span.SpanContext().TraceID().String(),
		"span_id":  span.SpanContext().SpanID().String(),
	})
}
func (z *zerologAdapter) DebugCtx(ctx context.Context, msg string) {
	z.withTrace(ctx).Debug().Msg(msg)
}

func (z *zerologAdapter) InfoCtx(ctx context.Context, msg string) {
	z.withTrace(ctx).Info().Msg(msg)
}

func (z *zerologAdapter) WarnCtx(ctx context.Context, msg string, errs ...error) {
	ev := z.withTrace(ctx).Warn()
	if len(errs) > 0 {
		ev = ev.Errs("errors", errs)
	}
	ev.Msg(msg)
}
func (z *zerologAdapter) ErrorCtx(ctx context.Context, msg string, errs ...error) {
	ev := z.withTrace(ctx).Error()
	if len(errs) > 0 {
		ev = ev.Errs("errors", errs)
	}
	ev.Msg(msg)
}

func (z *zerologAdapter) FatalCtx(ctx context.Context, msg string, errs ...error) {
	ev := z.withTrace(ctx).Fatal()
	if len(errs) > 0 {
		ev = ev.Errs("errors", errs)
	}
	ev.Msg(msg)
	os.Exit(1)
}

func (z *zerologAdapter) DebugFCtx(ctx context.Context, format string, args ...any) {
	z.withTrace(ctx).Debug().Msgf(format, args...)
}
func (z *zerologAdapter) InfoFCtx(ctx context.Context, format string, args ...any) {
	z.withTrace(ctx).Info().Msgf(format, args...)
}

func (z *zerologAdapter) WarnFCtx(ctx context.Context, format string, args ...any) {
	ev := z.withTrace(ctx).Warn()
	if len(args) > 0 {
		ev.Msgf(format, args...)
	} else {
		ev.Msg(format)
	}
	ev.Send()
}
func (z *zerologAdapter) ErrorFCtx(ctx context.Context, format string, args ...any) {
	z.withTrace(ctx).Error().Msgf(format, args...)
}

func (z *zerologAdapter) FatalFCtx(ctx context.Context, format string, args ...any) {
	ev := z.withTrace(ctx).Fatal()
	if len(args) > 0 {
		ev.Msgf(format, args...)
	} else {
		ev.Msg(format)
	}
	ev.Send()
	os.Exit(1)
}
