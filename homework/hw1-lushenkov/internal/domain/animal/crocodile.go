package animal

type Crocodile struct{ *Predator }

func NewCrocodile(name string, foodPerDayKg int) (*Crocodile, error) {
	p, err := NewPredator(name, "Crocodile", foodPerDayKg)
	if err != nil {
		return nil, err
	}
	return &Crocodile{Predator: p}, nil
}
