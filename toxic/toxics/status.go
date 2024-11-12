package toxics

import (
	"fmt"
	"net/http"

	"golang.org/x/exp/rand"
)

type ToxicStatus struct {
	Status  int    `mapstructure:"status"`
	Message string `mapstructure:"message"`
}

func StatusHandler(next http.Handler, probability float64, rawSpec map[string]any) http.Handler {
	var spec ToxicStatus
	if err := mapAnyToStruct(rawSpec, &spec); err != nil {
		panic(fmt.Sprintf("failed to decode status spec: %v", err))
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rand.Float64() > probability {
			next.ServeHTTP(w, r)
			return
		}
		w.WriteHeader(spec.Status)
		w.Write([]byte(spec.Message))
	})
}
