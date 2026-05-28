package system

import "testing"

func TestResolveTargetIndex_EmptyReturnsLast(t *testing.T) {
	layers := Layers()
	got := resolveTargetIndex(layers, "")
	want := len(layers) - 1
	if got != want {
		t.Fatalf("resolveTargetIndex(\"\") = %d, want %d", got, want)
	}
}

func TestResolveTargetIndex_L0MinimalFound(t *testing.T) {
	list := Layers()
	got := resolveTargetIndex(list, "L0-minimal")
	if got < 0 {
		t.Fatalf("expected L0-minimal layer to be found, got %d", got)
	}
	if list[got].Name() != "L0-minimal" {
		t.Fatalf("resolved layer name = %q, want L0-minimal", list[got].Name())
	}
}

func TestResolveTargetIndex_UnknownReturnsNegative(t *testing.T) {
	layers := Layers()
	if got := resolveTargetIndex(layers, "nonexistent"); got != -1 {
		t.Fatalf("resolveTargetIndex(\"nonexistent\") = %d, want -1", got)
	}
}
