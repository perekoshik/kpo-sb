package thing

type Lamp struct{ *Thing }

func NewLamp(name string) (*Lamp, error) {
	base, err := NewThing("Lamp " + name)
	if err != nil {
		return nil, err
	}
	return &Lamp{Thing: base}, nil
}
