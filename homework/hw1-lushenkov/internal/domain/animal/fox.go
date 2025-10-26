package animal

type Fox struct{ *Predator }

func NewFox(name string, foodPerDayKg int) (*Fox, error) {
	p, err := NewPredator(name, "Fox", foodPerDayKg)
	if err != nil {
		return nil, err
	}
	return &Fox{Predator: p}, nil
}
