package animal

type Parrot struct{ *Herbivore }

func NewParrot(name string, foodPerDayKg, kindness int) (*Parrot, error) {
	h, err := NewHerbivore(name, "Parrot", foodPerDayKg, kindness)
	if err != nil {
		return nil, err
	}
	return &Parrot{Herbivore: h}, nil
}
