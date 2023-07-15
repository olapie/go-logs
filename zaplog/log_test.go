package log

import (
	"testing"
	"time"
)

func TestWriteToFile(t *testing.T) {
	ReplaceGlobal(NewLogger(func(options *Options) {
		options.Development = true
		options.ConsoleTimeHidden = true
		options.Filename = "testdata/test.log"
		options.MaxFileSize = 1
		options.Console = false
	}))
	for i := 0; i < 1024*1024/16; i++ {
		Infoln(i, time.Now().String())
	}
}
