package command

import "github.com/perekoshik/oop-go-concepts/internal/services"

// DeleteOperationCommand removes operation by id.
type DeleteOperationCommand struct {
	Facade *services.OperationFacade
	ID     string
}

// Execute runs delete scenario.
func (c *DeleteOperationCommand) Execute() error {
	return c.Facade.DeleteOperation(c.ID)
}
