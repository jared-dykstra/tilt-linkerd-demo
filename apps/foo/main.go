package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	barURL = getEnv("BAR_URL", "http://bar")
	bazURL = getEnv("BAZ_URL", "http://baz")
)

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	barChan, bazChan := make(chan string), make(chan string)
	go func() {
		barChan <- getBar()
	}()
	go func() {
		bazChan <- getBaz()
	}()
	fmt.Fprintf(w, "Hello, World!\n")
	fmt.Fprintf(w, "bar=%s\n", <-barChan)
	fmt.Fprintf(w, "baz=%s\n", <-bazChan)
}

func getResource(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body)
}

func getBar() string {
	return getResource(barURL)
}

func getBaz() string {
	return getResource(bazURL)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
