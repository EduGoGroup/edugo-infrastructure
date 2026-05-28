package common

import "testing"

func TestMustParseUUID_Valid(t *testing.T) {
	const s = "11111111-2222-3333-4444-555555555555"
	got := MustParseUUID(s)
	if got.String() != s {
		t.Fatalf("expected %s, got %s", s, got.String())
	}
}

func TestMustParseUUID_PanicsOnInvalid(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on invalid uuid, got none")
		}
	}()
	_ = MustParseUUID("not-a-uuid")
}
