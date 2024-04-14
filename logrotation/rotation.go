package logrotation

import (
	"io"
	"log"
	"log/slog"
	"os"
	"strings"

	"go.olapie.com/logs"
)

type Options struct {
	HandlerOptions slog.HandlerOptions
	FileWriterOptions
}

// Init default logger writing to d file and stderr
// If filename is empty, then a default filename is used
func Init(filename string, opts ...func(options *Options)) io.Closer {
	if filename == "" {
		log.Fatalln("filename is empty")
	}

	options := &Options{
		FileWriterOptions: *defaultFileWriterOptions(),
	}

	for _, opt := range opts {
		opt(options)
	}

	if isDebugging() {
		options.HandlerOptions.Level = slog.LevelDebug
	}

	f := NewFileWriter(filename, func(o *FileWriterOptions) {
		*o = options.FileWriterOptions
	})
	h := logs.MultiHandler(slog.NewJSONHandler(f, &options.HandlerOptions),
		slog.NewTextHandler(os.Stderr, &options.HandlerOptions))
	slog.SetDefault(slog.New(h))
	return f
}

func isDebugging() bool {
	debug := strings.ToLower(os.Getenv(EnvDebug))
	return debug == "1" || debug == "true" || debug == "enabled"
}

const (
	EnvDebug = "OLA_LOG_DEBUG" // 1, true, or enabled
)
