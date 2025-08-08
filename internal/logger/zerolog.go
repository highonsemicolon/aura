package logger

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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

func UnaryServerZerologInterceptor(log Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		span := trace.SpanFromContext(ctx)
		sc := span.SpanContext()

		logger := log.WithFields(map[string]any{
			"method":   info.FullMethod,
			"trace_id": sc.TraceID().String(),
			"span_id":  sc.SpanID().String(),
		})
		logger.Info("gRPC request started")

		resp, err = handler(ctx, req)

		if err != nil {
			st := status.Convert(err)
			if st.Code() == codes.OK {
				log.InfoF("gRPC handler finished | method=%s", info.FullMethod)
				logger.WithFields(map[string]any{
					"status": st.Code(),
				}).Info("gRPC handler finished")
			} else {
				logger.WithFields(map[string]any{
					"status": st.Code(),
					"error":  st.Message(),
				}).Error("gRPC handler finished with error")
			}
		} else {
			logger.WithFields(map[string]any{
				"status": codes.OK,
			}).Info("gRPC handler finished successfully")
		}

		return resp, err
	}
}

func (z *zerologAdapter) Debug(msg string) {
	z.logger.Debug().Msg(msg)
}

func (z *zerologAdapter) Info(msg string) {
	z.logger.Info().Msg(msg)
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
