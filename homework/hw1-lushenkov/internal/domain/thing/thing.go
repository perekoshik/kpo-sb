package thing

import (
	"fmt"
)

type Thing struct {
	name         string
	inventoryNum int
}

func NewThing(name string) (*Thing, error) {
	if name == "" {
		return nil, fmt.Errorf("name must not be empty")
	}
	return &Thing{name: name, inventoryNum: 0}, nil
}

func (t *Thing) AssignInventoryNumber(n int) error {
	if n < 1 {
		return fmt.Errorf("inventory number must be >= 1")
	}
	if t.inventoryNum != 0 {
		return fmt.Errorf("inventory number already set (%d)", t.inventoryNum)
	}
	t.inventoryNum = n
	return nil
}

func (t *Thing) Name() string         { return t.name }
func (t *Thing) InventoryNumber() int { return t.inventoryNum }
func (t *Thing) DisplayName() string  { return t.name }
