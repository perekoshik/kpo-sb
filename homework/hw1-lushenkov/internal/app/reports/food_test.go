package reports_test

import (
	"context"
	"testing"

	"github.com/perekoshik/zoo/internal/app/reports"
	"github.com/perekoshik/zoo/internal/domain/animal"
	"github.com/perekoshik/zoo/internal/infra/repos/mem"
)

func TestFoodSummary_Empty(t *testing.T) {
	repo := mem.NewAnimalsRepo()
	svc := reports.NewFoodService(repo)

	sum, err := svc.Summary(context.Background())
	if err != nil {
		t.Fatalf("summary error: %v", err)
	}
	if sum.TotalKgPerDay != 0 || len(sum.Rows) != 0 {
		t.Fatalf("expected empty summary, got %+v", sum)
	}
}

func TestFoodSummary_NonEmpty(t *testing.T) {
	repo := mem.NewAnimalsRepo()
	svc := reports.NewFoodService(repo)

	r1, _ := animal.NewRabbit("Bunny", 3, 7)
	r2, _ := animal.NewTiger("ShereKhan", 30)
	repo.Add(context.Background(), r1)
	repo.Add(context.Background(), r2)

	sum, err := svc.Summary(context.Background())
	if err != nil {
		t.Fatalf("summary error: %v", err)
	}
	if sum.TotalKgPerDay != 33 {
		t.Fatalf("expected 33, got %d", sum.TotalKgPerDay)
	}
	if len(sum.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(sum.Rows))
	}
}
