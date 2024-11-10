package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
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

type ToxicStatus struct {
	Status  int    `mapstructure:"status"`
	Message string `mapstructure:"message"`
}

type ToxicDelay struct {
	Delay    string `mapstructure:"delay"`
	JitterMs string `mapstructure:"jitter"`
}

func log(message string, args ...any) {
	output := message
	for i := 0; i < len(args); i += 2 {
		output = fmt.Sprintf("%s %s=\"%v\"", output, args[i], args[i+1])
	}
	fmt.Println(output)
}

func main() {
	log("starting toxic")

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
			if err := http.ListenAndServe(server.Listen, serverHandler(server)); err != nil {
				log("failed to start server", "name", server.Name, "error", err)
			}
		}(server)
	}
	wg.Wait()
}

func serverHandler(server Server) http.Handler {
	log("creating server", "name", server.Name, "listen", server.Listen, "upstream", server.Upstream)
	upstream, err := url.Parse(server.Upstream)
	if err != nil {
		panic(fmt.Sprintf("failed to parse upstream url: %v", err))
	}
	proxy := httputil.NewSingleHostReverseProxy(upstream)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log("handling request", "name", server.Name)
		var handler http.Handler = proxy
		for i := len(server.Toxics) - 1; i >= 0; i-- {
			if rand.Float64() > server.Toxics[i].Probability {
				continue
			}
			toxic := server.Toxics[i]
			switch toxic.Kind {
			case "delay":
				handler = delayHandler(handler, delaySpecFromAny(toxic.Spec))
			case "status":
				handler = statusHandler(handler, statusSpecFromAny(toxic.Spec))
			}
		}
		handler.ServeHTTP(w, r)
	})
}

func delaySpecFromAny(spec any) ToxicDelay {
	var result ToxicDelay
	if err := mapstructure.Decode(spec, &result); err != nil {
		panic(fmt.Sprintf("failed to decode delay spec: %v", err))
	}
	return result
}

func delayHandler(next http.Handler, spec ToxicDelay) http.Handler {
	delay, err := time.ParseDuration(spec.Delay)
	if err != nil {
		panic(fmt.Sprintf("failed to parse delay: %v", err))
	}
	jitter, err := time.ParseDuration(spec.JitterMs)
	if err != nil {
		panic(fmt.Sprintf("failed to parse jitter: %v", err))
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log("delaying request", "delay", spec.Delay, "jitter", spec.JitterMs)
		time.Sleep(delay + time.Duration(rand.Intn(int(jitter))))
		next.ServeHTTP(w, r)
	})
}

func statusSpecFromAny(spec any) ToxicStatus {
	var result ToxicStatus
	if err := mapstructure.Decode(spec, &result); err != nil {
		panic(fmt.Sprintf("failed to decode status spec: %v", err))
	}
	return result
}

func statusHandler(next http.Handler, spec ToxicStatus) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log("returning status", "status", spec.Status, "message", spec.Message)
		w.WriteHeader(spec.Status)
		w.Write([]byte(spec.Message))
	})
}
