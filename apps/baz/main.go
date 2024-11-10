package main

import (
	"fmt"
	"net/http"
	"sort"
)

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	content := "Hello from baz! Here are the request details:\n\n"
	content += fmt.Sprintf("Host: %s\n", r.Host)
	content += fmt.Sprintf("Path: %s\n", r.URL.Path)
	content += fmt.Sprintf("Method: %s\n", r.Method)
	content += fmt.Sprintf("Protocol: %s\n", r.Proto)
	content += fmt.Sprintf("RemoteAddr: %s\n", r.RemoteAddr)
	content += fmt.Sprintf("RequestURI: %s\n", r.RequestURI)

	content += "Headers:\n"
	for _, k := range sortedHeaderKeys(r.Header) {
		content += fmt.Sprintf("  %s: %s\n", k, r.Header.Get(k))
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
