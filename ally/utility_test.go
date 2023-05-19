package ally

import (
	"testing"
	"time"
)

func TestAPIStatus(t *testing.T) {
	got := APIStatus().Time.Time.Round(time.Minute)
	now := time.Now().Round(time.Minute)

	if !got.Equal(now) {
		t.Errorf("ServerStatus() = %v, want %v", got, now)
	}
}

func TestAPIVersion(t *testing.T) {
	got := APIVersion()
	if got.Version != "1.0" {
		t.Errorf("ServerVersion() = %v, want %v", got, "1.0")
	}
}
