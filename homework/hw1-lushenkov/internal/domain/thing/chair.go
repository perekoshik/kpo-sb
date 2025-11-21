package thing

type Chair struct{ *Thing }

func NewChair(name string) (*Chair, error) {
	base, err := NewThing("Chair " + name)
	if err != nil {
		return nil, err
	}
	return &Chair{Thing: base}, nil
}
