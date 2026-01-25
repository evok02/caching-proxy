package proxy

import (
	"github.com/evok02/cacher/internal/storage"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func NewServer(infoLogs *os.File, errorLogs *os.File) *http.Server {
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")

	rdb, err := storage.InitDB()
	if err != nil {
		log.Printf("couldnt connect to db: %s", err.Error())
	}

	infoWriter := io.MultiWriter(os.Stdout, infoLogs)
	errorWriter := io.MultiWriter(os.Stdout, errorLogs)

	cfg := NewApiConfig(rdb, infoWriter, errorWriter)

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
