package env

import (
	"testing"
)

func TestDiff_Added(t *testing.T) {
	existing := map[string]string{"A": "1"}
	incoming := map[string]string{"A": "1", "B": "2"}

	d := Diff(existing, incoming)

	if len(d.Added) != 1 || d.Added["B"] != "2" {
		t.Errorf("expected B=2 in Added, got %v", d.Added)
	}
	if len(d.Removed) != 0 {
		t.Errorf("expected no removals, got %v", d.Removed)
	}
	if len(d.Changed) != 0 {
		t.Errorf("expected no changes, got %v", d.Changed)
	}
}

func TestDiff_Removed(t *testing.T) {
	existing := map[string]string{"A": "1", "B": "2"}
	incoming := map[string]string{"A": "1"}

	d := Diff(existing, incoming)

	if len(d.Removed) != 1 || d.Removed["B"] != "2" {
		t.Errorf("expected B=2 in Removed, got %v", d.Removed)
	}
	if len(d.Added) != 0 || len(d.Changed) != 0 {
		t.Errorf("unexpected added/changed entries")
	}
}

func TestDiff_Changed(t *testing.T) {
	existing := map[string]string{"A": "old"}
	incoming := map[string]string{"A": "new"}

	d := Diff(existing, incoming)

	if len(d.Changed) != 1 || d.Changed["A"] != "new" {
		t.Errorf("expected A=new in Changed, got %v", d.Changed)
	}
	if len(d.Added) != 0 || len(d.Removed) != 0 {
		t.Errorf("unexpected added/removed entries")
	}
}

func TestDiff_IsEmpty(t *testing.T) {
	existing := map[string]string{"A": "1"}
	incoming := map[string]string{"A": "1"}

	d := Diff(existing, incoming)
	if !d.IsEmpty() {
		t.Error("expected diff to be empty")
	}
}

func TestDiff_NilExisting(t *testing.T) {
	incoming := map[string]string{"X": "y", "Z": "w"}

	d := Diff(nil, incoming)

	if len(d.Added) != 2 {
		t.Errorf("expected 2 added keys, got %d", len(d.Added))
	}
	if len(d.Removed) != 0 || len(d.Changed) != 0 {
		t.Error("expected no removed or changed keys")
	}
}

func TestDiff_EmptyBoth(t *testing.T) {
	d := Diff(map[string]string{}, map[string]string{})
	if !d.IsEmpty() {
		t.Error("expected empty diff for two empty maps")
	}
}
