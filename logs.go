package logs

import (
	"io"
	"log"
	"log/slog"
	"os"
	"strings"
)

const (
	EnvDebug = "OLA_LOG_DEBUG" // 1, true, or enabled
)

// Init default logger writing to stderr and filename if it's not empty
func Init(filename string, opts ...func(options *slog.HandlerOptions)) io.Closer {
	options := new(slog.HandlerOptions)
	for _, opt := range opts {
		opt(options)
	}

	if isDebugging() {
		options.Level = slog.LevelDebug
	}

	if filename == "" {
		slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, options)))
		return nil
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("open file %s: %v\n", filename, err)
	}

	h := MultiHandler(slog.NewJSONHandler(f, options),
		slog.NewTextHandler(os.Stderr, options))
	slog.SetDefault(slog.New(h))
	return f
}

type RotationOptions struct {
	HandlerOptions slog.HandlerOptions
	RotateFileWriterOptions
}

// InitRotate init default logger writing to rotated file and stderr
// If filename is empty, then a default filename is used
func InitRotate(filename string, opts ...func(options *RotationOptions)) io.Closer {
	if filename == "" {
		log.Fatalln("filename is empty")
	}

	options := &RotationOptions{
		RotateFileWriterOptions: *defaultRotateFileWriterOptions(),
	}

	for _, opt := range opts {
		opt(options)
	}

	if isDebugging() {
		options.HandlerOptions.Level = slog.LevelDebug
	}

	f := NewRotateFileWriter(filename, func(o *RotateFileWriterOptions) {
		*o = options.RotateFileWriterOptions
	})
	h := MultiHandler(slog.NewJSONHandler(f, &options.HandlerOptions),
		slog.NewTextHandler(os.Stderr, &options.HandlerOptions))
	slog.SetDefault(slog.New(h))
	return f
}

func isDebugging() bool {
	debug := strings.ToLower(os.Getenv(EnvDebug))
	return debug == "1" || debug == "true" || debug == "enabled"
}
