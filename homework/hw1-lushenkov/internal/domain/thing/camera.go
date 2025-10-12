package thing

type Camera struct{ *Thing }

func NewCamera(name string) (*Camera, error) {
	base, err := NewThing("Camera " + name)
	if err != nil {
		return nil, err
	}
	return &Camera{Thing: base}, nil
}
