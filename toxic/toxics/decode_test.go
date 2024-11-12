package toxics

import (
	"testing"
	"time"
)

func TestStringToDurationHookFunc(t *testing.T) {
	spec := map[string]any{
		"delay":  "1s",
		"jitter": "2s",
	}
	var toxic ToxicDelay
	if err := mapAnyToStruct(spec, &toxic); err != nil {
		t.Fatalf("failed to decode delay spec: %v", err)
	}
	if toxic.Delay != time.Second {
		t.Errorf("expected delay to be 1s, got %v", toxic.Delay)
	}
	if toxic.Jitter != 2*time.Second {
		t.Errorf("expected jitter to be 2s, got %v", toxic.Jitter)
	}
}
