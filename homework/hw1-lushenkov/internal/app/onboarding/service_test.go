package onboarding_test

import (
	"context"
	"testing"

	"github.com/perekoshik/zoo/internal/app/onboarding"
	"github.com/perekoshik/zoo/internal/domain/animal"
	domainvet "github.com/perekoshik/zoo/internal/domain/vetclinic"
	"github.com/perekoshik/zoo/internal/infra/repos/mem"
)

type allowAllClinic struct{}

func (allowAllClinic) CheckHealth(ctx context.Context, a animal.Alive) (domainvet.Decision, error) {
	return domainvet.Decision{Healthy: true}, nil
}

type denyAllClinic struct{}

func (denyAllClinic) CheckHealth(ctx context.Context, a animal.Alive) (domainvet.Decision, error) {
	return domainvet.Decision{Healthy: false, Reason: "test-deny"}, nil
}

func TestAdmit_AcceptsAndAssignsNumber(t *testing.T) {
	repo := mem.NewAnimalsRepo()
	svc := onboarding.NewService(allowAllClinic{}, repo)

	rabbit, _ := animal.NewRabbit("Bunny", 3, 7)
	res, err := svc.Admit(context.Background(), rabbit)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Accepted {
		t.Fatalf("expected accepted, got deny: %v", res.Reason)
	}
	if got := res.Animal.InventoryNumber(); got != 1 {
		t.Fatalf("expected inv #1, got %d", got)
	}
}

func TestAdmit_DeniesWithoutConsumingNumber(t *testing.T) {
	repo := mem.NewAnimalsRepo()
	svcDeny := onboarding.NewService(denyAllClinic{}, repo)
	svcAllow := onboarding.NewService(allowAllClinic{}, repo)

	elephant, _ := animal.NewElephant("Dumbo", 100, 9)
	res, err := svcDeny.Admit(context.Background(), elephant)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Accepted {
		t.Fatalf("expected deny")
	}

	rabbit, _ := animal.NewRabbit("Bunny", 3, 7)
	res2, err := svcAllow.Admit(context.Background(), rabbit)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := res2.Animal.InventoryNumber(); got != 1 {
		t.Fatalf("expected inv #1 after denial, got %d", got)
	}
}
