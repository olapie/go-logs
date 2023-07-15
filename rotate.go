package logs

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/natefinch/lumberjack/v3"
)

var _ io.WriteCloser = (*RotateFileWriter)(nil)

type RotateFileWriterOptions struct {
	// MaxAge is the maximum time to retain old log files based on the timestamp
	// encoded in their filename. The default is not to remove old log files
	// based on age.
	MaxAge time.Duration

	// MaxBackups is the maximum number of old log files to retain. The default
	// is to retain all old log files (though MaxAge may still cause them to get
	// deleted.)
	MaxBackups int

	// LocalTime determines if the time used for formatting the timestamps in
	// backup files is the computer's local time. The default is to use UTC
	// time.
	LocalTime bool

	// Compress determines if the rotated log files should be compressed
	// using gzip. The default is not to perform compression.
	Compress bool

	// MaxSize is the maximum bytes of a log file before being rotated. The default value is 512M
	MaxSize int64
}

type RotateFileWriter struct {
	ll *lumberjack.Roller
}

func NewRotateFileWriter(filename string, optFns ...func(options *RotateFileWriterOptions)) *RotateFileWriter {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		filename = filepath.Join(os.Args[0], "log")
	}

	opts := &RotateFileWriterOptions{
		MaxBackups: 32,
		MaxAge:     30 * time.Hour * 24, // 28 days
		LocalTime:  false,
		Compress:   true,
		MaxSize:    512 * 1024 * 1024,
	}

	for _, fn := range optFns {
		fn(opts)
	}

	ll, err := lumberjack.NewRoller(filename, opts.MaxSize, &lumberjack.Options{
		MaxAge:     opts.MaxAge,
		MaxBackups: opts.MaxBackups,
		LocalTime:  opts.LocalTime,
		Compress:   opts.Compress,
	})

	if err != nil {
		log.Fatalln(err)
	}

	return &RotateFileWriter{
		ll: ll,
	}
}

func (f *RotateFileWriter) Close() error {
	return f.ll.Close()
}

func (f *RotateFileWriter) Write(p []byte) (n int, err error) {
	return f.ll.Write(p)
}
