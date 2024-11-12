package toxics

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDelayHandler(t *testing.T) {
	handler := DelayHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}), 1, map[string]any{"delay": "5ms", "jitter": "5ms"})

	recorder := httptest.NewRecorder()
	start := time.Now()
	handler.ServeHTTP(recorder, httptest.NewRequest("GET", "/", nil))
	elapsed := time.Since(start)
	if recorder.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, recorder.Code)
	}
	if elapsed < 5*time.Millisecond {
		t.Errorf("expected delay to be at least 10ms, got %v", elapsed)
	}
}
