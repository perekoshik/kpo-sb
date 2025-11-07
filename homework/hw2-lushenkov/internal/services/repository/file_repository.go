package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"sync"

	"github.com/perekoshik/oop-go-concepts/internal/domain/models"
)

// FileRepository persists data in JSON file. It is intended to be wrapped by caching proxy.
type FileRepository struct {
	path string
	mu   sync.Mutex
}

// NewFileRepository configures repository with file path.
func NewFileRepository(path string) *FileRepository {
	return &FileRepository{path: path}
}

func (r *FileRepository) ListAccounts() ([]models.BankAccount, error) {
	data, err := r.load()
	if err != nil {
		return nil, err
	}
	return append([]models.BankAccount(nil), data.Accounts...), nil
}

func (r *FileRepository) GetAccount(id string) (*models.BankAccount, error) {
	data, err := r.load()
	if err != nil {
		return nil, err
	}
	for _, acc := range data.Accounts {
		if acc.ID == id {
			copy := acc
			return &copy, nil
		}
	}
	return nil, fmt.Errorf("account %s not found", id)
}

func (r *FileRepository) CreateAccount(account models.BankAccount) error {
	return r.withData(func(data *StorageData) error {
		data.Accounts = append(data.Accounts, account)
		return nil
	})
}

func (r *FileRepository) UpdateAccount(account models.BankAccount) error {
	return r.withData(func(data *StorageData) error {
		for i, acc := range data.Accounts {
			if acc.ID == account.ID {
				data.Accounts[i] = account
				return nil
			}
		}
		return fmt.Errorf("account %s not found", account.ID)
	})
}

func (r *FileRepository) DeleteAccount(id string) error {
	return r.withData(func(data *StorageData) error {
		for i, acc := range data.Accounts {
			if acc.ID == id {
				data.Accounts = append(data.Accounts[:i], data.Accounts[i+1:]...)
				return nil
			}
		}
		return fmt.Errorf("account %s not found", id)
	})
}

func (r *FileRepository) ListCategories() ([]models.Category, error) {
	data, err := r.load()
	if err != nil {
		return nil, err
	}
	return append([]models.Category(nil), data.Categories...), nil
}

func (r *FileRepository) GetCategory(id string) (*models.Category, error) {
	data, err := r.load()
	if err != nil {
		return nil, err
	}
	for _, cat := range data.Categories {
		if cat.ID == id {
			copy := cat
			return &copy, nil
		}
	}
	return nil, fmt.Errorf("category %s not found", id)
}

func (r *FileRepository) CreateCategory(category models.Category) error {
	return r.withData(func(data *StorageData) error {
		data.Categories = append(data.Categories, category)
		return nil
	})
}

func (r *FileRepository) UpdateCategory(category models.Category) error {
	return r.withData(func(data *StorageData) error {
		for i, cat := range data.Categories {
			if cat.ID == category.ID {
				data.Categories[i] = category
				return nil
			}
		}
		return fmt.Errorf("category %s not found", category.ID)
	})
}

func (r *FileRepository) DeleteCategory(id string) error {
	return r.withData(func(data *StorageData) error {
		for i, cat := range data.Categories {
			if cat.ID == id {
				data.Categories = append(data.Categories[:i], data.Categories[i+1:]...)
				return nil
			}
		}
		return fmt.Errorf("category %s not found", id)
	})
}

func (r *FileRepository) ListOperations() ([]models.Operation, error) {
	data, err := r.load()
	if err != nil {
		return nil, err
	}
	return append([]models.Operation(nil), data.Operations...), nil
}

func (r *FileRepository) GetOperation(id string) (*models.Operation, error) {
	data, err := r.load()
	if err != nil {
		return nil, err
	}
	for _, op := range data.Operations {
		if op.ID == id {
			copy := op
			return &copy, nil
		}
	}
	return nil, fmt.Errorf("operation %s not found", id)
}

func (r *FileRepository) CreateOperation(operation models.Operation) error {
	return r.withData(func(data *StorageData) error {
		data.Operations = append(data.Operations, operation)
		return nil
	})
}

func (r *FileRepository) UpdateOperation(operation models.Operation) error {
	return r.withData(func(data *StorageData) error {
		for i, op := range data.Operations {
			if op.ID == operation.ID {
				data.Operations[i] = operation
				return nil
			}
		}
		return fmt.Errorf("operation %s not found", operation.ID)
	})
}

func (r *FileRepository) DeleteOperation(id string) error {
	return r.withData(func(data *StorageData) error {
		for i, op := range data.Operations {
			if op.ID == id {
				data.Operations = append(data.Operations[:i], data.Operations[i+1:]...)
				return nil
			}
		}
		return fmt.Errorf("operation %s not found", id)
	})
}

// ReplaceAll overwrites storage with provided snapshot.
func (r *FileRepository) ReplaceAll(data StorageData) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.saveUnsafe(data)
}

func (r *FileRepository) withData(fn func(data *StorageData) error) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	data, err := r.loadUnsafe()
	if err != nil {
		return err
	}

	if err := fn(&data); err != nil {
		return err
	}

	return r.saveUnsafe(data)
}

func (r *FileRepository) load() (StorageData, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.loadUnsafe()
}

func (r *FileRepository) loadUnsafe() (StorageData, error) {
	file, err := os.Open(r.path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return StorageData{}, nil
		}
		return StorageData{}, fmt.Errorf("open storage file: %w", err)
	}
	defer file.Close()

	var data StorageData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		if errors.Is(err, io.EOF) {
			return StorageData{}, nil
		}
		return StorageData{}, fmt.Errorf("decode storage: %w", err)
	}
	return data, nil
}

func (r *FileRepository) saveUnsafe(data StorageData) error {
	file, err := os.Create(r.path)
	if err != nil {
		return fmt.Errorf("create storage file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("encode storage: %w", err)
	}
	return nil
}
