package mem

import (
	"context"
	"fmt"
	"sync"

	"github.com/perekoshik/zoo/internal/domain/ports"
)

type ThingsRepo struct {
	mu         sync.Mutex
	nextNumber int
	items      []ports.ThingEntity
	indexByNum map[int]ports.ThingEntity
}

func NewThingsRepo() *ThingsRepo {
	return &ThingsRepo{
		nextNumber: 1,
		items:      make([]ports.ThingEntity, 0, 16),
		indexByNum: make(map[int]ports.ThingEntity),
	}
}

func (r *ThingsRepo) Add(ctx context.Context, t ports.ThingEntity) (ports.ThingEntity, error) {
	_ = ctx

	r.mu.Lock()
	defer r.mu.Unlock()

	num := r.nextNumber
	r.nextNumber++

	if err := t.AssignInventoryNumber(num); err != nil {
		return nil, fmt.Errorf("assign inventory number: %w", err)
	}
	if _, exists := r.indexByNum[num]; exists {
		return nil, fmt.Errorf("inventory number %d already exists", num)
	}

	r.items = append(r.items, t)
	r.indexByNum[num] = t
	return t, nil
}

func (r *ThingsRepo) List(ctx context.Context) ([]ports.ThingEntity, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()

	out := make([]ports.ThingEntity, len(r.items))
	copy(out, r.items)
	return out, nil
}

func (r *ThingsRepo) Count(ctx context.Context) (int, error) {
	_ = ctx
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.items), nil
}
