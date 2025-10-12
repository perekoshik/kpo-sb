package mem

import (
	"context"
	"testing"

	"github.com/perekoshik/zoo/internal/domain/animal"
	"github.com/perekoshik/zoo/internal/domain/thing"
)

func TestInventoryQuery_SkewsAnimalsAndThings(t *testing.T) {
	animals := NewAnimalsRepo()
	things := NewThingsRepo()

	// Add 2 animals
	r1, _ := animal.NewRabbit("Bunny", 3, 7)
	t1, _ := animal.NewTiger("ShereKhan", 30)
	animals.Add(context.Background(), r1)
	animals.Add(context.Background(), t1)

	// Add 2 things
	tb, _ := thing.NewTable("TN-5")
	cam, _ := thing.NewCamera("GateCam-01")
	things.Add(context.Background(), tb)
	things.Add(context.Background(), cam)

	q := NewInventoryQueryMem(animals, things)
	all, err := q.ListAllInventories(context.Background())
	if err != nil {
		t.Fatalf("query err: %v", err)
	}
	if len(all) != 4 {
		t.Fatalf("expected 4 inventory items, got %d", len(all))
	}

	numCount := 0
	nameCount := 0
	for _, it := range all {
		if it.InventoryNumber() > 0 {
			numCount++
		}
		if it.DisplayName() != "" {
			nameCount++
		}
	}
	if numCount != 4 || nameCount != 4 {
		t.Fatalf("unexpected items: nums %d names %d", numCount, nameCount)
	}
}
