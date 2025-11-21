package reports

import (
	"context"
	"testing"

	"github.com/perekoshik/zoo/internal/domain/animal"
	"github.com/perekoshik/zoo/internal/domain/thing"
	"github.com/perekoshik/zoo/internal/infra/repos/mem"
)

func TestBalance_Empty(t *testing.T) {
	q := mem.NewInventoryQueryMem(mem.NewAnimalsRepo(), mem.NewThingsRepo())
	svc := NewBalanceService(q)

	rows, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if len(rows) != 0 {
		t.Fatalf("want empty, got %d", len(rows))
	}
}

func TestBalance_SortedByInventoryNumber(t *testing.T) {
	ar := mem.NewAnimalsRepo()
	tr := mem.NewThingsRepo()

	r1, _ := animal.NewRabbit("Bunny", 3, 7)  // inv #1 (animals)
	r2, _ := animal.NewTiger("ShereKhan", 30) // inv #2 (animals)
	ar.Add(context.Background(), r1)
	ar.Add(context.Background(), r2)

	tb, _ := thing.NewTable("TN-5")      // inv #1 (things)
	cam, _ := thing.NewCamera("Gate-01") // inv #2 (things)
	tr.Add(context.Background(), tb)
	tr.Add(context.Background(), cam)

	q := mem.NewInventoryQueryMem(ar, tr)
	svc := NewBalanceService(q)
	rows, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if len(rows) != 4 {
		t.Fatalf("want 4, got %d", len(rows))
	}

	for i := 1; i < len(rows); i++ {
		if rows[i-1].InventoryNumber > rows[i].InventoryNumber {
			t.Fatalf("not sorted at %d: %v > %v", i, rows[i-1], rows[i])
		}
	}
}
