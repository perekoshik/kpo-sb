package animal

import "fmt"

type Diet int

const (
	DietUnknown Diet = iota
	DietHerbivore
	DietPredator
	DietOmnivore
)

func (d Diet) String() string {
	switch d {
	case DietHerbivore:
		return "Herbivore"
	case DietPredator:
		return "Predator"
	case DietOmnivore:
		return "Omnivore"
	default:
		return "Unknown"
	}
}

type BaseAnimal struct {
	name         string
	species      string
	diet         Diet
	foodPerDayKg int
	inventoryNum int
}

func NewBaseAnimal(name, species string, diet Diet, foodPerDayKg int) (*BaseAnimal, error) {
	if name == "" {
		return nil, fmt.Errorf("name must not be empty")
	}
	if species == "" {
		return nil, fmt.Errorf("species must not be empty")
	}
	if diet == DietUnknown {
		return nil, fmt.Errorf("diet must be specified")
	}
	if foodPerDayKg < 0 {
		return nil, fmt.Errorf("foodPerDayKg must be >= 0")
	}
	return &BaseAnimal{
		name:         name,
		species:      species,
		diet:         diet,
		foodPerDayKg: foodPerDayKg,
		inventoryNum: 0,
	}, nil
}

func (a *BaseAnimal) AssignInventoryNumber(n int) error {
	if n < 1 {
		return fmt.Errorf("inventory number must be >= 1")
	}
	if a.inventoryNum != 0 {
		return fmt.Errorf("inventory number already set (%d)", a.inventoryNum)
	}
	a.inventoryNum = n
	return nil
}

func (a *BaseAnimal) Name() string         { return a.name }
func (a *BaseAnimal) Species() string      { return a.species }
func (a *BaseAnimal) Diet() Diet           { return a.diet }
func (a *BaseAnimal) FoodPerDayKg() int    { return a.foodPerDayKg } // Alive
func (a *BaseAnimal) InventoryNumber() int { return a.inventoryNum } // InventoryItem
func (a *BaseAnimal) DisplayName() string  { return fmt.Sprintf("%s '%s'", a.species, a.name) }
