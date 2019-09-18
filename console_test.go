package color

import (
	"bytes"
	"testing"
)

func TestValueCache(t *testing.T) {
	v := valueCache{
		cache:   make(valueMap),
		parent:  &bytes.Buffer{},
	}

	vNew := v.value(FgWhite)
	if vNew == nil {
		t.Fatal("should not have nil")
	}

	vCache := v.value(FgWhite)
	if vNew != vCache {
		t.Fatalf("should point at same value")
	}
}
