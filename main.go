package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	http := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	http.ListenAndServe()
}
