package tzdata

import (
	"testing"
	"time"
)

func TestGetTimezones(t *testing.T) {
	tzs := GetTimezones()
	if len(tzs) == 0 {
		t.Fatalf("Expected non-empty list of timezones from tzdata, got 0")
	}
	t.Logf("Successfully loaded %d IANA timezones from embedded tzdata", len(tzs))

	// Verify key timezones exist and are loadable by Go's time package
	expected := []string{"Asia/Jakarta", "Asia/Tokyo", "UTC", "America/New_York", "Europe/London"}
	tzMap := make(map[string]bool)
	for _, tz := range tzs {
		tzMap[tz] = true
	}

	for _, exp := range expected {
		if !tzMap[exp] {
			t.Errorf("Expected timezone %s not found in list", exp)
		}
		loc, err := time.LoadLocation(exp)
		if err != nil || loc == nil {
			t.Errorf("Failed to load location for %s with time.LoadLocation: %v", exp, err)
		}
	}
}
