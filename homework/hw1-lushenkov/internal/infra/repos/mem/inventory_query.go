package mem

import (
	"context"

	"github.com/perekoshik/zoo/internal/domain/inventory"
	"github.com/perekoshik/zoo/internal/domain/ports"
)

type InventoryQueryMem struct {
	animals ports.AnimalRepository
	things  ports.ThingRepository
}

func NewInventoryQueryMem(animals ports.AnimalRepository, things ports.ThingRepository) *InventoryQueryMem {
	return &InventoryQueryMem{animals: animals, things: things}
}

func (q *InventoryQueryMem) ListAllInventories(ctx context.Context) ([]inventory.InventoryItem, error) {
	as, err := q.animals.List(ctx)
	if err != nil {
		return nil, err
	}
	ts, err := q.things.List(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]inventory.InventoryItem, 0, len(as)+len(ts))
	for _, a := range as {
		out = append(out, a)
	}
	for _, t := range ts {
		out = append(out, t)
	}
	return out, nil
}
