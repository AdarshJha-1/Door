package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"
)

type backend struct {
	URL     string
	isAlive bool
}

func NewBackend(URL string) *backend {
	return &backend{
		URL:     URL,
		isAlive: false,
	}
}

func (b *backend) CheckHealth() {

	backendURL := url.URL{
		Scheme: "http",
		Host:   b.URL,
	}
	newReq, err := http.NewRequestWithContext(context.Background(), http.MethodGet, backendURL.String(), nil)
	if err != nil {
		b.isAlive = false
		return
	}

	newReq.Header.Add("Content/Type", "application/json")

	res, err := client.Do(newReq)
	if err != nil {
		b.isAlive = false
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		b.isAlive = false
		return
	}
	b.isAlive = true
}

type Backends struct {
	backends []*backend
}

var client = &http.Client{}

// my thinking of round robin
var idx uint64 = 0

func (b *Backends) getBackendToReq(ctx context.Context) *backend {

	var backendToHit *backend
	select {
	case <-ctx.Done():
		return nil
	default:
		i := atomic.AddUint64(&idx, 1) - 1
		backendToHit := b.backends[i%uint64(len(b.backends))]
		if backendToHit.isAlive {
			break
		}
	}
	return backendToHit
}

func (b *Backends) proxyHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s | %s \n", r.Method, r.URL)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	backendToHit := b.getBackendToReq(ctx)
	if backendToHit != nil {
		http.Error(w, "all servers are down", http.StatusBadGateway)
		return
	}
	newURL := url.URL{
		Scheme:   "http",
		Host:     backendToHit.URL,
		Path:     r.URL.Path,
		RawQuery: r.URL.RawQuery,
	}

	newReq, err := http.NewRequestWithContext(r.Context(), r.Method, newURL.String(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	newReq.Header = r.Header.Clone()

	res, err := client.Do(newReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	defer res.Body.Close()
	for k, values := range res.Header {
		for _, v := range values {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(res.StatusCode)
	_, err = io.Copy(w, res.Body)
	if err != nil {
		log.Println(err)
	}
}

func main() {

	b := &Backends{
		backends: []*backend{
			NewBackend("localhost:8000"),
			NewBackend("localhost:8001"),
			NewBackend("localhost:8002"),
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", b.proxyHandler)

	server := &http.Server{
		Addr:    ":6969",
		Handler: mux,
	}

	log.Println("server is running...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
