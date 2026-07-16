package node

import (
	"net/http"
	"net/url"
	"sync/atomic"
)

type Node struct {
	URL   string
	Alive atomic.Bool
}

func New(URL string) *Node {
	return &Node{
		URL:   URL,
		Alive: atomic.Bool{},
	}
}

func (b *Node) CheckHealth(client *http.Client) {

	backendURL := url.URL{
		Scheme: "http",
		Host:   b.URL,
		Path:   "/health", // hardcoding it
	}
	newReq, err := http.NewRequest(http.MethodGet, backendURL.String(), nil)
	if err != nil {
		b.Alive.Store(false)
		return
	}

	newReq.Header.Add("Content-Type", "application/json")

	res, err := client.Do(newReq)
	if err != nil {
		b.Alive.Store(false)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		b.Alive.Store(false)
		return
	}
	b.Alive.Store(true)
}
