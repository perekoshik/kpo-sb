package animal

type Rabbit struct{ *Herbivore }

func NewRabbit(name string, foodPerDayKg, kindness int) (*Rabbit, error) {
	h, err := NewHerbivore(name, "Rabbit", foodPerDayKg, kindness)
	if err != nil {
		return nil, err
	}
	return &Rabbit{Herbivore: h}, nil
}
