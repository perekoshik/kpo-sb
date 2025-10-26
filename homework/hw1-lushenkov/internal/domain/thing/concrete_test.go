package thing

import "testing"

type invThing interface {
	Name() string
	AssignInventoryNumber(int) error
	InventoryNumber() int
	DisplayName() string
}

func TestConcreteThings(t *testing.T) {
	cases := []struct {
		name   string
		mk     func(string) (invThing, error)
		prefix string
	}{
		{
			name: "Table",
			mk: func(s string) (invThing, error) {
				return NewTable(s)
			},
			prefix: "Table ",
		},
		{
			name: "Computer",
			mk: func(s string) (invThing, error) {
				return NewComputer(s)
			},
			prefix: "Computer ",
		},
		{
			name: "Camera",
			mk: func(s string) (invThing, error) {
				return NewCamera(s)
			},
			prefix: "Camera ",
		},
		{
			name: "Cage",
			mk: func(s string) (invThing, error) {
				return NewCage(s)
			},
			prefix: "Cage ",
		},
		{
			name: "Chair",
			mk: func(s string) (invThing, error) {
				return NewChair(s)
			},
			prefix: "Chair ",
		},
		{
			name: "Feeder",
			mk: func(s string) (invThing, error) {
				return NewFeeder(s)
			},
			prefix: "Feeder ",
		},
		{
			name: "Lamp",
			mk: func(s string) (invThing, error) {
				return NewLamp(s)
			},
			prefix: "Lamp ",
		},
		{
			name: "WaterBowl",
			mk: func(s string) (invThing, error) {
				return NewWaterBowl(s)
			},
			prefix: "WaterBowl ",
		},
	}

	for _, tc := range cases {
		obj, err := tc.mk("X1")
		if err != nil {
			t.Fatalf("%s err: %v", tc.name, err)
		}
		if obj.Name() != tc.prefix+"X1" {
			t.Fatalf("%s name mismatch: %q", tc.name, obj.Name())
		}
		if err := obj.AssignInventoryNumber(42); err != nil {
			t.Fatalf("%s assign err: %v", tc.name, err)
		}
		if obj.InventoryNumber() != 42 {
			t.Fatalf("%s inv mismatch", tc.name)
		}
		if obj.DisplayName() == "" {
			t.Fatalf("%s display empty", tc.name)
		}
	}
}
