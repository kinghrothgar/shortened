package main

import (
	"github.com/bpina/shortened/data"
	"github.com/bpina/shortened/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	pool, err := data.GetConnectionPool()

	if err != nil {
		panic(err)
	}

	data.Database = pool

	router := mux.NewRouter()

	router.HandleFunc("/shorten", handlers.ShortenHandler)
	router.HandleFunc("/{id}", handlers.UrlRedirectHandler)

	http.Handle("/", router)

	panic(http.ListenAndServe(":8080", nil))
}
