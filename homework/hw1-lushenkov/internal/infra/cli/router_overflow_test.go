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

func newRouterForFlowTests() *Router {
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

func TestCmdAddThing_AndBalanceReport(t *testing.T) {
	r := newRouterForFlowTests()
	ctx := context.Background()

	out := captureStdout(func() { r.cmdAddThing(ctx, []string{"table", "TN-5"}) })
	if !strings.Contains(out, "OK") {
		t.Fatalf("expected OK for add thing table, got:\n%s", out)
	}
	out = captureStdout(func() { r.cmdAddThing(ctx, []string{"computer", "PC-42"}) })
	if !strings.Contains(out, "OK") {
		t.Fatalf("expected OK for add thing computer, got:\n%s", out)
	}

	_ = captureStdout(func() { r.cmdAddAnimal(ctx, []string{"rabbit", "Bunny", "3", "7"}) })
	_ = captureStdout(func() { r.cmdAddAnimal(ctx, []string{"tiger", "ShereKhan", "30"}) })

	out = captureStdout(func() { r.cmdReportBalance(ctx) })

	for _, s := range []string{"Table", "TN-5", "Computer", "PC-42", "Rabbit", "Bunny", "Tiger", "ShereKhan"} {
		if !strings.Contains(out, s) {
			t.Fatalf("balance report should contain %q, got:\n%s", s, out)
		}
	}
}

func TestReports_FoodAndInteractive(t *testing.T) {
	r := newRouterForFlowTests()
	ctx := context.Background()

	_ = captureStdout(func() { r.cmdAddAnimal(ctx, []string{"rabbit", "RB", "3", "8"}) })
	_ = captureStdout(func() { r.cmdAddAnimal(ctx, []string{"tiger", "TG", "30"}) })

	out := captureStdout(func() { r.cmdReportFood(ctx) })
	low := strings.ToLower(out)
	if !(strings.Contains(low, "33") && (strings.Contains(low, "kg") || strings.Contains(low, "кг"))) {
		t.Fatalf("food report should show total ~33 kg/day, got:\n%s", out)
	}

	out = captureStdout(func() { r.cmdReportInteractive(ctx) })
	if !strings.Contains(out, "RB") || strings.Contains(out, "TG") {
		t.Fatalf("interactive should contain rabbit and exclude tiger, got:\n%s", out)
	}
}

func TestCmdAddThing_UsageErrors(t *testing.T) {
	r := newRouterForFlowTests()
	ctx := context.Background()

	out := captureStdout(func() { r.cmdAddThing(ctx, []string{}) })
	if !(strings.Contains(strings.ToLower(out), "usage") && strings.Contains(out, "thing")) {
		t.Fatalf("expected usage for empty args, got:\n%s", out)
	}

	out = captureStdout(func() { r.cmdAddThing(ctx, []string{"table"}) })
	if !(strings.Contains(strings.ToLower(out), "usage") && strings.Contains(out, "thing")) {
		t.Fatalf("expected usage for single arg, got:\n%s", out)
	}
}

func TestFieldsKeepQuoted_Splits(t *testing.T) {
	line := `add animal rabbit "Bunny Jr." 3 8`
	fields := fieldsKeepQuoted(line)

	if len(fields) < 5 {
		t.Fatalf("expected >=5 fields, got %d: %#v", len(fields), fields)
	}

	found := false
	for _, f := range fields {
		if f == `Bunny Jr.` || f == `"Bunny Jr."` {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("expected joined quoted token 'Bunny Jr.' in fields: %#v", fields)
	}
}
