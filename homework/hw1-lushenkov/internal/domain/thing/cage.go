package thing

type Cage struct{ *Thing }

func NewCage(name string) (*Cage, error) {
	base, err := NewThing("Cage " + name)
	if err != nil {
		return nil, err
	}
	return &Cage{Thing: base}, nil
}
