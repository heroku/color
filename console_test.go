package color

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestConsoleWrite(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
	var cons Console
	var buff bytes.Buffer
	cons.current = &buff
	cons.Unset()
	if buff.String() != colorReset {
		t.Fatalf("got %q want %q", buff.String(), colorReset)
	}
}

func mockStd() (*Console, func() string) {
	r, w, _ := os.Pipe()
	console := NewConsole(w)
	return console, func() string {
		_ = w.Close()
		var b bytes.Buffer
		_, _ = io.Copy(&b, r)
		_ = r.Close()
		return b.String()
	}
}

func TestDisable(t *testing.T) {
	// can't be parallel since we're toggling global color which will interfere with concurrently running tests
	Disable(true)
	defer func() {
		Disable(false)
	}()
	cons, outf := mockStd()
	c := New(FgRed)
	_, _ = cons.Print(c, "foo")
	assertEqualS(t, outf(), "foo")
}

func TestEnable(t *testing.T) {
	cons, outf := mockStd()
	c := New(FgRed)
	_, _ = cons.Print(c, "foo")
	_, _ = cons.Println(c, "bar")
	_, _ = cons.Printf(c, "%d", 32)
	want := "\x1b[31mfoo\x1b[0m\x1b[31mbar\x1b[0m\n\x1b[31m32\x1b[0m"
	assertEqualS(t, outf(), want)
}

func TestConsoleEnable(t *testing.T) {
	// can't be parallel since we're toggling global color which will interfere with concurrently running tests
	Disable(true)
	defer func() {
		Disable(false)
	}()
	cons, outf := mockStd()
	c := New(FgRed)
	_, _ = cons.Print(c, "foo")
	cons.DisableColors(false)
	_, _ = cons.Print(c, "bar")
	want := "foo\x1b[31mbar\x1b[0m"
	assertEqualS(t, outf(), want)
}

func TestFileDescriptor(t *testing.T) {
	c := NewConsole(os.Stdout)
	if c.Fd() != os.Stdout.Fd() {
		t.Fatalf("fd mismatch stdout %X console %X", os.Stdout.Fd(), c.Fd())
	}
}
