package services

import (
	"fmt"
	"time"

	"github.com/perekoshik/oop-go-concepts/internal/domain/factory"
	"github.com/perekoshik/oop-go-concepts/internal/domain/models"
	"github.com/perekoshik/oop-go-concepts/internal/services/repository"
)

// OperationFacade orchestrates operation lifecycle and balance updates.
type OperationFacade struct {
	repo    repository.Repository
	factory *factory.DomainFactory
}

// NewOperationFacade initialises operation facade.
func NewOperationFacade(repo repository.Repository, factory *factory.DomainFactory) *OperationFacade {
	return &OperationFacade{repo: repo, factory: factory}
}

// AddOperation creates operation and updates account balance accordingly.
func (f *OperationFacade) AddOperation(opType models.OperationType, accountID, categoryID string, amount float64, date time.Time, description string) (*models.Operation, error) {
	if err := f.ensureCategoryMatchesOperation(categoryID, opType); err != nil {
		return nil, err
	}
	if _, err := f.repo.GetAccount(accountID); err != nil {
		return nil, err
	}
	operation, err := f.factory.CreateOperation(opType, accountID, categoryID, amount, date, description)
	if err != nil {
		return nil, err
	}
	if err := f.repo.CreateOperation(*operation); err != nil {
		return nil, err
	}
	if err := f.applyOperationEffect(operation, true); err != nil {
		return nil, err
	}
	return operation, nil
}

// UpdateOperation edits operation fields and re-applies balance impact.
func (f *OperationFacade) UpdateOperation(id string, opType models.OperationType, accountID, categoryID string, amount float64, date time.Time, description string) error {
	prev, err := f.repo.GetOperation(id)
	if err != nil {
		return err
	}
	if err := f.ensureCategoryMatchesOperation(categoryID, opType); err != nil {
		return err
	}
	if _, err := f.repo.GetAccount(accountID); err != nil {
		return err
	}

	if err := f.applyOperationEffect(prev, false); err != nil {
		return err
	}

	updated, err := f.factory.CreateOperation(opType, accountID, categoryID, amount, date, description)
	if err != nil {
		return err
	}
	updated.ID = id
	if err := f.repo.UpdateOperation(*updated); err != nil {
		return err
	}
	return f.applyOperationEffect(updated, true)
}

// DeleteOperation removes operation and reverts account balance impact.
func (f *OperationFacade) DeleteOperation(id string) error {
	operation, err := f.repo.GetOperation(id)
	if err != nil {
		return err
	}
	if err := f.repo.DeleteOperation(id); err != nil {
		return err
	}
	return f.applyOperationEffect(operation, false)
}

// ListOperations returns all operations.
func (f *OperationFacade) ListOperations() ([]models.Operation, error) {
	return f.repo.ListOperations()
}

func (f *OperationFacade) ensureCategoryMatchesOperation(categoryID string, opType models.OperationType) error {
	category, err := f.repo.GetCategory(categoryID)
	if err != nil {
		return err
	}
	switch category.Type {
	case models.CategoryTypeIncome:
		if opType != models.OperationTypeIncome {
			return fmt.Errorf("category %s is income, but operation is %s", category.Name, opType)
		}
	case models.CategoryTypeExpense:
		if opType != models.OperationTypeExpense {
			return fmt.Errorf("category %s is expense, but operation is %s", category.Name, opType)
		}
	default:
		return fmt.Errorf("unsupported category type %s", category.Type)
	}
	return nil
}

func (f *OperationFacade) applyOperationEffect(operation *models.Operation, apply bool) error {
	account, err := f.repo.GetAccount(operation.BankAccountID)
	if err != nil {
		return err
	}
	coef := 1.0
	if !apply {
		coef = -1.0
	}
	switch operation.Type {
	case models.OperationTypeIncome:
		account.Balance += coef * operation.Amount
	case models.OperationTypeExpense:
		account.Balance -= coef * operation.Amount
	default:
		return fmt.Errorf("unknown operation type %s", operation.Type)
	}
	return f.repo.UpdateAccount(*account)
}
