package proxy

import (
	"encoding/json"
	"github.com/evok02/cacher/storage"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

type ApiConfig struct {
	Storage     storage.RedisStorage
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	SyncChan    chan<- struct{}
}

type ErrorResposne struct {
	Error string `json:"error"`
}

type ResponseStatistics struct {
	Method string `json:"method"`
	Host   string `json:"host"`
	Addr   string `json:"addr"`
	Url    string `json:"url"`
	Body   []byte `json:"body"`
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

func (cfg *ApiConfig) HandleRequest(w http.ResponseWriter, r *http.Request) {
	target := os.Getenv("TARGET")
	cfg.InfoLogger.Printf("Incoming request method: %s\n", r.Method)
	cfg.InfoLogger.Printf("Incoming host: %s\n", r.Host)
	cfg.InfoLogger.Printf("Incoming addr: %s\n", r.RemoteAddr)
	cfg.InfoLogger.Printf("Incoming url: %s\n\n\n", r.URL.String())

	targetUrl, err := url.JoinPath(target, r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	proxyReq, err := http.NewRequest(r.Method, targetUrl, r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		cfg.ErrorLogger.Printf("NewRequest: %s", err.Error())
		return
	}

	proxyReq.Header = r.Header.Clone()

	cfg.InfoLogger.Printf("Making the request: %+v\n", proxyReq)
	res, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		writeError(w, 500, err)
		cfg.ErrorLogger.Printf("DefaultClient.Do: %s", err.Error())
		return
	}
	defer res.Body.Close()

	cfg.InfoLogger.Printf("Got response %d %s\n", res.StatusCode, res.Status)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		writeError(w, 500, err)
		cfg.ErrorLogger.Printf("io.ReadAll(res.Body): %s", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, ResponseStatistics{
		Method: proxyReq.Method,
		Host:   proxyReq.Host,
		Addr:   proxyReq.RemoteAddr,
		Url:    proxyReq.URL.String(),
		Body:   body,
	})

	// TODO: fix logger: doesnt flush :- (

	cfg.InfoLogger.Printf("Outcoming request method: %s\n", proxyReq.Method)
	cfg.InfoLogger.Printf("Outcoming host: %s\n", proxyReq.Host)
	cfg.InfoLogger.Printf("Outcoming url: %s\n", proxyReq.URL.String())
	cfg.InfoLogger.Printf("Resposne body: %s", body)
}
