package main

import (
	"log"
	"net/http"
	"time"

	"github.com/AdarshJha-1/Door/internal/backend"
	"github.com/AdarshJha-1/Door/internal/node"
)

func main() {

	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	nodes := []*node.Node{
		node.New("localhost:8000"),
		node.New("localhost:8001"),
		node.New("localhost:8002"),
	}
	backends := backend.New(nodes, client)

	backends.StartHealthChecker(5 * time.Second)

	mux := http.NewServeMux()
	mux.HandleFunc("/", backends.ProxyHandler)

	server := &http.Server{
		Addr:    ":6969",
		Handler: mux,
	}

	log.Println("server is running...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
