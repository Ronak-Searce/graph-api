package logger

import (
	"context"
	"io"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type levelName string

// log level names
const (
	NameDebug levelName = "DEBUG"
	NameInfo  levelName = "INFO"
	NameWarn  levelName = "WARN"
	NameError levelName = "ERROR"
)

var (
	global          *zap.SugaredLogger
	defaultLevel    = zap.NewAtomicLevelAt(zap.InfoLevel)
	messageCounters = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "logger_messages_total",
		Help: "Count of logged messages",
	}, []string{"level"})
)

func init() {
	SetLogger(New(defaultLevel))
	// metrics.MustRegister(messageCounters)
}

// InitWithName ...
func InitWithName(lvl levelName) bool {
	l := map[levelName]zapcore.Level{
		NameDebug: zapcore.DebugLevel,
		NameInfo:  zapcore.InfoLevel,
		NameWarn:  zapcore.WarnLevel,
		NameError: zapcore.ErrorLevel,
	}
	if v, ok := l[lvl]; ok {
		SetLevel(v)
		return true
	}
	SetLevel(defaultLevel.Level())
	return false
}

// InitWithString ...
func InitWithString(lvl string) bool {
	n := levelName(lvl)
	return InitWithName(n)
}

// New ...
func New(level zapcore.LevelEnabler, options ...zap.Option) *zap.SugaredLogger {
	return NewWithSink(level, os.Stdout, options...)
}

// NewWithSink ...
func NewWithSink(level zapcore.LevelEnabler, sink io.Writer, options ...zap.Option) *zap.SugaredLogger {
	var stackTraceLvls zap.LevelEnablerFunc = func(l zapcore.Level) bool {
		return l >= zap.ErrorLevel
	}
	options = append([]zap.Option{zap.AddStacktrace(stackTraceLvls)}, options...)
	if level == nil {
		level = defaultLevel
	}
	return zap.New(
		zapcore.NewCore(
			zapcore.NewJSONEncoder(zapcore.EncoderConfig{
				TimeKey:        "time", // https://cloud.google.com/logging/docs/agent/logging/configuration#timestamp-processing
				LevelKey:       "level",
				NameKey:        "logger",
				CallerKey:      "caller",
				MessageKey:     "message",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			}),
			zapcore.AddSync(sink),
			level,
		),
		options...,
	).Sugar()
}

// Logger ...
func Logger() *zap.SugaredLogger {
	return global
}

// SetLogger ...
func SetLogger(l *zap.SugaredLogger) {
	global = l
}

// Level ...
func Level() zapcore.Level {
	return defaultLevel.Level()
}

// SetLevel ...
func SetLevel(l zapcore.Level) {
	defaultLevel.SetLevel(l)
}

// Debug ...
func Debug(ctx context.Context, args ...interface{}) {
	FromContext(ctx).Debug(args...)
	messageCounters.WithLabelValues("debug").Inc()
}

// Debugf ...
func Debugf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Debugf(format, args...)
	messageCounters.WithLabelValues("debug").Inc()
}

// DebugKV ...
func DebugKV(ctx context.Context, message string, kvs ...interface{}) {
	FromContext(ctx).Debugw(message, kvs...)
	messageCounters.WithLabelValues("debug").Inc()
}

// Info ...
func Info(ctx context.Context, args ...interface{}) {
	FromContext(ctx).Info(args...)
	messageCounters.WithLabelValues("info").Inc()
}

// Infof ...
func Infof(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Infof(format, args...)
	messageCounters.WithLabelValues("info").Inc()
}

// InfoKV ...
func InfoKV(ctx context.Context, message string, kvs ...interface{}) {
	FromContext(ctx).Infow(message, kvs...)
	messageCounters.WithLabelValues("info").Inc()
}

// Warn ...
func Warn(ctx context.Context, args ...interface{}) {
	FromContext(ctx).Warn(args...)
	messageCounters.WithLabelValues("warn").Inc()
}

// Warnf ...
func Warnf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Warnf(format, args...)
	messageCounters.WithLabelValues("warn").Inc()
}

// WarnKV ...
func WarnKV(ctx context.Context, message string, kvs ...interface{}) {
	FromContext(ctx).Warnw(message, kvs...)
	messageCounters.WithLabelValues("warn").Inc()
}

// Error ...
func Error(ctx context.Context, args ...interface{}) {
	FromContext(ctx).Error(args...)
	messageCounters.WithLabelValues("error").Inc()
}

// Errorf ...
func Errorf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Errorf(format, args...)
	messageCounters.WithLabelValues("error").Inc()
}

// ErrorKV ...
func ErrorKV(ctx context.Context, message string, kvs ...interface{}) {
	FromContext(ctx).Errorw(message, kvs...)
	messageCounters.WithLabelValues("error").Inc()
}

// Fatal ...
func Fatal(ctx context.Context, args ...interface{}) {
	FromContext(ctx).Fatal(args...)
	messageCounters.WithLabelValues("fatal").Inc()
}

// Fatalf ...
func Fatalf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Fatalf(format, args...)
	messageCounters.WithLabelValues("fatal").Inc()
}

// FatalKV ...
func FatalKV(ctx context.Context, message string, kvs ...interface{}) {
	FromContext(ctx).Fatalw(message, kvs...)
	messageCounters.WithLabelValues("fatal").Inc()
}

// Panic ...
func Panic(ctx context.Context, args ...interface{}) {
	FromContext(ctx).Panic(args...)
	messageCounters.WithLabelValues("panic").Inc()
}

// Panicf ...
func Panicf(ctx context.Context, format string, args ...interface{}) {
	FromContext(ctx).Panicf(format, args...)
	messageCounters.WithLabelValues("panic").Inc()
}

// PanicKV ...
func PanicKV(ctx context.Context, message string, kvs ...interface{}) {
	FromContext(ctx).Panicw(message, kvs...)
	messageCounters.WithLabelValues("panic").Inc()
}
