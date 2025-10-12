package reports

import (
	"context"

	"github.com/perekoshik/zoo/internal/domain/animal"
	"github.com/perekoshik/zoo/internal/domain/ports"
)

type FoodRow struct {
	InventoryNumber int
	DisplayName     string
	KilogramsPerDay int
}

type FoodSummary struct {
	TotalKgPerDay int
	Rows          []FoodRow
}

type FoodService struct {
	animals ports.AnimalRepository
}

func NewFoodService(animals ports.AnimalRepository) *FoodService {
	return &FoodService{animals: animals}
}

func (s *FoodService) Summary(ctx context.Context) (FoodSummary, error) {
	list, err := s.animals.List(ctx)
	if err != nil {
		return FoodSummary{}, err
	}

	rows := make([]FoodRow, 0, len(list))
	total := 0
	for _, a := range list {
		alive := animal.Alive(a)
		kg := alive.FoodPerDayKg()
		rows = append(rows, FoodRow{
			InventoryNumber: a.InventoryNumber(),
			DisplayName:     a.DisplayName(),
			KilogramsPerDay: kg,
		})
		total += kg
	}
	return FoodSummary{TotalKgPerDay: total, Rows: rows}, nil
}
