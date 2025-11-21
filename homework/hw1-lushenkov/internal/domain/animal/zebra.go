package animal

type Zebra struct{ *Herbivore }

func NewZebra(name string, foodPerDayKg, kindness int) (*Zebra, error) {
	h, err := NewHerbivore(name, "Zebra", foodPerDayKg, kindness)
	if err != nil {
		return nil, err
	}
	return &Zebra{Herbivore: h}, nil
}
