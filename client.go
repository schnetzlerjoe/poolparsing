package client

var (
  "net/http"
)

var client = &http.Client{Timeout: time.Second * 10}