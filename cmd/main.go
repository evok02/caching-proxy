package main

import (
	"context"
	"github.com/evok02/cacher/internal/config"
	"github.com/evok02/cacher/internal/proxy"
	"log"
	"os"
	"time"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("%+v\n", cfg)

	infoLogs, err := os.OpenFile(cfg.Proxy.Logs.InfoLogsPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("couldnt open an infologs file: %s", err.Error())
	}
	defer infoLogs.Close()

	errorLogs, err := os.OpenFile(cfg.Proxy.Logs.ErrLogsPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Printf("couldnt open and errorLogs file: %s", err.Error())
	}
	defer errorLogs.Close()
	s := proxy.NewServer(infoLogs, errorLogs, cfg.Proxy)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ticker := time.NewTicker(time.Second)
	go proxy.SyncDaemon(ctx, ticker.C, infoLogs, errorLogs)

	log.Println("Listening on new connections...")
	s.ListenAndServe()
}
