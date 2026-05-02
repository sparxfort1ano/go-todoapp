// Package logger contains all needed tools to correctly initialize,
// set and close the Zap logger.
package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger wraps the zap.Logger to write to stored file path.
type Logger struct {
	*zap.Logger

	file *os.File
}

type contextKey string

const logKey contextKey = "log"

// IntoContext places the logger into the context and
// returns the new context.
func IntoContext(ctx context.Context, log *Logger) context.Context {
	return context.WithValue(ctx, logKey, log)
}

// FromContext checks if there is a logger in the context tree.
// It panics if there is no logger.
// Panic is allowed: the application is blind
// without a properly configured logger in the context.
func FromContext(ctx context.Context) *Logger {
	log, ok := ctx.Value(logKey).(*Logger)
	if !ok {
		panic("no logger in context")
	}

	return log
}

// NewLogger initializes and returns a new Zap logger.
// It writes structured logs to stdout and a specified file.
func NewLogger(cfg config) (*Logger, error) {
	zapLvl := zap.NewAtomicLevel()
	if err := zapLvl.UnmarshalText([]byte(cfg.Level)); err != nil {
		return nil, fmt.Errorf("unmarshal log level: %w", err)
	}

	if err := os.MkdirAll(cfg.Folder, 0755); err != nil {
		return nil, fmt.Errorf("mkdir log folder: %w", err)
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.000000")
	logFilePath := filepath.Join(
		cfg.Folder,
		fmt.Sprintf("%s.log", timestamp),
	)

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.000000")

	zapEncoder := zapcore.NewConsoleEncoder(zapConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLvl),
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLvl),
	)

	zapLogger := zap.New(core, zap.AddCaller())

	return &Logger{
		Logger: zapLogger,
		file:   logFile,
	}, nil
}

// With returns a new Logger that includes the provided stuctured fields
// (e.g., request_id, url) in all subsequent log entries.
func (l *Logger) With(field ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(field...),
		file:   l.file,
	}
}

// Close closes the file used by the logger.
func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Println("failed to close application logger:", err)
	}
}
