package main

import (
	"fmt"
	"io"
	"net/http"
	"sort"
)

func main() {
	http.ListenAndServe(":8080", http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://baz/foo/bar")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	content := "Hello from bar! Here are the request details:\n\n"
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	content += fmt.Sprintf("Response from baz:\nstatus: %s\nbody: %s\n", resp.Status, string(body))

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
