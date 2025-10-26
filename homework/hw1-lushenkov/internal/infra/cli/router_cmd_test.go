package cli

import (
	"context"
	"strings"
	"testing"

	"github.com/perekoshik/zoo/internal/app/config"
	"github.com/perekoshik/zoo/internal/app/onboarding"
	"github.com/perekoshik/zoo/internal/app/reports"
	"github.com/perekoshik/zoo/internal/infra/repos/mem"
	vetclinicimpl "github.com/perekoshik/zoo/internal/infra/vetclinic"
)

func newRouterForAddTests() *Router {
	animals := mem.NewAnimalsRepo()
	things := mem.NewThingsRepo()
	invQ := mem.NewInventoryQueryMem(animals, things)

	cfg := config.Default()
	clinic := vetclinicimpl.NewSimpleClinic(cfg)
	onboard := onboarding.NewService(clinic, animals)

	return &Router{
		cfg:        cfg,
		onboard:    onboard,
		animals:    animals,
		things:     things,
		invQuery:   invQ,
		foodSvc:    reports.NewFoodService(animals),
		interSvc:   reports.NewInteractiveService(cfg, animals),
		balanceSvc: reports.NewBalanceService(invQ),
	}
}

func containsAny(s string, needles ...string) bool {
	for _, n := range needles {
		if strings.Contains(s, n) {
			return true
		}
	}
	return false
}

func TestCmdAddAnimal_HappyAndErrors(t *testing.T) {
	r := newRouterForAddTests()
	ctx := context.Background()

	out := captureStdout(func() {
		r.cmdAddAnimal(ctx, []string{"rabbit", "Bunny", "3", "8"})
	})
	if !(strings.Contains(out, "OK") && strings.Contains(out, "Rabbit") && strings.Contains(out, "Bunny")) {
		t.Fatalf("unexpected output for herbivore happy path:\n%s", out)
	}

	out = captureStdout(func() {
		r.cmdAddAnimal(ctx, []string{"tiger", "ShereKhan", "30"})
	})
	if !(strings.Contains(out, "OK") && strings.Contains(out, "Tiger") && strings.Contains(out, "ShereKhan")) {
		t.Fatalf("unexpected output for predator happy path:\n%s", out)
	}

	out = captureStdout(func() {
		r.cmdAddAnimal(ctx, []string{})
	})
	if !containsAny(out, "usage", "Usage", "USAGE") || !strings.Contains(out, "animal") {
		t.Fatalf("expected usage for empty args, got:\n%s", out)
	}

	out = captureStdout(func() {
		r.cmdAddAnimal(ctx, []string{"tiger", "X", "-1"})
	})
	low := strings.ToLower(out)
	if !containsAny(low, "err", "error", "invalid", "food") {
		t.Fatalf("expected error for negative food, got:\n%s", out)
	}

	out = captureStdout(func() {
		r.cmdAddAnimal(ctx, []string{"rabbit", "RB", "1"})
	})
	low = strings.ToLower(out)
	if !containsAny(low, "usage", "kindness", "arg") {
		t.Fatalf("expected error for missing kindness, got:\n%s", out)
	}

	out = captureStdout(func() {
		r.cmdAddAnimal(ctx, []string{"rabbit", "RB", "1", "11"})
	})
	low = strings.ToLower(out)
	if !containsAny(low, "kindness", "0..10", "range", "invalid", "error") {
		t.Fatalf("expected error for kindness out of range, got:\n%s", out)
	}
}
