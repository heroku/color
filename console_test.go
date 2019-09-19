package color

import (
	"bytes"
	"testing"
)

func TestValueCache(t *testing.T) {
	v := valueCache{
		cache:  make(valueMap),
		parent: &bytes.Buffer{},
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

func TestConsoleWrite(t *testing.T) {
	var out bytes.Buffer
	c := Console{
		colorable: &out,
	}
	n, err := c.Write([]byte("foo"))
	if err != nil {
		t.Fatal("no error expected", err)
	}
	if n != len("foo") {
		t.Fatalf("expected len %d got %d", len("foo"), n)
	}
	if out.String() != "foo" {
		t.Fatalf("got %q expected %q", out.String(), "foo")
	}
}
