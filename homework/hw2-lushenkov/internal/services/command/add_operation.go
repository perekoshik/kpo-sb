package command

import (
	"time"

	"github.com/perekoshik/oop-go-concepts/internal/domain/models"
	"github.com/perekoshik/oop-go-concepts/internal/services"
)

// AddOperationCommand runs scenario of attaching new operation to account.
type AddOperationCommand struct {
	Facade      *services.OperationFacade
	Type        models.OperationType
	AccountID   string
	CategoryID  string
	Amount      float64
	Date        time.Time
	Description string
	OnCreated   func(*models.Operation)
}

// Execute invokes operation facade to create operation.
func (c *AddOperationCommand) Execute() error {
	operation, err := c.Facade.AddOperation(c.Type, c.AccountID, c.CategoryID, c.Amount, c.Date, c.Description)
	if err != nil {
		return err
	}
	if c.OnCreated != nil {
		c.OnCreated(operation)
	}
	return nil
}
