package proxy

import (
	"github.com/evok02/cacher/internal/config"
	"github.com/evok02/cacher/internal/storage"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func NewServer(infoLogs *os.File, errorLogs *os.File, cfg config.ProxyConfig) *http.Server {
	port := cfg.Port
	host := cfg.Host

	rdb, err := storage.InitDB(cfg.DB)
	if err != nil {
		log.Printf("couldnt connect to db: %s", err.Error())
	}

	infoWriter := io.MultiWriter(os.Stdout, infoLogs)
	errorWriter := io.MultiWriter(os.Stdout, errorLogs)

	apiCfg := NewApiConfig(rdb, infoWriter, errorWriter)
	apiCfg.SetTarget(cfg.Target)

	mux := http.NewServeMux()
	mux.HandleFunc("/", apiCfg.HandleRequest)
	return &http.Server{
		Addr:         host + ":" + port,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 5,
		Handler:      mux,
		ErrorLog:     apiCfg.ErrorLogger,
	}
}
