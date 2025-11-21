package onboarding

import (
	"context"
	"fmt"

	"github.com/perekoshik/zoo/internal/domain/ports"
	"github.com/perekoshik/zoo/internal/domain/vetclinic"
)

type Service struct {
	clinic vetclinic.Clinic
	repo   ports.AnimalRepository
}

func NewService(clinic vetclinic.Clinic, repo ports.AnimalRepository) *Service {
	return &Service{clinic: clinic, repo: repo}
}

type Result struct {
	Accepted bool
	Reason   string
	Animal   ports.AnimalEntity
}

func (s *Service) Admit(ctx context.Context, a ports.AnimalEntity) (Result, error) {
	dec, err := s.clinic.CheckHealth(ctx, a)
	if err != nil {
		return Result{}, fmt.Errorf("vet check failed: %w", err)
	}
	if !dec.Healthy {
		return Result{Accepted: false, Reason: dec.Reason}, nil
	}
	stored, err := s.repo.Add(ctx, a)
	if err != nil {
		return Result{}, fmt.Errorf("add animal failed: %w", err)
	}
	return Result{Accepted: true, Animal: stored}, nil
}
