package proxy

import (
	"github.com/evok02/cacher/storage"
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

	cfg := NewApiConfig(rdb)
	log.Printf("%+v\n", cfg)

	mux := http.NewServeMux()
	mux.HandleFunc("/", cfg.HandleRequest)
	return &http.Server{
		Addr:         host + ":" + port,
		ReadTimeout:  time.Second * 3,
		WriteTimeout: time.Second * 5,
		Handler:      mux,
	}
}
