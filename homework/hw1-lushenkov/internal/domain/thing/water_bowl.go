package thing

type WaterBowl struct{ *Thing }

func NewWaterBowl(name string) (*WaterBowl, error) {
	base, err := NewThing("WaterBowl " + name)
	if err != nil {
		return nil, err
	}
	return &WaterBowl{Thing: base}, nil
}
