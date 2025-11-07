package repository

import "github.com/perekoshik/oop-go-concepts/internal/domain/models"

// StorageData represents persisted snapshot of all domain entities.
type StorageData struct {
	Accounts   []models.BankAccount `json:"accounts"`
	Categories []models.Category    `json:"categories"`
	Operations []models.Operation   `json:"operations"`
}

// Repository provides CRUD operations over domain entities.
type Repository interface {
	// Accounts
	ListAccounts() ([]models.BankAccount, error)
	GetAccount(id string) (*models.BankAccount, error)
	CreateAccount(account models.BankAccount) error
	UpdateAccount(account models.BankAccount) error
	DeleteAccount(id string) error

	// Categories
	ListCategories() ([]models.Category, error)
	GetCategory(id string) (*models.Category, error)
	CreateCategory(category models.Category) error
	UpdateCategory(category models.Category) error
	DeleteCategory(id string) error

	// Operations
	ListOperations() ([]models.Operation, error)
	GetOperation(id string) (*models.Operation, error)
	CreateOperation(operation models.Operation) error
	UpdateOperation(operation models.Operation) error
	DeleteOperation(id string) error

	// ReplaceAll atomically swaps data snapshot.
	ReplaceAll(data StorageData) error
}
