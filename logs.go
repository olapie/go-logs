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
	var w io.Writer
	var closer io.Closer
	if filename != "" {
		f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Fatalln("cannot open file", filename, err)
		}
		w = io.MultiWriter(f, os.Stderr)
		closer = f
	} else {
		w = os.Stderr
	}

	var l *slog.Logger
	debug := strings.ToLower(os.Getenv("DEBUG"))
	if debug == "1" || debug == "true" || debug == "enabled" {
		l = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}))
	} else {
		l = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		}))
	}
	slog.SetDefault(l)
	return closer
}

// InitRotate init default logger writing to rotated file and stderr
// If filename is empty, then a default filename is used
func InitRotate(filename string, optFns ...func(options *RotateFileWriterOptions)) io.Closer {
	f := NewRotateFileWriter(filename, optFns...)
	w := io.MultiWriter(f, os.Stderr)

	var l *slog.Logger
	debug := strings.ToLower(os.Getenv("DEBUG"))
	if debug == "1" || debug == "true" || debug == "enabled" {
		l = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}))
	} else {
		l = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelInfo,
		}))
	}
	slog.SetDefault(l)
	return f
}
