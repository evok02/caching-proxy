package proxy

type ErrorResposne struct {
	Error string `json:"error"`
}

type ResponseCache struct {
	Status        string              `json:"status"`
	StatusCode    int                 `json:"status_code"`
	Header        map[string][]string `json:"header"`
	ContentLength int64               `json:"content_length"`
}

type RequestCache struct {
	Method string              `json:"method"`
	URL    string              `json:"url"`
	Origin string              `json:"origin"`
	Header map[string][]string `json:"header"`
}
