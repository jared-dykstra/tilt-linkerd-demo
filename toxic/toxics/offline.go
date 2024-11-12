package toxics

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"golang.org/x/exp/rand"
)

type ToxicOffline struct {
	// Interval is how often the offline period will trigger
	Interval time.Duration `mapstructure:"interval"`
	// Duration is how long the offline period will last
	Duration time.Duration `mapstructure:"duration"`
	// StartJitter is the amount of time to wait before starting the offline period
	StartJitter time.Duration `mapstructure:"start_jitter"`
}

// OfflineHandler will return 503s during an offline window
func OfflineHandler(next http.Handler, probability float64, rawSpec map[string]any) http.Handler {
	var spec ToxicOffline
	if err := mapAnyToStruct(rawSpec, &spec); err != nil {
		panic(fmt.Sprintf("failed to decode offline spec: %v", err))
	}

	var (
		offline bool
	)

	randomStartDelay := time.Duration(rand.Intn(int(spec.StartJitter)))

	log.Info("starting offline handler", "duration", spec.Duration, "interval", spec.Interval, "delay", randomStartDelay)
	go func() {
		time.Sleep(randomStartDelay)
		for range time.Tick(spec.Interval) {
			offline = true
			log.Info("going offline", "duration", spec.Duration, "hostname", os.Getenv("HOSTNAME"))
			time.Sleep(spec.Duration)
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
