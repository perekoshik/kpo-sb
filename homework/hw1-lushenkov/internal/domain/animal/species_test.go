package animal

import "testing"

func TestHerbivoreInterfaces(t *testing.T) {
	h, err := NewHerbivore("Bunny", "Rabbit", 3, 7)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if h.FoodPerDayKg() != 3 {
		t.Fatalf("want 3 kg/day")
	}
	if h.Kindness0to10() != 7 {
		t.Fatalf("want kindness 7")
	}
	if err := h.AssignInventoryNumber(5); err != nil {
		t.Fatalf("assign err: %v", err)
	}
	if h.InventoryNumber() != 5 {
		t.Fatalf("want inv 5")
	}
}

func TestPredatorInterfaces(t *testing.T) {
	p, err := NewPredator("Hunter", "Tiger", 25)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if p.FoodPerDayKg() != 25 {
		t.Fatalf("want 25")
	}
	if _, ok := any(p).(KindnessCarrier); ok {
		t.Fatal("predator must not implement KindnessCarrier")
	}
}

func TestSpeciesConstructors(t *testing.T) {
	ele, _ := NewElephant("Dumbo", 120, 9)
	if ele.Species() != "Elephant" {
		t.Fatalf("species mismatch: %s", ele.Species())
	}
	lion, _ := NewLion("King", 30)
	if lion.Species() != "Lion" {
		t.Fatalf("species mismatch: %s", lion.Species())
	}
}
