package animal

type Elephant struct{ *Herbivore }

func NewElephant(name string, foodPerDayKg, kindness int) (*Elephant, error) {
	h, err := NewHerbivore(name, "Elephant", foodPerDayKg, kindness)
	if err != nil {
		return nil, err
	}
	return &Elephant{Herbivore: h}, nil
}
