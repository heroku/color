package color

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"sync"
	"testing"
)

type helperFunc func(fmt string, a ...interface{})
type helperFuncString func(fmt string, a ...interface{}) string

type mockConsole struct {
	sync.Mutex
	bytes.Buffer
	*valueCache
}

func newMockConsole() *mockConsole {
	var mock mockConsole
	mock.valueCache = &valueCache{
		cache:  make(valueMap),
		parent: &mock,
	}
	return &mock
}
func (c *mockConsole) Write(b []byte) (int, error) {
	c.Lock()
	defer c.Unlock()
	return c.Buffer.Write(b)
}

func (c *mockConsole) String() string {
	c.Lock()
	defer c.Unlock()
	return c.Buffer.String()
}

func assertEqualS(t *testing.T, want, got string) {
	t.Helper()
	if got != want {
		t.Logf("got  %q", got)
		t.Logf("want %q", want)
		t.Fatal("mismatch")
	}
}

func TestColor(t *testing.T) {
	t.Parallel()
	tt := []struct {
		text string
		code Attribute
	}{
		{"black", FgBlack},
		{text: "red", code: FgRed},
		{text: "green", code: FgGreen},
		{text: "yellow", code: FgYellow},
		{text: "blue", code: FgBlue},
		{text: "magent", code: FgMagenta},
		{text: "cyan", code: FgCyan},
		{text: "white", code: FgWhite},
		{text: "hblack", code: FgHiBlack},
		{text: "hred", code: FgHiRed},
		{text: "hgreen", code: FgHiGreen},
		{text: "hyellow", code: FgHiYellow},
		{text: "hblue", code: FgHiBlue},
		{text: "hmagent", code: FgHiMagenta},
		{text: "hcyan", code: FgHiCyan},
		{text: "hwhite", code: FgHiWhite},
	}

	for _, tc := range tt {
		t.Run(tc.text, func(t *testing.T) {
			t.Run(tc.text+" print", func(t *testing.T) {
				var buff bytes.Buffer
				v, _ := New(&buff, tc.code)
				f := v.PrintFunc()
				f(tc.text)
				got := buff.String()
				t.Log(got)
				want := fmt.Sprintf("%s[%dm%s%s[0m", escape, tc.code, tc.text, escape)
				if got != want {
					t.Logf("got  %q", got)
					t.Logf("want %q", want)
					t.Fatal()
				}
			})

			t.Run(tc.text+" printf", func(t *testing.T) {
				var buff bytes.Buffer
				v, _ := New(&buff, tc.code)
				f := v.PrintfFunc()
				f("%q", tc.text)
				got := buff.String()
				t.Log(got)
				want := fmt.Sprintf("%s[%dm%q%s[0m", escape, tc.code, tc.text, escape)
				if got != want {
					t.Logf("got  %q", got)
					t.Logf("want %q", want)
					t.Fatal()
				}
			})

			t.Run(tc.text+" println", func(t *testing.T) {
				var buff bytes.Buffer
				v, _ := New(&buff, tc.code)
				f := v.PrintlnFunc()
				f(tc.text)
				got := buff.String()
				t.Log(got)
				want := fmt.Sprintf("%s[%dm%s%s[0m\n", escape, tc.code, tc.text, escape)
				if got != want {
					t.Logf("got  %q", got)
					t.Logf("want %q", want)
					t.Fatal()
				}
			})

		})
	}
}

func TestIoFuncs(t *testing.T) {
	t.Parallel()
	tt := []struct {
		name string
		test func(out io.Writer, v *Value)
		want string
	}{
		{
			name: "FprintFunc",
			test: func(out io.Writer, v *Value) {
				v.FprintFunc()(out, "white sprint")
			},
			want: "\x1b[37mwhite sprint\x1b[0m",
		},
		{
			name: "FprintfFunc",
			test: func(out io.Writer, v *Value) {
				v.FprintfFunc()(out, "%q", "white sprintf")
			},
			want: "\x1b[37m\"white sprintf\"\x1b[0m",
		},
		{
			name: "FprintlnFunc",
			test: func(out io.Writer, v *Value) {
				v.FprintlnFunc()(out, "white sprintln")
			},
			want: "\x1b[37mwhite sprintln\x1b[0m\n",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			v, _ := New(&bytes.Buffer{}, FgWhite)
			tc.test(&out, v)
			got := out.String()
			t.Log(got)
			if got != tc.want {
				t.Logf("want %q", tc.want)
				t.Logf("got  %q", got)
				t.Fatal()
			}
		})
	}
}

func TestMultiAttribute(t *testing.T) {
	t.Parallel()
	var out bytes.Buffer
	v, _ := New(&out, FgWhite, Bold, Underline)
	_, _ = v.Print("bold white")
	want := "\x1b[37;1;4mbold white\x1b[0m"
	got := out.String()
	t.Log(got)
	if got != want {
		t.Fatalf("want %q got %q", want, got)
	}
}

