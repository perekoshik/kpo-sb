package di

import (
	"github.com/perekoshik/zoo/internal/app/config"
	"github.com/perekoshik/zoo/internal/app/onboarding"
	"github.com/perekoshik/zoo/internal/app/reports"
	"github.com/perekoshik/zoo/internal/domain/ports"
	domainvet "github.com/perekoshik/zoo/internal/domain/vetclinic"
	"github.com/perekoshik/zoo/internal/infra/cli"
	"github.com/perekoshik/zoo/internal/infra/repos/mem"
	vetclinicimpl "github.com/perekoshik/zoo/internal/infra/vetclinic"
	"go.uber.org/dig"
)

func Build() (*dig.Container, error) {
	c := dig.New()

	if err := c.Provide(config.Default); err != nil {
		return nil, err
	}

	if err := c.Provide(mem.NewAnimalsRepo, dig.As(new(ports.AnimalRepository))); err != nil {
		return nil, err
	}
	if err := c.Provide(mem.NewThingsRepo, dig.As(new(ports.ThingRepository))); err != nil {
		return nil, err
	}
	if err := c.Provide(mem.NewInventoryQueryMem, dig.As(new(ports.InventoryQuery))); err != nil {
		return nil, err
	}

	if err := c.Provide(vetclinicimpl.NewSimpleClinic, dig.As(new(domainvet.Clinic))); err != nil {
		return nil, err
	}

	if err := c.Provide(onboarding.NewService); err != nil {
		return nil, err
	}
	if err := c.Provide(reports.NewFoodService); err != nil {
		return nil, err
	}
	if err := c.Provide(reports.NewInteractiveService); err != nil {
		return nil, err
	}
	if err := c.Provide(reports.NewBalanceService); err != nil {
		return nil, err
	}

	if err := c.Provide(cli.NewRouter); err != nil {
		return nil, err
	}

	return c, nil
}
