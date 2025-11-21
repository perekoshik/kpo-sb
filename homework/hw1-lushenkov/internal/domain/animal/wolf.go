package animal

type Wolf struct{ *Predator }

func NewWolf(name string, foodPerDayKg int) (*Wolf, error) {
	p, err := NewPredator(name, "Wolf", foodPerDayKg)
	if err != nil {
		return nil, err
	}
	return &Wolf{Predator: p}, nil
}
