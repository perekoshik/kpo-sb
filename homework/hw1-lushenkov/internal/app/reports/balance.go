package reports

import (
	"context"
	"sort"

	"github.com/perekoshik/zoo/internal/domain/inventory"
	"github.com/perekoshik/zoo/internal/domain/ports"
)

type BalanceRow struct {
	InventoryNumber int
	DisplayName     string
}

type BalanceService struct {
	q ports.InventoryQuery
}

func NewBalanceService(q ports.InventoryQuery) *BalanceService {
	return &BalanceService{q: q}
}

func (s *BalanceService) List(ctx context.Context) ([]BalanceRow, error) {
	items, err := s.q.ListAllInventories(ctx)
	if err != nil {
		return nil, err
	}
	rows := make([]BalanceRow, 0, len(items))
	for _, it := range items {
		rows = append(rows, BalanceRow{
			InventoryNumber: it.InventoryNumber(),
			DisplayName:     it.DisplayName(),
		})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].InventoryNumber < rows[j].InventoryNumber })
	return rows, nil
}

func toInventoryRows(items []inventory.InventoryItem) []BalanceRow {
	rows := make([]BalanceRow, 0, len(items))
	for _, it := range items {
		rows = append(rows, BalanceRow{
			InventoryNumber: it.InventoryNumber(),
			DisplayName:     it.DisplayName(),
		})
	}
	return rows
}
