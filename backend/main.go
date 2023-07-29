package main

import (
	"log"
	"net/http"
)

var (
	portHTTP     string = "8080"
	CustomRouter *http.ServeMux
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	corsMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Auth-Token, X-Requested-With, Application, json")

			// if i remember correct, firefox lags happens without this
			if r.Method == "OPTIONS" {
				log.Println("pre-fetch request") // todo: remove debug
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	CustomRouter = http.NewServeMux()
	registerHandlers()
	log.Println("starting forum at http://localhost:" + portHTTP + "/")
	log.Println("starting websocket at ws://localhost:" + portHTTP + "/ws")
	log.Fatal(http.ListenAndServe(":"+portHTTP, corsMiddleware(CustomRouter)))
}
