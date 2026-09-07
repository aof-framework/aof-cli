package model

import "testing"

func TestNormalizeProfile(t *testing.T) {
	for in, want := range map[string]string{"core": ProfileCore, "AOF-Governed": ProfileGoverned, "assured": ProfileAssured} {
		if got := NormalizeProfile(in); got != want {
			t.Fatalf("%q: got %q want %q", in, got, want)
		}
	}
	if got := NormalizeProfile("invalid"); got != "" {
		t.Fatalf("expected empty, got %q", got)
	}
}
