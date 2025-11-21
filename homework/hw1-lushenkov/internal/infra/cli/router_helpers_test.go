package cli

import "testing"

func TestFieldsKeepQuoted(t *testing.T) {
	in := `add thing table "TN 5" `
	got := fieldsKeepQuoted(in)
	want := []string{"add", "thing", "table", "TN 5"}
	if len(got) != len(want) {
		t.Fatalf("len: want %d got %d (%v)", len(want), len(got), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("idx %d: want %q got %q", i, want[i], got[i])
		}
	}
}

func TestCanonicalSpecies(t *testing.T) {
	if canonicalSpecies("  rAbBit ") != "Rabbit" {
		t.Fatal("expected Rabbit")
	}
}

func TestIsHerbivoreSpecies(t *testing.T) {
	if !isHerbivoreSpecies("rabbit") {
		t.Fatal("rabbit should be herbivore")
	}
	if isHerbivoreSpecies("tiger") {
		t.Fatal("tiger should not be herbivore")
	}
}
