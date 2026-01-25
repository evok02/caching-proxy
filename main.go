package main

import (
	"context"
	"github.com/evok02/cacher/internal/proxy"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

func main() {
	godotenv.Load()
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
	s := proxy.NewServer(infoLogs, errorLogs)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := time.NewTicker(time.Second)
	go proxy.SyncDaemon(ctx, ticker.C, infoLogs, errorLogs)

	log.Println("Listening on new connections...")
	s.ListenAndServe()
}
