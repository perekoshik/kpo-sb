package animal

import "fmt"

type Herbivore struct {
	*BaseAnimal
	kindness int
}

func NewHerbivore(name, species string, foodPerDayKg, kindness int) (*Herbivore, error) {
	if kindness < 0 || kindness > 10 {
		return nil, fmt.Errorf("kindness must be in [0..10]")
	}
	base, err := NewBaseAnimal(name, species, DietHerbivore, foodPerDayKg)
	if err != nil {
		return nil, err
	}
	return &Herbivore{BaseAnimal: base, kindness: kindness}, nil
}

func (h *Herbivore) Kindness0to10() int { return h.kindness }
