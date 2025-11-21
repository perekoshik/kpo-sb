package command

import (
	"time"

	"github.com/perekoshik/oop-go-concepts/internal/domain/models"
	"github.com/perekoshik/oop-go-concepts/internal/services"
)

// EditOperationCommand updates operation fields.
type EditOperationCommand struct {
	Facade      *services.OperationFacade
	ID          string
	Type        models.OperationType
	AccountID   string
	CategoryID  string
	Amount      float64
	Date        time.Time
	Description string
}

// Execute performs update.
func (c *EditOperationCommand) Execute() error {
	return c.Facade.UpdateOperation(c.ID, c.Type, c.AccountID, c.CategoryID, c.Amount, c.Date, c.Description)
}
