package animal

type Tiger struct{ *Predator }

func NewTiger(name string, foodPerDayKg int) (*Tiger, error) {
	p, err := NewPredator(name, "Tiger", foodPerDayKg)
	if err != nil {
		return nil, err
	}
	return &Tiger{Predator: p}, nil
}
