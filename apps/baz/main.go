package main

import (
	"fmt"
	"net/http"

	"math/rand"
)

var (
	flakeProbability = 50
)

func main() {
	http.ListenAndServe(":8080", flakeWrapper(http.HandlerFunc(handler)))
}

func flakeWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rand.Intn(100) < flakeProbability {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello from baz!")
}
