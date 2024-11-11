package main

import (
	"fmt"
	"net/http"
	"os"
	"sort"
)

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	content := fmt.Sprintf(`
Pod details:
  Hostname: %s

Request details:
  Host: %s
  Path: %s
  Method: %s
  Protocol: %s
  RemoteAddr: %s
  RequestURI: %s
  Headers:
`, os.Getenv("HOSTNAME"), r.Host, r.URL.Path, r.Method, r.Proto, r.RemoteAddr, r.RequestURI)
	for _, k := range sortedHeaderKeys(r.Header) {
		content += fmt.Sprintf("    %s: %s\n", k, r.Header.Get(k))
	}

	fmt.Fprintln(w, content)
}

func sortedHeaderKeys(headers http.Header) []string {
	keys := make([]string, 0, len(headers))
	for k := range headers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
