package animal

type Predator struct {
	*BaseAnimal
}

func NewPredator(name, species string, foodPerDayKg int) (*Predator, error) {
	base, err := NewBaseAnimal(name, species, DietPredator, foodPerDayKg)
	if err != nil {
		return nil, err
	}
	return &Predator{BaseAnimal: base}, nil
}
