package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"toxic/toxics"

	"github.com/charmbracelet/log"
)

type Server struct {
	Name     string  `json:"name"`
	Listen   string  `json:"listen"`
	Upstream string  `json:"upstream"`
	Toxics   []Toxic `json:"toxics"`
}

type Toxic struct {
	Kind        string         `json:"kind"`
	Probability float64        `json:"probability"`
	Spec        map[string]any `json:"spec"`
}

func main() {
	log.Info("starting toxic")

	var configPath string
	flag.StringVar(&configPath, "config", "/config/toxic.json", "path to the config file")
	flag.Parse()

	if configPath == "" {
		panic("config is required")
	}

	config, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Sprintf("failed to read config: %v", err))
	}

	var configFile []Server
	if err := json.Unmarshal(config, &configFile); err != nil {
		panic(fmt.Sprintf("failed to unmarshal config: %v", err))
	}

	var wg sync.WaitGroup
	for _, server := range configFile {
		wg.Add(1)
		go func(server Server) {
			defer wg.Done()
			if err := http.ListenAndServe(server.Listen, createServerHandler(server)); err != nil {
				log.Info("failed to start server", "name", server.Name, "error", err)
			}
		}(server)
	}
	wg.Wait()
}

func createServerHandler(server Server) http.Handler {
	log.Info("creating server", "name", server.Name, "listen", server.Listen, "upstream", server.Upstream)
	upstream, err := url.Parse(server.Upstream)
	if err != nil {
		panic(fmt.Sprintf("failed to parse upstream url: %v", err))
	}
	return wrapHandler(httputil.NewSingleHostReverseProxy(upstream), server.Toxics)
}

func wrapHandler(handler http.Handler, toxiclist []Toxic) http.Handler {
	for i := len(toxiclist) - 1; i >= 0; i-- {
		toxic := toxiclist[i]
		switch toxic.Kind {
		case "delay":
			handler = toxics.DelayHandler(handler, toxic.Probability, toxic.Spec)
		case "status":
			handler = toxics.StatusHandler(handler, toxic.Probability, toxic.Spec)
		case "offline":
			handler = toxics.OfflineHandler(handler, toxic.Probability, toxic.Spec)
		}
	}
	return handler
}
