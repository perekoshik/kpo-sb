package mem

import (
	"context"
	"fmt"
	"sync"

	"github.com/perekoshik/zoo/internal/domain/ports"
)

type AnimalsRepo struct {
	mu         sync.Mutex
	nextNumber int
	items      []ports.AnimalEntity
	indexByNum map[int]ports.AnimalEntity
}

func NewAnimalsRepo() *AnimalsRepo {
	return &AnimalsRepo{
		nextNumber: 1,
		items:      make([]ports.AnimalEntity, 0, 16),
		indexByNum: make(map[int]ports.AnimalEntity),
	}
}

func (r *AnimalsRepo) Add(ctx context.Context, a ports.AnimalEntity) (ports.AnimalEntity, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	num := r.nextNumber
	r.nextNumber++

	if err := a.AssignInventoryNumber(num); err != nil {
		return nil, fmt.Errorf("assign inventory number: %w", err)
	}
	if _, exists := r.indexByNum[num]; exists {
		return nil, fmt.Errorf("inventory number %d already exists", num)
	}

	r.items = append(r.items, a)
	r.indexByNum[num] = a
	return a, nil
}

func (r *AnimalsRepo) List(ctx context.Context) ([]ports.AnimalEntity, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]ports.AnimalEntity, len(r.items))
	copy(out, r.items)
	return out, nil
}

func (r *AnimalsRepo) Count(ctx context.Context) (int, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.items), nil
}
