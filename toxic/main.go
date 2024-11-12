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

	"github.com/charmbracelet/log"
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
	Delay  string `mapstructure:"delay"`
	Jitter string `mapstructure:"jitter"`
}

type ToxicOffline struct {
	// Interval is how often the offline period will trigger
	Interval string `mapstructure:"interval"`
	// Duration is how long the offline period will last
	Duration string `mapstructure:"duration"`
	// StartJitter is the amount of time to wait before starting the offline period
	StartJitter string `mapstructure:"start_jitter"`
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
			if err := http.ListenAndServe(server.Listen, serverHandler(server)); err != nil {
				log.Info("failed to start server", "name", server.Name, "error", err)
			}
		}(server)
	}
	wg.Wait()
}

func serverHandler(server Server) http.Handler {
	log.Info("creating server", "name", server.Name, "listen", server.Listen, "upstream", server.Upstream)
	upstream, err := url.Parse(server.Upstream)
	if err != nil {
		panic(fmt.Sprintf("failed to parse upstream url: %v", err))
	}
	proxy := httputil.NewSingleHostReverseProxy(upstream)

	var handler http.Handler = proxy

	// Iterate through toxics in reverse order
	for i := len(server.Toxics) - 1; i >= 0; i-- {
		toxic := server.Toxics[i]
		switch toxic.Kind {
		case "delay":
			handler = delayHandler(handler, toxic.Probability, delaySpecFromAny(toxic.Spec))
		case "status":
			handler = statusHandler(handler, toxic.Probability, statusSpecFromAny(toxic.Spec))
		case "offline":
			handler = offlineHandler(handler, toxic.Probability, offlineSpecFromAny(toxic.Spec))
		}
	}

	return handler
}

func delaySpecFromAny(spec any) ToxicDelay {
	var result ToxicDelay
	if err := mapstructure.Decode(spec, &result); err != nil {
		panic(fmt.Sprintf("failed to decode delay spec: %v", err))
	}
	return result
}

func delayHandler(next http.Handler, probability float64, spec ToxicDelay) http.Handler {
	delay, err := time.ParseDuration(spec.Delay)
	if err != nil {
		panic(fmt.Sprintf("failed to parse delay: %v", err))
	}
	jitter, err := time.ParseDuration(spec.Jitter)
	if err != nil {
		panic(fmt.Sprintf("failed to parse jitter: %v", err))
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rand.Float64() > probability {
			next.ServeHTTP(w, r)
			return
		}
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

func statusHandler(next http.Handler, probability float64, spec ToxicStatus) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rand.Float64() > probability {
			next.ServeHTTP(w, r)
			return
		}
		w.WriteHeader(spec.Status)
		w.Write([]byte(spec.Message))
	})
}

func offlineSpecFromAny(spec any) ToxicOffline {
	var result ToxicOffline
	if err := mapstructure.Decode(spec, &result); err != nil {
		panic(fmt.Sprintf("failed to decode offline spec: %v", err))
	}
	return result
}

// offlineHandler will return 503s during an offline window
func offlineHandler(next http.Handler, probability float64, spec ToxicOffline) http.Handler {
	var (
		offline bool
	)
	duration, err := time.ParseDuration(spec.Duration)
	if err != nil {
		panic(fmt.Sprintf("failed to parse duration: %v", err))
	}
	startJitter, err := time.ParseDuration(spec.StartJitter)
	if err != nil {
		panic(fmt.Sprintf("failed to parse start jitter: %v", err))
	}
	interval, err := time.ParseDuration(spec.Interval)
	if err != nil {
		panic(fmt.Sprintf("failed to parse interval: %v", err))
	}

	randomStartDelay := time.Duration(rand.Intn(int(startJitter)))

	log.Info("starting offline handler", "duration", duration, "interval", interval, "delay", randomStartDelay)
	go func() {
		time.Sleep(randomStartDelay)
		for range time.Tick(interval) {
			offline = true
			log.Info("going offline", "duration", duration, "hostname", os.Getenv("HOSTNAME"))
			time.Sleep(duration)
			offline = false
			log.Info("back online", "hostname", os.Getenv("HOSTNAME"))
		}
	}()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rand.Float64() > probability {
			next.ServeHTTP(w, r)
			return
		}
		// Return a 503 if we're offline
		if offline {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("Service Unavailable"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
