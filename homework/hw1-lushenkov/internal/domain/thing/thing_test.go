package thing

import "testing"

func TestNewThing_Validates(t *testing.T) {
	if _, err := NewThing(""); err == nil {
		t.Fatal("expected error on empty name")
	}
	tb, err := NewThing("Table TN-5")
	if err != nil || tb == nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestThing_AssignInventoryNumber(t *testing.T) {
	tb, _ := NewThing("Table TN-5")
	if err := tb.AssignInventoryNumber(0); err == nil {
		t.Fatal("expected error on <1")
	}
	if err := tb.AssignInventoryNumber(7); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if tb.InventoryNumber() != 7 {
		t.Fatalf("want 7, got %d", tb.InventoryNumber())
	}
	if err := tb.AssignInventoryNumber(8); err == nil {
		t.Fatal("expected error on re-assign")
	}
}
