package mem

import (
	"context"
	"sync"
	"testing"

	"github.com/perekoshik/zoo/internal/domain/animal"
)

type testPred struct{ *animal.Predator }

func newTestPred(name string, food int) *testPred {
	p, _ := animal.NewPredator(name, "TestPred", food)
	return &testPred{Predator: p}
}

func TestAnimalsRepo_ParallelAdd_UniqueNumbers(t *testing.T) {
	repo := NewAnimalsRepo()
	const n = 200

	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			a := newTestPred("A", 1+i%3)
			if _, err := repo.Add(context.Background(), a); err != nil {
				t.Errorf("add failed: %v", err)
			}
		}()
	}
	wg.Wait()

	list, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("list err: %v", err)
	}
	if len(list) != n {
		t.Fatalf("expected %d items, got %d", n, len(list))
	}

	seen := make(map[int]struct{}, n)
	for _, a := range list {
		num := a.InventoryNumber()
		if num < 1 {
			t.Fatalf("invalid inv number %d", num)
		}
		if _, ok := seen[num]; ok {
			t.Fatalf("duplicate inv number %d", num)
		}
		seen[num] = struct{}{}
	}
}
