package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	barURL = getEnv("BAR_URL", "http://bar")
	bazURL = getEnv("BAZ_URL", "http://baz/")
)

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	barChan, bazChan := make(chan Response), make(chan Response)
	go func() {
		barChan <- getBar()
	}()
	go func() {
		bazChan <- getBaz()
	}()

	barResponse, bazResponse := <-barChan, <-bazChan

	if barResponse.Status != http.StatusOK || bazResponse.Status != http.StatusOK {
		w.WriteHeader(http.StatusInternalServerError)
	}

	fmt.Fprintf(w, "%s\n%s", barResponse, bazResponse)
}

func getResource(url string) Response {
	resp, err := http.Get(url)
	if err != nil {
		return Response{
			URL:   url,
			Error: err,
		}
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{
			URL:   url,
			Error: err,
		}
	}
	return Response{
		Status: resp.StatusCode,
		URL:    url,
		Body:   string(body),
	}
}

func getBar() Response {
	return getResource(barURL)
}

func getBaz() Response {
	return getResource(bazURL)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

type Response struct {
	Status int
	URL    string
	Body   string
	Error  error
}

func (r Response) String() string {
	return fmt.Sprintf(`
---
   url: %s
status: %d
  body: %s
 error: %v
	`, r.URL, r.Status, r.Body, r.Error)
}
