package thing

type Feeder struct{ *Thing }

func NewFeeder(name string) (*Feeder, error) {
	base, err := NewThing("Feeder " + name)
	if err != nil {
		return nil, err
	}
	return &Feeder{Thing: base}, nil
}
