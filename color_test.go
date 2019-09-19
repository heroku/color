package color

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"testing"
)

type helperFunc func(fmt string, a ...interface{})
type helperFuncString func(fmt string, a ...interface{}) string

type mockConsole struct {
	bytes.Buffer
	*valueCache
}

func newMockConsole() *mockConsole {
	var mc mockConsole
	vc := &valueCache{
		cache:  make(valueMap),
		parent: &mc,
	}
	mc.valueCache = vc
	return &mc
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
				buff := newMockConsole()
				f := NewWithWriter(buff, tc.code).PrintFunc()
				f(tc.text)
				got := buff.String()
				t.Log(got)
				want := fmt.Sprintf("%s%dm%s%s0m", escape, tc.code, tc.text, escape)
				if got != want {
					t.Logf("got  %q", got)
					t.Logf("want %q", want)
					t.Fatal()
				}
			})

			t.Run(tc.text+" printf", func(t *testing.T) {
				buff := newMockConsole()
				f := NewWithWriter(buff, tc.code).PrintfFunc()
				f("%q", tc.text)
				got := buff.String()
				t.Log(got)
				want := fmt.Sprintf("%s%dm%q%s0m", escape, tc.code, tc.text, escape)
				if got != want {
					t.Logf("got  %q", got)
					t.Logf("want %q", want)
					t.Fatal()
				}
			})

			t.Run(tc.text+" println", func(t *testing.T) {
				buff := newMockConsole()
				NewWithWriter(buff, tc.code).PrintlnFunc()(tc.text)
				got := buff.String()
				t.Log(got)
				want := fmt.Sprintf("%s%dm%s%s0m\n", escape, tc.code, tc.text, escape)
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
		test func(out writerValuer, v *Color)
		want string
	}{
		{
			name: "FprintFunc",
			test: func(out writerValuer, v *Color) {
				v.FprintFunc()(out, "white sprint")
			},
			want: "\x1b[37mwhite sprint\x1b[0m",
		},
		{
			name: "FprintfFunc",
			test: func(out writerValuer, v *Color) {
				v.FprintfFunc()(out, "%q", "white sprintf")
			},
			want: "\x1b[37m\"white sprintf\"\x1b[0m",
		},
		{
			name: "FprintlnFunc",
			test: func(out writerValuer, v *Color) {
				v.FprintlnFunc()(out, "white sprintln")
			},
			want: "\x1b[37mwhite sprintln\x1b[0m\n",
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			out := newMockConsole()
			v := NewWithWriter(newMockConsole(), FgWhite)
			tc.test(out, v)
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
	out := newMockConsole()
	v := NewWithWriter(out, FgWhite, Bold, Underline)
	_, _ = v.Print("bold white")
	want := "\x1b[37;1;4mbold white\x1b[0m"
	got := out.String()
	t.Log(got)
	if got != want {
		t.Fatalf("want %q got %q", want, got)
	}
}

func TestMissingAttribute(t *testing.T) {
	t.Parallel()
	out := newMockConsole()
	v := NewWithWriter(out)
	_, _ = v.Print("no color")
	want := "\x1b[0mno color\x1b[0m"
	got := out.String()
	t.Log(got)
	if got != want {
		t.Fatalf("want %q got %q", want, got)
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
		t.Run(tc.code.String(), func(t *testing.T) {
			want := fmt.Sprintf("\x1b[%dmcolor - %q\x1b[0m", tc.code, tc.code)
			got := tc.test("color - %q", tc.code)
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
		want := fmt.Sprintf("\x1b[%dmcolor - %q\x1b[0m\n", tc.code, tc.code)

		t.Run(tc.code.String()+"_stdout", func(t *testing.T) {
			var buff bytes.Buffer
			cons := Stdout()
			oldWriter := cons.colorable
			cons.colorable = &buff
			defer func() {
				cons.colorable = oldWriter
			}()

			tc.testStdout("color - %q", tc.code)
			t.Log(buff.String())
			assertEqualS(t, want, buff.String())
		})

		t.Run(tc.code.String()+"_stderr", func(t *testing.T) {
			var buff bytes.Buffer
			cons := Stderr()
			oldWriter := cons.colorable
			cons.colorable = &buff
			defer func() {
				cons.colorable = oldWriter
			}()

			tc.testStdErr("color - %q", tc.code)
			assertEqualS(t, want, buff.String())
		})
	}
}

func BenchmarkColorFuncs(b *testing.B) {
	cons := Stdout()
	oldWriter := cons.colorable
	cons.colorable = ioutil.Discard
	defer func() {
		cons.colorable = oldWriter
	}()

	for i := 0; i < b.N; i++ {
		Black("hello from %s", "black")
		Green("hello from %s", "green")
		Red("hello from %q.  i'm %d", "red", 23)
	}
}

func BenchmarkColorFuncsParallel(b *testing.B) {
	cons := Stdout()
	oldWriter := cons.colorable
	cons.colorable = ioutil.Discard
	defer func() {
		cons.colorable = oldWriter
	}()

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

type mockWriter struct {
	*valueCache
	w io.Writer
}

func newMockWriter(w io.Writer) *mockWriter {
	var mc mockWriter
	vc := &valueCache{
		cache:  make(valueMap),
		parent: &mc,
	}
	mc.valueCache = vc
	mc.w = w
	return &mc
}

func (w *mockWriter) Write(b []byte) (int, error) {
	return w.Write(b)
}

func BenchmarkColorStruct(b *testing.B) {
	attrs := []Attribute{
		FgBlack,
		FgRed,
		FgGreen,
		FgYellow,
		FgBlue,
		FgMagenta,
		FgCyan,
		FgWhite,
	}

	for i := 0; i < b.N; i++ {
		w := NewWithWriter(newMockWriter(ioutil.Discard), attrs...)
		for i := 0; i < 100; i++ {
			w.Println("jjjjjjfjfjfjfjfjfjfjfjfjfjfjfjfjfjfjfjfjfjfjfjfjfjfjfjfjfjfjfjfjfjfjfjfjfjfjjf")
		}
	}
}

func ExampleRed() {
	// Print some red text.
	Red("Hello red!")
}

func ExampleColorer_Println() {
	// Output underlined white text to stdout.
	New(FgWhite, Underline).Println("I'm underlined and white!")
}

func ExampleColorer_SprintFunc() {
	// Create functions that add color information
	emphasized := New(FgRed, Bold, Underline).SprintFunc()
	fmt.Println("Wow, this is", emphasized("exciting!"))
}
