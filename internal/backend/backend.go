package backend

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/AdarshJha-1/Door/internal/node"
)

type Backends struct {
	Nodes     []*node.Node
	Client    *http.Client
	NodeCount uint64
	idx       uint64
}

func New(nodes []*node.Node, client *http.Client) *Backends {
	return &Backends{
		Nodes:     nodes,
		Client:    client,
		NodeCount: uint64(len(nodes)),
		idx:       0,
	}
}

func (b *Backends) getBackendToReq() *node.Node {
	for range b.Nodes {
		i := atomic.AddUint64(&b.idx, 1) - 1
		backendToHit := b.Nodes[i%b.NodeCount]
		if backendToHit.Alive.Load() {
			return backendToHit
		}
	}

	return nil
}

func (b *Backends) ProxyHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s | %s \n", r.Method, r.URL)

	for attempts := 0; attempts < int(b.NodeCount); attempts++ {
		backendToHit := b.getBackendToReq()
		if backendToHit == nil {
			break
		}
		newURL := url.URL{
			Scheme:   "http",
			Host:     backendToHit.URL,
			Path:     r.URL.Path,
			RawQuery: r.URL.RawQuery,
		}

		newReq, err := http.NewRequestWithContext(r.Context(), r.Method, newURL.String(), r.Body)
		if err != nil {
			backendToHit.Alive.Store(false)
			continue
		}

		newReq.Header = r.Header.Clone()

		res, err := b.Client.Do(newReq)
		if err != nil {
			backendToHit.Alive.Store(false)
			continue
		}
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

		res.Body.Close()
		return
	}
	http.Error(w, "all servers are down", http.StatusBadGateway)
}

func (b *Backends) StartHealthChecker(seconds time.Duration) {
	go func() {
		ticker := time.NewTicker(seconds)
		defer ticker.Stop()
		for {
			for _, n := range b.Nodes {
				go n.CheckHealth(b.Client)
			}

			<-ticker.C
		}
	}()
}
