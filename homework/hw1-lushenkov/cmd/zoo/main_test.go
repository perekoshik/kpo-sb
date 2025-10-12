package main

import (
	"testing"

	"github.com/perekoshik/zoo/internal/infra/di"
)

// smoke-тест: контейнер собирается (не запускаем интерактив).
func TestContainerBuilds(t *testing.T) {
	if _, err := di.Build(); err != nil {
		t.Fatalf("di build err: %v", err)
	}
}
