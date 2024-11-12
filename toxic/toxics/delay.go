package toxics

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/exp/rand"
)

type ToxicDelay struct {
	Delay  time.Duration `mapstructure:"delay"`
	Jitter time.Duration `mapstructure:"jitter"`
}

func DelayHandler(next http.Handler, probability float64, rawSpec map[string]any) http.Handler {
	var spec ToxicDelay
	if err := mapAnyToStruct(rawSpec, &spec); err != nil {
		panic(fmt.Sprintf("failed to decode delay spec: %v", err))
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rand.Float64() > probability {
			next.ServeHTTP(w, r)
			return
		}
		time.Sleep(spec.Delay + time.Duration(rand.Intn(int(spec.Jitter))))
		next.ServeHTTP(w, r)
	})
}
