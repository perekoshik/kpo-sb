package animal

type Monkey struct{ *Herbivore }

func NewMonkey(name string, foodPerDayKg, kindness int) (*Monkey, error) {
	h, err := NewHerbivore(name, "Monkey", foodPerDayKg, kindness)
	if err != nil {
		return nil, err
	}
	return &Monkey{Herbivore: h}, nil
}
