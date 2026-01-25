package main

import (
	"github.com/evok02/cacher/internal/proxy"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	godotenv.Load()
	s := proxy.NewServer()

	log.Println("Listening on new connections...")
	s.ListenAndServe()
}
