package di

import (
	"context"
	"testing"

	"github.com/perekoshik/zoo/internal/infra/cli"
)

func TestBuild_ResolveRouter(t *testing.T) {
	c, err := Build()
	if err != nil {
		t.Fatalf("build err: %v", err)
	}
	err = c.Invoke(func(r *cli.Router) {
		_ = r
	})
	if err != nil {
		t.Fatalf("resolve router err: %v", err)
	}

	err = c.Invoke(func(r *cli.Router) {
		_ = context.Background()
	})
	if err != nil {
		t.Fatalf("invoke err: %v", err)
	}
}
