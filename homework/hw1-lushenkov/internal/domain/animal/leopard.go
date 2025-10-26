package animal

type Leopard struct{ *Predator }

func NewLeopard(name string, foodPerDayKg int) (*Leopard, error) {
	p, err := NewPredator(name, "Leopard", foodPerDayKg)
	if err != nil {
		return nil, err
	}
	return &Leopard{Predator: p}, nil
}
