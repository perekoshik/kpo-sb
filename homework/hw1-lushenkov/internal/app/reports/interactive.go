package reports

import (
	"context"

	"github.com/perekoshik/zoo/internal/app/config"
	"github.com/perekoshik/zoo/internal/domain/animal"
	"github.com/perekoshik/zoo/internal/domain/ports"
)

type InteractiveRow struct {
	InventoryNumber int
	DisplayName     string
	Kindness        int
}

type InteractiveService struct {
	cfg     config.Config
	animals ports.AnimalRepository
}

func NewInteractiveService(cfg config.Config, animals ports.AnimalRepository) *InteractiveService {
	return &InteractiveService{cfg: cfg, animals: animals}
}

func (s *InteractiveService) List(ctx context.Context) ([]InteractiveRow, error) {
	list, err := s.animals.List(ctx)
	if err != nil {
		return nil, err
	}

	thr := s.cfg.InteractiveKindnessThreshold
	out := make([]InteractiveRow, 0)
	for _, a := range list {
		if kc, ok := a.(animal.KindnessCarrier); ok {
			k := kc.Kindness0to10()
			if k > thr-1 {
				out = append(out, InteractiveRow{
					InventoryNumber: a.InventoryNumber(),
					DisplayName:     a.DisplayName(),
					Kindness:        k,
				})
			}
		}
	}
	return out, nil
}
