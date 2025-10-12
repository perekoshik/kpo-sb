package config

import "testing"

func TestDefaultAndAdmissionLimitFlag(t *testing.T) {
	c := Default()
	if c.InteractiveKindnessThreshold != 6 {
		t.Fatalf("want threshold 6, got %d", c.InteractiveKindnessThreshold)
	}
	if c.AdmissionLimitEnabled() {
		t.Fatal("limit must be disabled by default")
	}
	c.MaxFoodPerDayForAdmission = 10
	if !c.AdmissionLimitEnabled() {
		t.Fatal("limit must be enabled after set")
	}
}
