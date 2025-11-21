package services

import (
	"fmt"

	"github.com/perekoshik/oop-go-concepts/internal/domain/factory"
	"github.com/perekoshik/oop-go-concepts/internal/domain/models"
	"github.com/perekoshik/oop-go-concepts/internal/services/repository"
)

// AccountFacade encapsulates operations over bank accounts.
type AccountFacade struct {
	repo    repository.Repository
	factory *factory.DomainFactory
}

// NewAccountFacade builds account facade with dependencies.
func NewAccountFacade(repo repository.Repository, factory *factory.DomainFactory) *AccountFacade {
	return &AccountFacade{repo: repo, factory: factory}
}

// CreateAccount creates account with validation.
func (f *AccountFacade) CreateAccount(name string, initialBalance float64) (*models.BankAccount, error) {
	account, err := f.factory.CreateBankAccount(name, initialBalance)
	if err != nil {
		return nil, err
	}
	if err := f.repo.CreateAccount(*account); err != nil {
		return nil, err
	}
	return account, nil
}

// UpdateAccountName changes human-readable name of account.
func (f *AccountFacade) UpdateAccountName(id, newName string) error {
	account, err := f.repo.GetAccount(id)
	if err != nil {
		return err
	}
	if newName == "" {
		return fmt.Errorf("account name cannot be empty")
	}
	account.Name = newName
	return f.repo.UpdateAccount(*account)
}

// DeleteAccount removes account when there are no dependent operations.
func (f *AccountFacade) DeleteAccount(id string) error {
	operations, err := f.repo.ListOperations()
	if err != nil {
		return err
	}
	for _, op := range operations {
		if op.BankAccountID == id {
			return fmt.Errorf("cannot delete account with existing operations")
		}
	}
	return f.repo.DeleteAccount(id)
}

// ListAccounts returns all accounts.
func (f *AccountFacade) ListAccounts() ([]models.BankAccount, error) {
	return f.repo.ListAccounts()
}

// RecalculateBalance recomputes balance from operations.
func (f *AccountFacade) RecalculateBalance(id string) error {
	account, err := f.repo.GetAccount(id)
	if err != nil {
		return err
	}
	operations, err := f.repo.ListOperations()
	if err != nil {
		return err
	}
	var balance float64
	for _, op := range operations {
		if op.BankAccountID != id {
			continue
		}
		switch op.Type {
		case models.OperationTypeIncome:
			balance += op.Amount
		case models.OperationTypeExpense:
			balance -= op.Amount
		}
	}
	account.Balance = balance
	return f.repo.UpdateAccount(*account)
}
