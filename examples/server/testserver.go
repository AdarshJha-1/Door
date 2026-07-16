package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 3 {
		log.Fatal("args 1st -> PORT (:8000), 2nd -> name, ", len(args))
	}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /users", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.String())
		w.WriteHeader(http.StatusOK)
		res := struct {
			Users []string `json:"users"`
		}{
			Users: []string{
				"Luffy",
				"Zoro",
				"Ms. UwU",
				"Babu Rao",
			},
		}

		if err := json.NewEncoder(w).Encode(&res); err != nil {
			http.Error(w, "error res json", http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.String())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Healthy Server"))
	})

	server := &http.Server{
		Addr:    args[1],
		Handler: mux,
	}

	log.Printf("server %s is running...\n", args[2])
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
