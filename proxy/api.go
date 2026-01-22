package proxy

import (
	"encoding/json"
	"fmt"
	"github.com/evok02/cacher/storage"
	"log"
	"net/http"
	"net/url"
	"os"
)

type ApiConfig struct {
	Storage storage.RedisStorage
}

type ErrorResposne struct {
	Error string `json:"error"`
}

func NewApiConfig(rdb storage.RedisStorage) *ApiConfig {
	return &ApiConfig{
		Storage: rdb,
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

	log.Println(r.URL.Path)
	targetUrl, err := url.JoinPath(target, r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	parsedUrl, err := url.Parse(targetUrl)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	log.Println(parsedUrl.String())

	proxyReq := &http.Request{
		Method: r.Method,
		Header: r.Header,
		URL:    parsedUrl,
		Body:   r.Body,
	}

	res, err := http.DefaultClient.Do(proxyReq)
	if err != nil {
		writeError(w, 500, err)
		return
	}

	var body []byte
	_, err = res.Body.Read(body)
	if err != nil {
		writeError(w, 500, err)
		return
	}

	_, err = w.Write(body)
	if err != nil {
		writeError(w, 500, err)
		return
	}

	fmt.Printf("this is the body %s: ", body)

}
