package animal

type Lion struct{ *Predator }

func NewLion(name string, foodPerDayKg int) (*Lion, error) {
	p, err := NewPredator(name, "Lion", foodPerDayKg)
	if err != nil {
		return nil, err
	}
	return &Lion{Predator: p}, nil
}
