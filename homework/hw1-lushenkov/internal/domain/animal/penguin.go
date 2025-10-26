package animal

type Penguin struct{ *Predator }

func NewPenguin(name string, foodPerDayKg int) (*Penguin, error) {
	p, err := NewPredator(name, "Penguin", foodPerDayKg)
	if err != nil {
		return nil, err
	}
	return &Penguin{Predator: p}, nil
}
