package repository

import (
	"fmt"
	"sync"

	"github.com/perekoshik/oop-go-concepts/internal/domain/models"
)

// ProxyRepository keeps in-memory cache while delegating persistence to wrapped repository.
type ProxyRepository struct {
	base Repository

	loadOnce sync.Once
	loadErr  error

	mu         sync.RWMutex
	accounts   map[string]models.BankAccount
	categories map[string]models.Category
	operations map[string]models.Operation
}

// NewProxyRepository constructs proxy with lazy cache initialisation.
func NewProxyRepository(base Repository) *ProxyRepository {
	return &ProxyRepository{base: base}
}

func (p *ProxyRepository) ensureLoaded() error {
	p.loadOnce.Do(func() {
		accounts, err := p.base.ListAccounts()
		if err != nil {
			p.loadErr = fmt.Errorf("load accounts: %w", err)
			return
		}
		categories, err := p.base.ListCategories()
		if err != nil {
			p.loadErr = fmt.Errorf("load categories: %w", err)
			return
		}
		operations, err := p.base.ListOperations()
		if err != nil {
			p.loadErr = fmt.Errorf("load operations: %w", err)
			return
		}
		p.mu.Lock()
		defer p.mu.Unlock()
		p.accounts = make(map[string]models.BankAccount)
		for _, acc := range accounts {
			p.accounts[acc.ID] = acc
		}
		p.categories = make(map[string]models.Category)
		for _, cat := range categories {
			p.categories[cat.ID] = cat
		}
		p.operations = make(map[string]models.Operation)
		for _, op := range operations {
			p.operations[op.ID] = op
		}
	})
	return p.loadErr
}

func (p *ProxyRepository) ListAccounts() ([]models.BankAccount, error) {
	if err := p.ensureLoaded(); err != nil {
		return nil, err
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	result := make([]models.BankAccount, 0, len(p.accounts))
	for _, acc := range p.accounts {
		result = append(result, acc)
	}
	return result, nil
}

func (p *ProxyRepository) GetAccount(id string) (*models.BankAccount, error) {
	if err := p.ensureLoaded(); err != nil {
		return nil, err
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	acc, ok := p.accounts[id]
	if !ok {
		return nil, fmt.Errorf("account %s not found", id)
	}
	copy := acc
	return &copy, nil
}

func (p *ProxyRepository) CreateAccount(account models.BankAccount) error {
	if err := p.base.CreateAccount(account); err != nil {
		return err
	}
	if err := p.ensureLoaded(); err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.accounts[account.ID] = account
	return nil
}

func (p *ProxyRepository) UpdateAccount(account models.BankAccount) error {
	if err := p.base.UpdateAccount(account); err != nil {
		return err
	}
	if err := p.ensureLoaded(); err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.accounts[account.ID] = account
	return nil
}

func (p *ProxyRepository) DeleteAccount(id string) error {
	if err := p.base.DeleteAccount(id); err != nil {
		return err
	}
	if err := p.ensureLoaded(); err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.accounts, id)
	return nil
}

func (p *ProxyRepository) ListCategories() ([]models.Category, error) {
	if err := p.ensureLoaded(); err != nil {
		return nil, err
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	result := make([]models.Category, 0, len(p.categories))
	for _, cat := range p.categories {
		result = append(result, cat)
	}
	return result, nil
}

func (p *ProxyRepository) GetCategory(id string) (*models.Category, error) {
	if err := p.ensureLoaded(); err != nil {
		return nil, err
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	cat, ok := p.categories[id]
	if !ok {
		return nil, fmt.Errorf("category %s not found", id)
	}
	copy := cat
	return &copy, nil
}

func (p *ProxyRepository) CreateCategory(category models.Category) error {
	if err := p.base.CreateCategory(category); err != nil {
		return err
	}
	if err := p.ensureLoaded(); err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.categories[category.ID] = category
	return nil
}

func (p *ProxyRepository) UpdateCategory(category models.Category) error {
	if err := p.base.UpdateCategory(category); err != nil {
		return err
	}
	if err := p.ensureLoaded(); err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.categories[category.ID] = category
	return nil
}

func (p *ProxyRepository) DeleteCategory(id string) error {
	if err := p.base.DeleteCategory(id); err != nil {
		return err
	}
	if err := p.ensureLoaded(); err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.categories, id)
	return nil
}

func (p *ProxyRepository) ListOperations() ([]models.Operation, error) {
	if err := p.ensureLoaded(); err != nil {
		return nil, err
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	result := make([]models.Operation, 0, len(p.operations))
	for _, op := range p.operations {
		result = append(result, op)
	}
	return result, nil
}

func (p *ProxyRepository) GetOperation(id string) (*models.Operation, error) {
	if err := p.ensureLoaded(); err != nil {
		return nil, err
	}
	p.mu.RLock()
	defer p.mu.RUnlock()
	op, ok := p.operations[id]
	if !ok {
		return nil, fmt.Errorf("operation %s not found", id)
	}
	copy := op
	return &copy, nil
}

func (p *ProxyRepository) CreateOperation(operation models.Operation) error {
	if err := p.base.CreateOperation(operation); err != nil {
		return err
	}
	if err := p.ensureLoaded(); err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.operations[operation.ID] = operation
	return nil
}

func (p *ProxyRepository) UpdateOperation(operation models.Operation) error {
	if err := p.base.UpdateOperation(operation); err != nil {
		return err
	}
	if err := p.ensureLoaded(); err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.operations[operation.ID] = operation
	return nil
}

func (p *ProxyRepository) DeleteOperation(id string) error {
	if err := p.base.DeleteOperation(id); err != nil {
		return err
	}
	if err := p.ensureLoaded(); err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.operations, id)
	return nil
}

// ReplaceAll swaps cache and persistence with provided snapshot.
func (p *ProxyRepository) ReplaceAll(data StorageData) error {
	if err := p.base.ReplaceAll(data); err != nil {
		return err
	}
	if err := p.ensureLoaded(); err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.accounts = make(map[string]models.BankAccount, len(data.Accounts))
	for _, acc := range data.Accounts {
		p.accounts[acc.ID] = acc
	}
	p.categories = make(map[string]models.Category, len(data.Categories))
	for _, cat := range data.Categories {
		p.categories[cat.ID] = cat
	}
	p.operations = make(map[string]models.Operation, len(data.Operations))
	for _, op := range data.Operations {
		p.operations[op.ID] = op
	}
	return nil
}
