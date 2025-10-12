package reports_test

import (
	"context"
	"testing"

	"github.com/perekoshik/zoo/internal/app/config"
	"github.com/perekoshik/zoo/internal/app/reports"
	"github.com/perekoshik/zoo/internal/domain/animal"
	"github.com/perekoshik/zoo/internal/infra/repos/mem"
)

func TestInteractive_FilterByKindnessAndHerbivore(t *testing.T) {
	repo := mem.NewAnimalsRepo()
	cfg := config.Default()

	h1, _ := animal.NewHerbivore("Nice", "Rabbit", 2, 6)
	h2, _ := animal.NewHerbivore("Border", "Rabbit", 2, 5)
	p1, _ := animal.NewPredator("Hunter", "Tiger", 25)

	repo.Add(context.Background(), h1)
	repo.Add(context.Background(), h2)
	repo.Add(context.Background(), p1)

	svc := reports.NewInteractiveService(cfg, repo)
	rows, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 eligible, got %d", len(rows))
	}
	if rows[0].DisplayName != h1.DisplayName() {
		t.Fatalf("unexpected row: %+v", rows[0])
	}
}
