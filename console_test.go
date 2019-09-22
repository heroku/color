package color

import (
	"bytes"
	"testing"
)

func TestColorCache(t *testing.T) {
	v := colorCache{
		cache:  make(colorMap),
		parent: newMockConsole(),
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

func TestConsoleWrite(t *testing.T) {
	var out bytes.Buffer
	c := Console{
		current: &out,
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

func TestSet(t *testing.T) {
	var cons Console
	var buff bytes.Buffer
	cons.current = &buff

	col := New(FgRed)
	cons.Set(col)
	if buff.String() != col.colorStart {
		t.Fatalf("want %q got %q", col.colorStart, buff.String())
	}
}

func TestUnset(t *testing.T) {
	var cons Console
	var buff bytes.Buffer
	cons.current = &buff
	cons.Unset()
	if buff.String() != colorReset {
		t.Fatalf("got %q want %q", buff.String(), colorReset)
	}
}
