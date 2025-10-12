package cli

import (
	"bytes"
	"context"
	"os"
	"testing"

	"github.com/perekoshik/zoo/internal/app/config"
	"github.com/perekoshik/zoo/internal/app/reports"
	"github.com/perekoshik/zoo/internal/domain/animal"
	"github.com/perekoshik/zoo/internal/domain/thing"
	"github.com/perekoshik/zoo/internal/infra/repos/mem"
)

func captureStdout(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	fn()
	_ = w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String()
}

func TestRouter_ReportCommands_Print_NoFakes(t *testing.T) {
	ctx := context.Background()

	animalsRepo := mem.NewAnimalsRepo()
	thingsRepo := mem.NewThingsRepo()

	rabbit, _ := animal.NewRabbit("Bunny", 3, 8)
	tiger, _ := animal.NewTiger("ShereKhan", 30)
	table, _ := thing.NewTable("TN-5")

	if _, err := animalsRepo.Add(ctx, rabbit); err != nil {
		t.Fatal(err)
	}
	if _, err := animalsRepo.Add(ctx, tiger); err != nil {
		t.Fatal(err)
	}
	if _, err := thingsRepo.Add(ctx, table); err != nil {
		t.Fatal(err)
	}

	invQ := mem.NewInventoryQueryMem(animalsRepo, thingsRepo)
	foodSvc := reports.NewFoodService(animalsRepo)
	interSvc := reports.NewInteractiveService(config.Default(), animalsRepo)
	balanceSvc := reports.NewBalanceService(invQ)

	r := &Router{
		cfg:        config.Default(),
		animals:    animalsRepo,
		things:     thingsRepo,
		invQuery:   invQ,
		foodSvc:    foodSvc,
		interSvc:   interSvc,
		balanceSvc: balanceSvc,
	}

	out := captureStdout(func() { r.cmdReportBalance(ctx) })
	if !bytes.Contains([]byte(out), []byte("#1 Table TN-5")) {
		t.Fatalf("balance: expected Table TN-5, got %q", out)
	}

	out = captureStdout(func() { r.cmdReportFood(ctx) })
	if !bytes.Contains([]byte(out), []byte("TOTAL: 33 kg/day")) {
		t.Fatalf("food: expected total 33, got %q", out)
	}

	out = captureStdout(func() { r.cmdReportInteractive(ctx) })
	if !bytes.Contains([]byte(out), []byte("kindness 8")) {
		t.Fatalf("interactive: expected kindness 8, got %q", out)
	}
}

func TestRouter_ReportEmpty_NoFakes(t *testing.T) {
	ctx := context.Background()

	animalsRepo := mem.NewAnimalsRepo()
	thingsRepo := mem.NewThingsRepo()
	invQ := mem.NewInventoryQueryMem(animalsRepo, thingsRepo)

	r := &Router{
		cfg:        config.Default(),
		animals:    animalsRepo,
		things:     thingsRepo,
		invQuery:   invQ,
		foodSvc:    reports.NewFoodService(animalsRepo),
		interSvc:   reports.NewInteractiveService(config.Default(), animalsRepo),
		balanceSvc: reports.NewBalanceService(invQ),
	}

	check := func(fn func(context.Context), want string) {
		out := captureStdout(func() { fn(ctx) })
		if !bytes.Contains([]byte(out), []byte(want)) {
			t.Fatalf("want %q, got %q", want, out)
		}
	}
	check(r.cmdReportBalance, "(empty balance)")
	check(r.cmdReportFood, "(no animals)")
	check(r.cmdReportInteractive, "(no animals eligible")
}
