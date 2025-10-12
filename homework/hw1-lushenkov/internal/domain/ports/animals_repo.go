package ports

import (
	"context"

	"github.com/perekoshik/zoo/internal/domain/animal"
	"github.com/perekoshik/zoo/internal/domain/inventory"
)

type AnimalEntity interface {
	animal.Alive
	inventory.InventoryItem

	Name() string
	Species() string
	AssignInventoryNumber(n int) error
}

type AnimalRepository interface {
	Add(ctx context.Context, a AnimalEntity) (AnimalEntity, error)
	List(ctx context.Context) ([]AnimalEntity, error)
	Count(ctx context.Context) (int, error)
}
