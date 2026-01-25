package proxy

import (
	"context"
	"io"
	"log"
	"os"
)

func SyncDaemon(ctx context.Context, out <-chan struct{}, files ...*os.File) {
	select {
	case <-out:
		for _, file := range files {
			log.Println("got signal")
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
