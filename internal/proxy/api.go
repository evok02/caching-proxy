package proxy

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/evok02/cacher/internal/storage"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type ApiConfig struct {
	Storage     storage.RedisStorage
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	SyncChan    chan<- struct{}
}

func NewApiConfig(rdb storage.RedisStorage, sync chan<- struct{}, infoOut, errorOut io.Writer) *ApiConfig {
	return &ApiConfig{
		Storage:     rdb,
		InfoLogger:  NewInfoLogger(infoOut),
		ErrorLogger: NewErrorLogger(errorOut),
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, v error) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResposne{
		Error: v.Error(),
	})
}

func hashRequest(req RequestCache) (string, error) {
	b, err := json.Marshal(req)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", sha256.Sum256(b)), nil
}

func (cfg *ApiConfig) HandleRequest(w http.ResponseWriter, r *http.Request) {
	target := os.Getenv("TARGET")
	cfg.InfoLogger.Printf("Incoming request method: %s\n", r.Method)
	cfg.InfoLogger.Printf("Incoming host: %s\n", r.Host)
	cfg.InfoLogger.Printf("Incoming addr: %s\n", r.RemoteAddr)
	cfg.InfoLogger.Printf("Incoming url: %s\n", r.URL.String())

	targetUrl, err := url.JoinPath(target, r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		cfg.SyncChan <- struct{}{}
		return
	}

	proxyReq, err := http.NewRequest(r.Method, targetUrl, r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		cfg.ErrorLogger.Printf("NewRequest: %s\n", err.Error())
		cfg.SyncChan <- struct{}{}
		return
	}

	proxyReq.Header = r.Header.Clone()

	hashReq, err := hashRequest(RequestCache{
		Method: proxyReq.Method,
		URL:    proxyReq.URL.String(),
		Header: proxyReq.Header,
	})

	dbRes, err := cfg.Storage.Get(context.TODO(), hashReq).Result()

	if err != nil {
		cfg.ErrorLogger.Printf("couldnt hit the cache: %s\n", err.Error())
	} else {
		cfg.InfoLogger.Println("hit the cache~")
		writeJSON(w, http.StatusOK, dbRes)
		return
	}

	cfg.InfoLogger.Printf("Making the request %s %s", proxyReq.Method, proxyReq.URL.String())
	res, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		writeError(w, 500, err)
		cfg.ErrorLogger.Printf("DefaultClient.Do: %s\n", err.Error())
		cfg.SyncChan <- struct{}{}
		return
	}
	defer res.Body.Close()

	cfg.InfoLogger.Printf("Got response %d %s\n", res.StatusCode, res.Status)

	resStruct := ResponseCache{
		Status:        res.Status,
		StatusCode:    res.StatusCode,
		Header:        res.Header,
		ContentLength: res.ContentLength,
	}

	writeJSON(w, http.StatusOK, resStruct)
	resBytes, err := json.Marshal(resStruct)
	if err != nil {
		writeError(w, 500, err)
	}

	_, err = cfg.Storage.Set(context.TODO(), hashReq, resBytes, 6*time.Hour).Result()

	if err != nil {
		cfg.ErrorLogger.Printf("unable to store response: %s\n", err.Error())
	} else {
		cfg.InfoLogger.Printf("set new value to the cache\n")
	}
}
