package animal

import "testing"

func TestAllSpeciesConstructors(t *testing.T) {
	hcases := []struct {
		name    string
		mk      func(string, int, int) (interface{ Species() string }, error)
		wantSpc string
	}{
		{"Elephant", func(n string, f, k int) (interface{ Species() string }, error) { return NewElephant(n, f, k) }, "Elephant"},
		{"Giraffe", func(n string, f, k int) (interface{ Species() string }, error) { return NewGiraffe(n, f, k) }, "Giraffe"},
		{"Zebra", func(n string, f, k int) (interface{ Species() string }, error) { return NewZebra(n, f, k) }, "Zebra"},
		{"Panda", func(n string, f, k int) (interface{ Species() string }, error) { return NewPanda(n, f, k) }, "Panda"},
		{"Parrot", func(n string, f, k int) (interface{ Species() string }, error) { return NewParrot(n, f, k) }, "Parrot"},
		{"Deer", func(n string, f, k int) (interface{ Species() string }, error) { return NewDeer(n, f, k) }, "Deer"},
		{"Hippopotamus", func(n string, f, k int) (interface{ Species() string }, error) { return NewHippopotamus(n, f, k) }, "Hippopotamus"},
		{"Rabbit", func(n string, f, k int) (interface{ Species() string }, error) {
			a, err := NewRabbit(n, f, k)
			if err != nil {
				return nil, err
			}
			return a, nil
		}, "Rabbit"},
		{"Monkey", func(n string, f, k int) (interface{ Species() string }, error) {
			a, err := NewMonkey(n, f, k)
			if err != nil {
				return nil, err
			}
			return a, nil
		}, "Monkey"},
	}
	for _, tc := range hcases {
		a, err := tc.mk("X", 3, 7)
		if err != nil {
			t.Fatalf("%s err: %v", tc.name, err)
		}
		if a.Species() != tc.wantSpc {
			t.Fatalf("%s species: %s", tc.name, a.Species())
		}
	}

	pcases := []struct {
		name    string
		mk      func(string, int) (interface{ Species() string }, error)
		wantSpc string
	}{
		{"Lion", func(n string, f int) (interface{ Species() string }, error) { return NewLion(n, f) }, "Lion"},
		{"Leopard", func(n string, f int) (interface{ Species() string }, error) { return NewLeopard(n, f) }, "Leopard"},
		{"Bear", func(n string, f int) (interface{ Species() string }, error) { return NewBear(n, f) }, "Bear"},
		{"Crocodile", func(n string, f int) (interface{ Species() string }, error) { return NewCrocodile(n, f) }, "Crocodile"},
		{"Penguin", func(n string, f int) (interface{ Species() string }, error) { return NewPenguin(n, f) }, "Penguin"},
		{"Fox", func(n string, f int) (interface{ Species() string }, error) { return NewFox(n, f) }, "Fox"},
		{"Tiger", func(n string, f int) (interface{ Species() string }, error) {
			a, err := NewTiger(n, f)
			if err != nil {
				return nil, err
			}
			return a, nil
		}, "Tiger"},
		{"Wolf", func(n string, f int) (interface{ Species() string }, error) {
			a, err := NewWolf(n, f)
			if err != nil {
				return nil, err
			}
			return a, nil
		}, "Wolf"},
	}
	for _, tc := range pcases {
		a, err := tc.mk("Y", 10)
		if err != nil {
			t.Fatalf("%s err: %v", tc.name, err)
		}
		if a.Species() != tc.wantSpc {
			t.Fatalf("%s species: %s", tc.name, a.Species())
		}
	}
}
