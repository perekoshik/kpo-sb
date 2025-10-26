package thing

type Computer struct{ *Thing }

func NewComputer(name string) (*Computer, error) {
	base, err := NewThing("Computer " + name)
	if err != nil {
		return nil, err
	}
	return &Computer{Thing: base}, nil
}
