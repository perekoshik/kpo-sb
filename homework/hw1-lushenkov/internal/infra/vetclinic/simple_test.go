package vetclinicimpl

import (
	"context"
	"testing"

	"github.com/perekoshik/zoo/internal/app/config"
	"github.com/perekoshik/zoo/internal/domain/animal"
)

func TestSimpleClinic_NoLimitMeansAlwaysHealthy(t *testing.T) {
	cfg := config.Default()
	clinic := NewSimpleClinic(cfg)

	a, _ := animal.NewPredator("Hungry", "Tiger", 999)
	dec, err := clinic.CheckHealth(context.Background(), a)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if !dec.Healthy {
		t.Fatalf("expected healthy, got deny: %s", dec.Reason)
	}
}

func TestSimpleClinic_WithLimit(t *testing.T) {
	cfg := config.Default()
	cfg.MaxFoodPerDayForAdmission = 50
	clinic := NewSimpleClinic(cfg)

	ok, _ := animal.NewRabbit("Bunny", 10, 8)
	deny, _ := animal.NewElephant("Dumbo", 120, 9)

	dec1, _ := clinic.CheckHealth(context.Background(), ok)
	if !dec1.Healthy {
		t.Fatal("expected healthy")
	}

	dec2, _ := clinic.CheckHealth(context.Background(), deny)
	if dec2.Healthy {
		t.Fatal("expected deny")
	}
	if dec2.Reason == "" {
		t.Fatal("expected reason on deny")
	}
}
