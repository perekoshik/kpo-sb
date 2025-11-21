package animal

type Bear struct{ *Predator }

func NewBear(name string, foodPerDayKg int) (*Bear, error) {
	p, err := NewPredator(name, "Bear", foodPerDayKg)
	if err != nil {
		return nil, err
	}
	return &Bear{Predator: p}, nil
}
