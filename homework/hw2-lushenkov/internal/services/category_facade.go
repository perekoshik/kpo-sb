package services

import (
	"fmt"

	"github.com/perekoshik/oop-go-concepts/internal/domain/factory"
	"github.com/perekoshik/oop-go-concepts/internal/domain/models"
	"github.com/perekoshik/oop-go-concepts/internal/services/repository"
)

// CategoryFacade coordinates category lifecycle operations.
type CategoryFacade struct {
	repo    repository.Repository
	factory *factory.DomainFactory
}

// NewCategoryFacade wires dependencies for category facade.
func NewCategoryFacade(repo repository.Repository, factory *factory.DomainFactory) *CategoryFacade {
	return &CategoryFacade{repo: repo, factory: factory}
}

// CreateCategory creates and persists new category.
func (f *CategoryFacade) CreateCategory(name string, catType models.CategoryType) (*models.Category, error) {
	category, err := f.factory.CreateCategory(name, catType)
	if err != nil {
		return nil, err
	}
	if err := f.repo.CreateCategory(*category); err != nil {
		return nil, err
	}
	return category, nil
}

// UpdateCategoryName updates category name.
func (f *CategoryFacade) UpdateCategoryName(id, newName string) error {
	if newName == "" {
		return fmt.Errorf("category name cannot be empty")
	}
	category, err := f.repo.GetCategory(id)
	if err != nil {
		return err
	}
	category.Name = newName
	return f.repo.UpdateCategory(*category)
}

// DeleteCategory removes category without operations.
func (f *CategoryFacade) DeleteCategory(id string) error {
	operations, err := f.repo.ListOperations()
	if err != nil {
		return err
	}
	for _, op := range operations {
		if op.CategoryID == id {
			return fmt.Errorf("cannot delete category with linked operations")
		}
	}
	return f.repo.DeleteCategory(id)
}

// ListCategories returns all categories.
func (f *CategoryFacade) ListCategories() ([]models.Category, error) {
	return f.repo.ListCategories()
}
