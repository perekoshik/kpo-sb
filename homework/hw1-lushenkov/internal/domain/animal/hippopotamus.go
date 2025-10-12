package animal

type Hippopotamus struct{ *Herbivore }

func NewHippopotamus(name string, foodPerDayKg, kindness int) (*Hippopotamus, error) {
	h, err := NewHerbivore(name, "Hippopotamus", foodPerDayKg, kindness)
	if err != nil {
		return nil, err
	}
	return &Hippopotamus{Herbivore: h}, nil
}
