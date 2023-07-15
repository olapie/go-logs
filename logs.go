package logs

import (
	"io"
	"log"
	"log/slog"
	"os"
	"strings"
)

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
			Level: slog.LevelDebug,
		}))
	} else {
		l = slog.New(slog.NewJSONHandler(w, nil))
	}
	slog.SetDefault(l)
	return closer
}
