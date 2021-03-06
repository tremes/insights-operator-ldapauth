package server

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

// HandleHTTP handle all routes, used for proxying them to controller
func (s Server) HandleHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("Proxying to", s.Proxy)
	tempURL, _ := url.Parse(s.Proxy)
	req.URL.Host = tempURL.Host
	req.URL.Scheme = tempURL.Scheme
	req.Host = tempURL.Host
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
