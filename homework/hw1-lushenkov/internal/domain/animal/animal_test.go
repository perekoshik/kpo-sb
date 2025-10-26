package animal

import "testing"

func TestNewBaseAnimal_Validates(t *testing.T) {
	if _, err := NewBaseAnimal("", "Rabbit", DietHerbivore, 1); err == nil {
		t.Fatal("expected error on empty name")
	}
	if _, err := NewBaseAnimal("Bunny", "", DietHerbivore, 1); err == nil {
		t.Fatal("expected error on empty species")
	}
	if _, err := NewBaseAnimal("Bunny", "Rabbit", DietUnknown, 1); err == nil {
		t.Fatal("expected error on unknown diet")
	}
	if _, err := NewBaseAnimal("Bunny", "Rabbit", DietHerbivore, -1); err == nil {
		t.Fatal("expected error on negative food")
	}
	if a, err := NewBaseAnimal("Bunny", "Rabbit", DietHerbivore, 2); err != nil || a == nil {
		t.Fatalf("expected ok, got err=%v", err)
	}
}

func TestAssignInventoryNumber_Once(t *testing.T) {
	a, _ := NewBaseAnimal("Bunny", "Rabbit", DietHerbivore, 2)
	if err := a.AssignInventoryNumber(0); err == nil {
		t.Fatal("expected error on <1")
	}
	if err := a.AssignInventoryNumber(42); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got := a.InventoryNumber(); got != 42 {
		t.Fatalf("want 42, got %d", got)
	}
	if err := a.AssignInventoryNumber(43); err == nil {
		t.Fatal("expected error on re-assign")
	}
}

func TestDiet_String(t *testing.T) {
	if DietUnknown.String() != "Unknown" ||
		DietHerbivore.String() != "Herbivore" ||
		DietPredator.String() != "Predator" ||
		DietOmnivore.String() != "Omnivore" {
		t.Fatal("unexpected Diet.String values")
	}
}
