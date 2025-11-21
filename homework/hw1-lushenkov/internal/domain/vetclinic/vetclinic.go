package vetclinic

import (
	"context"

	"github.com/perekoshik/zoo/internal/domain/animal"
)

type Decision struct {
	Healthy bool
	Reason  string
}

type Clinic interface {
	CheckHealth(ctx context.Context, a animal.Alive) (Decision, error)
}
