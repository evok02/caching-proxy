package proxy

import (
	"context"
	"github.com/evok02/cacher/storage"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func NewServer() *http.Server {
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")

	rdb, err := storage.InitDB()
	if err != nil {
		log.Printf("couldnt connect to db: %s", err.Error())
	}

	infoLogs, err := os.OpenFile("./log/info.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("couldnt open an infologs file: %s", err.Error())
	}
	defer infoLogs.Close()

	errorLogs, err := os.OpenFile("./log/error.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("couldnt open and errorLogs file: %s", err.Error())
	}
	defer errorLogs.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	syncChan := make(chan struct{})
	go SyncDaemon(ctx, syncChan, infoLogs, errorLogs)

	infoWriter := io.MultiWriter(os.Stdout, infoLogs)
	errorWriter := io.MultiWriter(os.Stdout, errorLogs)

	cfg := NewApiConfig(rdb, syncChan, infoWriter, errorWriter)

	mux := http.NewServeMux()
	mux.HandleFunc("/", cfg.HandleRequest)
	return &http.Server{
		Addr:         host + ":" + port,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 5,
		Handler:      mux,
		ErrorLog:     cfg.ErrorLogger,
	}
}
