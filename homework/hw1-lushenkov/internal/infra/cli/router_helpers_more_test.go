package cli

import "testing"

func TestGenericThing_AssignAndDisplay(t *testing.T) {
	tg, err := newGenericThing("Table TN-5")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if tg.DisplayName() != "Table TN-5" {
		t.Fatalf("display mismatch: %s", tg.DisplayName())
	}
	if err := tg.AssignInventoryNumber(9); err != nil {
		t.Fatalf("assign err: %v", err)
	}
	if tg.InventoryNumber() != 9 {
		t.Fatalf("want 9, got %d", tg.InventoryNumber())
	}
	if err := tg.AssignInventoryNumber(10); err == nil {
		t.Fatal("expected error on reassign")
	}
}

func TestFieldsKeepQuoted_Complex(t *testing.T) {
	in := `add animal rabbit "Bunny Jr." 3 8`
	got := fieldsKeepQuoted(in)
	want := []string{"add", "animal", "rabbit", "Bunny Jr.", "3", "8"}
	if len(got) != len(want) {
		t.Fatalf("len mismatch: %v", got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("idx %d want %q got %q", i, want[i], got[i])
		}
	}
}
