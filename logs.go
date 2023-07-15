package logs

import (
	"io"
	"log"
	"log/slog"
	"os"
	"strings"
)

// Init default logger writing to stderr and filename if it's not empty
func Init(filename string) io.Closer {
	debug := strings.ToLower(os.Getenv("DEBUG"))
	var options *slog.HandlerOptions
	if debug == "1" || debug == "true" || debug == "enabled" {
		options = &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
	}

	if filename == "" {
		if options != nil {
			slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, options)))
		}
		return nil
	}

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln("cannot open file", filename, err)
	}
	h := MultiHandler(slog.NewJSONHandler(f, options),
		slog.NewTextHandler(os.Stderr, options))
	slog.SetDefault(slog.New(h))
	return f
}

// InitRotate init default logger writing to rotated file and stderr
// If filename is empty, then a default filename is used
func InitRotate(filename string, optFns ...func(options *RotateFileWriterOptions)) io.Closer {
	if filename == "" {
		log.Fatalln("filename is empty")
	}

	debug := strings.ToLower(os.Getenv("DEBUG"))
	var options *slog.HandlerOptions
	if debug == "1" || debug == "true" || debug == "enabled" {
		options = &slog.HandlerOptions{
			Level: slog.LevelDebug,
		}
	}

	f := NewRotateFileWriter(filename, optFns...)
	h := MultiHandler(slog.NewJSONHandler(f, options),
		slog.NewTextHandler(os.Stderr, options))
	slog.SetDefault(slog.New(h))
	return f
}
