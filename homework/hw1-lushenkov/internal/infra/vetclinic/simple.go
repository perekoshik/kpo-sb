package vetclinicimpl

import (
	"context"
	"fmt"

	"github.com/perekoshik/zoo/internal/app/config"
	"github.com/perekoshik/zoo/internal/domain/animal"
	domainvet "github.com/perekoshik/zoo/internal/domain/vetclinic"
)

type SimpleClinic struct {
	cfg config.Config
}

func NewSimpleClinic(cfg config.Config) *SimpleClinic {
	return &SimpleClinic{cfg: cfg}
}

func (c *SimpleClinic) CheckHealth(ctx context.Context, a animal.Alive) (domainvet.Decision, error) {
	_ = ctx

	if !c.cfg.AdmissionLimitEnabled() {
		return domainvet.Decision{Healthy: true}, nil
	}

	limit := c.cfg.MaxFoodPerDayForAdmission
	if a.FoodPerDayKg() <= limit {
		return domainvet.Decision{Healthy: true}, nil
	}

	return domainvet.Decision{
		Healthy: false,
		Reason:  fmt.Sprintf("daily food %d kg exceeds limit %d kg", a.FoodPerDayKg(), limit),
	}, nil
}
