package color

import "testing"

func TestColorCache(t *testing.T) {
	v := colorCache{
		cache: make(colorMap),
	}

	vNew := v.value(FgWhite)
	if vNew == nil {
		t.Fatal("should not have nil")
	}

	vCache := v.value(FgWhite)
	if vNew != vCache {
		t.Fatalf("should point at same value")
	}

	// create two colors, with same attributes in different order, they should point to same color
	c1 := New(FgRed, BgWhite, Underline)
	c2 := New(Underline, FgRed, BgWhite)
	if c1 != c2 {
		t.Fatal("expect c2 to be same as c1")
	}
}
