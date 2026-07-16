package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func main() {

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

	server := &http.Server{
		Addr:    ":8000",
		Handler: mux,
	}

	log.Println("server is running...")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
