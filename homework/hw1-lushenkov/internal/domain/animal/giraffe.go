package animal

type Giraffe struct{ *Herbivore }

func NewGiraffe(name string, foodPerDayKg, kindness int) (*Giraffe, error) {
	h, err := NewHerbivore(name, "Giraffe", foodPerDayKg, kindness)
	if err != nil {
		return nil, err
	}
	return &Giraffe{Herbivore: h}, nil
}
