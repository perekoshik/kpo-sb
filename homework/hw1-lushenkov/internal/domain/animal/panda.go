package animal

type Panda struct{ *Herbivore }

func NewPanda(name string, foodPerDayKg, kindness int) (*Panda, error) {
	h, err := NewHerbivore(name, "Panda", foodPerDayKg, kindness)
	if err != nil {
		return nil, err
	}
	return &Panda{Herbivore: h}, nil
}
