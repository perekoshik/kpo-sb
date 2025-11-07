package factory

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/perekoshik/oop-go-concepts/internal/domain/models"
)

// IDGenerator abstracts id generation strategy for domain objects.
type IDGenerator interface {
	NewID() (string, error)
}

// UUIDGenerator produces UUIDv4 identifiers.
type UUIDGenerator struct{}

// NewID returns fresh UUID string.
func (UUIDGenerator) NewID() (string, error) {
	return newUUID()
}

// DomainFactory centralises creation and validation of domain entities.
type DomainFactory struct {
	idGen IDGenerator
}

// NewDomainFactory constructs factory with provided ID generator.
func NewDomainFactory(idGen IDGenerator) *DomainFactory {
	return &DomainFactory{idGen: idGen}
}

// CreateBankAccount validates and returns a new BankAccount aggregate.
func (f *DomainFactory) CreateBankAccount(name string, initialBalance float64) (*models.BankAccount, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("account name is required")
	}
	id, err := f.idGen.NewID()
	if err != nil {
		return nil, fmt.Errorf("generate account id: %w", err)
	}
	if initialBalance < 0 {
		return nil, errors.New("initial balance cannot be negative")
	}
	return &models.BankAccount{ID: id, Name: name, Balance: initialBalance}, nil
}

// CreateCategory validates and returns a new Category entity.
func (f *DomainFactory) CreateCategory(name string, categoryType models.CategoryType) (*models.Category, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("category name is required")
	}
	if categoryType != models.CategoryTypeIncome && categoryType != models.CategoryTypeExpense {
		return nil, errors.New("invalid category type")
	}
	id, err := f.idGen.NewID()
	if err != nil {
		return nil, fmt.Errorf("generate category id: %w", err)
	}
	return &models.Category{ID: id, Type: categoryType, Name: name}, nil
}

// CreateOperation validates inputs and returns a new Operation entity.
func (f *DomainFactory) CreateOperation(opType models.OperationType, accountID, categoryID string, amount float64, date time.Time, description string) (*models.Operation, error) {
	if opType != models.OperationTypeIncome && opType != models.OperationTypeExpense {
		return nil, errors.New("invalid operation type")
	}
	if strings.TrimSpace(accountID) == "" {
		return nil, errors.New("bank account id is required")
	}
	if strings.TrimSpace(categoryID) == "" {
		return nil, errors.New("category id is required")
	}
	if amount <= 0 {
		return nil, errors.New("operation amount must be positive")
	}
	if date.IsZero() {
		return nil, errors.New("operation date is required")
	}
	id, err := f.idGen.NewID()
	if err != nil {
		return nil, fmt.Errorf("generate operation id: %w", err)
	}
	return &models.Operation{
		ID:            id,
		Type:          opType,
		BankAccountID: accountID,
		CategoryID:    categoryID,
		Amount:        amount,
		Date:          date,
		Description:   strings.TrimSpace(description),
	}, nil
}

func newUUID() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", fmt.Errorf("read random: %w", err)
	}
	// Set version (4) and variant bits according to RFC 4122.
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	dst := make([]byte, 36)
	hex.Encode(dst[0:8], b[0:4])
	dst[8] = '-'
	hex.Encode(dst[9:13], b[4:6])
	dst[13] = '-'
	hex.Encode(dst[14:18], b[6:8])
	dst[18] = '-'
	hex.Encode(dst[19:23], b[8:10])
	dst[23] = '-'
	hex.Encode(dst[24:36], b[10:16])
	return string(dst), nil
}