func TestValueError(t *testing.T) {
	t.Parallel()
	_, err := New(nil, FgWhite)
	if err != ErrMissingWriter {
		t.Fatal("expected", ErrMissingWriter.String())
	}
	var unused bytes.Buffer
	_, err = New(&unused)
	if err != ErrMissingRequiredAttribute {
		t.Fatal("expected", ErrMissingRequiredAttribute.String())
	}
}

func TestStringHelperFuncs(t *testing.T) {
	t.Parallel()
	tt := []struct {
		code Attribute
		test helperFuncString
	}{
		{FgBlack, BlackString},
		{FgRed, RedString},
		{FgGreen, GreenString},
		{FgYellow, YellowString},
		{FgBlue, BlueString},
		{FgMagenta, MagentaString},
		{FgCyan, CyanString},
		{FgHiBlack, HiBlackString},
		{FgHiRed, HiRedString},
		{FgHiGreen, HiGreenString},
		{FgHiYellow, HiYellowString},
		{FgHiBlue, HiBlueString},
		{FgHiMagenta, HiMagentaString},
		{FgHiCyan, HiCyanString},
		{FgHiWhite, HiWhiteString},
		{FgWhite, WhiteString},
	}

	for _, tc := range tt {
		t.Run(tc.code.Name(), func(t *testing.T) {
			want := fmt.Sprintf("\x1b[%dmcolor - %q\x1b[0m", tc.code, tc.code.Name())
			got := tc.test("color - %q", tc.code.Name())
			assertEqualS(t, want, got)
		})
	}
}

func TestHelperStdoutFuncs(t *testing.T) {
	t.Parallel()
	tt := []struct {
		code       Attribute
		testStdout helperFunc
		testStdErr helperFunc
	}{
		{FgBlack, Black, BlackE},
		{FgRed, Red, RedE},
		{FgGreen, Green, GreenE},
		{FgYellow, Yellow, YellowE},
		{FgBlue, Blue, BlueE},
		{FgMagenta, Magenta, MagentaE},
		{FgCyan, Cyan, CyanE},
		{FgWhite, White, WhiteE},
		{FgHiBlack, HiBlack, HiBlackE},
		{FgHiRed, HiRed, HiRedE},
		{FgHiGreen, HiGreen, HiGreenE},
		{FgHiYellow, HiYellow, HiYellowE},
		{FgHiBlue, HiBlue, HiBlueE},
		{FgHiMagenta, HiMagenta, HiMagentaE},
		{FgHiCyan, HiCyan, HiCyanE},
		{FgHiWhite, HiWhite, HiWhiteE},
	}

	for _, tc := range tt {
		want := fmt.Sprintf("\x1b[%dmcolor - %q\x1b[0m\n", tc.code, tc.code.Name())

		t.Run(tc.code.Name()+"_stdout", func(t *testing.T) {
			mock := newMockConsole()
			oldOut := printOut
			printOut = mock
			defer func() {
				printOut = oldOut
			}()

			tc.testStdout("color - %q", tc.code.Name())
			t.Log(mock.String())
			assertEqualS(t, want, mock.String())
		})

		t.Run(tc.code.Name()+"_stderr", func(t *testing.T) {
			mock := newMockConsole()
			oldOut := printErr
			printErr = mock
			defer func() {
				printErr = oldOut
			}()

			tc.testStdErr("color - %q", tc.code.Name())
			assertEqualS(t, want, mock.String())
		})
	}
}

type benchConsole struct {
	sync.Mutex
	*valueCache
	io.Writer
}

func newBenchConsole() *benchConsole {
	var mock benchConsole
	mock.valueCache = &valueCache{
		cache:  make(valueMap),
		parent: &mock,
	}
	mock.Writer = ioutil.Discard
	return &mock
}
func (c *benchConsole) Write(b []byte) (int, error) {
	c.Lock()
	defer c.Unlock()
	return c.Writer.Write(b)
}

func BenchmarkColorFuncs(b *testing.B) {
	oldCons := printOut
	printOut = newBenchConsole()
	defer func() { printOut = oldCons }()

	for i := 0; i < b.N; i++ {
		Black("hello from %s", "black")
		Green("hello from %s", "green")
		Red("hello from %q.  i'm %d", "red", 23)
	}
}

func BenchmarkColorFuncsParallel(b *testing.B) {
	oldCons := printOut
	printOut = newBenchConsole()
	defer func() { printOut = oldCons }()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Black("hello from %s xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", "black")
			Green("hello from %s yyyyyhfhhdhehehehehhskdkdkdkdkdkdkdkkkekkekekekeekekkk", "green")
			Red("hello from %q.  i'm %d", "red", 23)
			Black("more blach stuff")
			Green("hello from %s xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx", "green")
			Red("hello from %q.  i'm %d", "red", 23)
		}
	})
}
