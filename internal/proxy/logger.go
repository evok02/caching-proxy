package proxy

import (
	"context"
	"io"
	"log"
	"os"
	"time"
)

func SyncDaemon(ctx context.Context, ticks <-chan time.Time, files ...*os.File) {
	select {
	case <-ticks:
		for _, file := range files {
			file.Sync()
		}
	case <-ctx.Done():
		return
	}
}

func NewErrorLogger(out io.Writer) *log.Logger {
	return log.New(out, "ERROR: ", log.LstdFlags)
}

func NewInfoLogger(out io.Writer) *log.Logger {
	return log.New(out, "INFO: ", log.LstdFlags)
}
