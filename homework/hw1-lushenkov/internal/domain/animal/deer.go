package animal

type Deer struct{ *Herbivore }

func NewDeer(name string, foodPerDayKg, kindness int) (*Deer, error) {
	h, err := NewHerbivore(name, "Deer", foodPerDayKg, kindness)
	if err != nil {
		return nil, err
	}
	return &Deer{Herbivore: h}, nil
}
