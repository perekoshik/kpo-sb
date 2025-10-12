package thing

type Table struct{ *Thing }

func NewTable(name string) (*Table, error) {
	base, err := NewThing("Table " + name)
	if err != nil {
		return nil, err
	}
	return &Table{Thing: base}, nil
}
