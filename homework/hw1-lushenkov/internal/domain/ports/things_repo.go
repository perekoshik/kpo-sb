package ports

import (
	"context"

	"github.com/perekoshik/zoo/internal/domain/inventory"
)

type ThingEntity interface {
	inventory.InventoryItem

	Name() string
	AssignInventoryNumber(n int) error
}

type ThingRepository interface {
	Add(ctx context.Context, t ThingEntity) (ThingEntity, error)
	List(ctx context.Context) ([]ThingEntity, error)
	Count(ctx context.Context) (int, error)
}
