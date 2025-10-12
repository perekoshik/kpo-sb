package ports

import (
	"context"

	"github.com/perekoshik/zoo/internal/domain/inventory"
)

type InventoryQuery interface {
	ListAllInventories(ctx context.Context) ([]inventory.InventoryItem, error)
}
