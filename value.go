package color

import (
	"fmt"
	"io"
	"strings"
)

type ErrTypeColor string

func (et ErrTypeColor) Error() string { return string(et) }

const (
	escape    = "\x1b"
	lineFeed  = "\n"
	delimiter = ";"

	ErrMissingRequiredAttribute = ErrTypeColor("must provide one or more attributes")
)

type Attributes []Attribute

func (a Attributes) String() string {
	var bld strings.Builder
	bld.Grow(len(a) * 4)
	for i := 0; i < len(a); i++ {
		if bld.Len() > 0 {
			_, _ = bld.WriteString(delimiter)
		}
		bld.WriteString(a[i].String())
	}
	return bld.String()
}

type Value struct {
	//params Attributes
	out                  io.Writer
	colorStart, colorEnd string
}

func New(out io.Writer, attrs ...Attribute) (*Value, error) {
	if len(attrs) == 0 {
		return nil, ErrMissingRequiredAttribute
	}

	return &Value{
		colorStart: fmt.Sprintf("%s[%sm", escape, attrs),
		colorEnd:   fmt.Sprintf("%s[%sm", escape, Reset),
		out:        out,
	}, nil
}

func (v *Value) Fprint(out io.Writer, a ...interface{}) (int, error) {
	_, _ = fmt.Fprint(out, v.colorStart)
	n, err := fmt.Fprint(out, a...)
	_, _ = fmt.Fprint(out, v.colorEnd)
	return n, err
}

func (v *Value) Fprintf(out io.Writer, format string, a ...interface{}) (int, error) {
	_, _ = fmt.Fprint(out, v.colorStart)
	n, err := fmt.Fprintf(out, format, a...)
	_, _ = fmt.Fprint(out, v.colorEnd)
	return n, err
}

func (v *Value) Fprintln(out io.Writer, a ...interface{}) (int, error) {
	_, _ = fmt.Fprint(out, v.colorStart)
	n, err := fmt.Fprintln(out, a...)
	_, _ = fmt.Fprint(out, v.colorEnd)
	return n, err
}

func (v *Value) Print(a ...interface{}) (int, error) {
	_, _ = fmt.Fprint(v.out, v.colorStart)
	n, err := fmt.Fprint(v.out, a...)
	_, _ = fmt.Fprint(v.out, v.colorEnd)
	return n, err
}

func (v *Value) Printf(format string, a ...interface{}) (int, error) {
	_, _ = fmt.Fprint(v.out, v.colorStart)
	n, err := fmt.Fprintf(v.out, format, a...)
	_, _ = fmt.Fprint(v.out, v.colorEnd)
	return n, err
}

func (v *Value) Println(a ...interface{}) (int, error) {
	_, _ = fmt.Fprint(v.out, v.colorStart)
	n, err := fmt.Fprintln(v.out, a...)
	_, _ = fmt.Fprint(v.out, v.colorEnd)
	return n, err
}

func (v *Value) Sprint(a ...interface{}) string {
	return v.wrap(fmt.Sprint(a...))
}

func (v *Value) Sprintf(format string, a ...interface{}) string {
	return v.wrap(fmt.Sprintf(format, a...))
}

func (v *Value) Sprintln(a ...interface{}) string {
	return v.wrap(fmt.Sprintln(a...))
}

func (v *Value) FprintFunc() func(out io.Writer, a ...interface{}) {
	return func(out io.Writer, a ...interface{}) {
		_, _ = v.Fprint(out, a...)
	}
}

func (v *Value) FprintfFunc() func(out io.Writer, format string, a ...interface{}) {
	return func(out io.Writer, format string, a ...interface{}) {
		_, _ = v.Fprintf(out, format, a...)
	}
}

func (v *Value) FprintlnFunc() func(out io.Writer, a ...interface{}) {
	return func(out io.Writer, a ...interface{}) {
		_, _ = v.Fprintln(out, a...)
	}
}

func (v *Value) PrintFunc() func(a ...interface{}) {
	return func(a ...interface{}) {
		_, _ = v.Print(a...)
	}
}

func (v *Value) PrintfFunc() func(format string, a ...interface{}) {
	return func(format string, a ...interface{}) {
		_, _ = v.Printf(format, a...)
	}
}

func (v *Value) PrintlnFunc() func(a ...interface{}) {
	return func(a ...interface{}) {
		_, _ = v.Println(a...)
	}
}

func (v *Value) SprintFunc() func(a ...interface{}) string {
	return func(a ...interface{}) string {
		return v.wrap(fmt.Sprint(a...))
	}
}

func (v *Value) SprintfFunc() func(format string, a ...interface{}) string {
	return func(format string, a ...interface{}) string {
		return v.wrap(fmt.Sprintf(format, a...))
	}
}

func (v *Value) SprintlnFunc() func(a ...interface{}) string {
	return func(a ...interface{}) string {
		return v.wrap(fmt.Sprintln(a...))
	}
}

func (v *Value) wrap(s string) string {
	var b strings.Builder
	b.Grow(len(v.colorStart) + len(s) + len(v.colorEnd))
	b.WriteString(v.colorStart)
	b.WriteString(s)
	b.WriteString(v.colorEnd)
	return b.String()
}

type writerValuer interface {
	io.Writer
	value(Attribute) *Value
}

func colorPrint(out writerValuer, format string, attr Attribute, a ...interface{}) {
	v := out.value(attr)
	if !strings.HasSuffix(format, lineFeed) {
		format += lineFeed
	}
	_, _ = v.Fprintf(out, format, a...)
}

func colorString(format string, attr Attribute, a ...interface{}) string {
	v, _ := New(nil, attr)
	return v.Sprintf(format, a...)
}

func Black(format string, a ...interface{})  { colorPrint(Stdout(), format, FgBlack, a...) }
func BlackE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgBlack, a...) }

func Red(format string, a ...interface{})  { colorPrint(Stdout(), format, FgRed, a...) }
func RedE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgRed, a...) }

func Green(format string, a ...interface{})  { colorPrint(Stdout(), format, FgGreen, a...) }
func GreenE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgGreen, a...) }

func Yellow(format string, a ...interface{})  { colorPrint(Stdout(), format, FgYellow, a...) }
func YellowE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgYellow, a...) }

func Blue(format string, a ...interface{})  { colorPrint(Stdout(), format, FgBlue, a...) }
func BlueE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgBlue, a...) }

func Magenta(format string, a ...interface{})  { colorPrint(Stdout(), format, FgMagenta, a...) }
func MagentaE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgMagenta, a...) }

func Cyan(format string, a ...interface{})  { colorPrint(Stdout(), format, FgCyan, a...) }
func CyanE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgCyan, a...) }

func White(format string, a ...interface{})  { colorPrint(Stdout(), format, FgWhite, a...) }
func WhiteE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgWhite, a...) }

func BlackString(format string, a ...interface{}) string { return colorString(format, FgBlack, a...) }

func RedString(format string, a ...interface{}) string { return colorString(format, FgRed, a...) }

func GreenString(format string, a ...interface{}) string { return colorString(format, FgGreen, a...) }

func YellowString(format string, a ...interface{}) string { return colorString(format, FgYellow, a...) }

func BlueString(format string, a ...interface{}) string { return colorString(format, FgBlue, a...) }

func MagentaString(format string, a ...interface{}) string {
	return colorString(format, FgMagenta, a...)
}

func CyanString(format string, a ...interface{}) string { return colorString(format, FgCyan, a...) }

func WhiteString(format string, a ...interface{}) string { return colorString(format, FgWhite, a...) }

func HiBlack(format string, a ...interface{})  { colorPrint(Stdout(), format, FgHiBlack, a...) }
func HiBlackE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgHiBlack, a...) }

func HiRed(format string, a ...interface{})  { colorPrint(Stdout(), format, FgHiRed, a...) }
func HiRedE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgHiRed, a...) }

func HiGreen(format string, a ...interface{})  { colorPrint(Stdout(), format, FgHiGreen, a...) }
func HiGreenE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgHiGreen, a...) }

func HiYellow(format string, a ...interface{})  { colorPrint(Stdout(), format, FgHiYellow, a...) }
func HiYellowE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgHiYellow, a...) }

func HiBlue(format string, a ...interface{})  { colorPrint(Stdout(), format, FgHiBlue, a...) }
func HiBlueE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgHiBlue, a...) }

func HiMagenta(format string, a ...interface{})  { colorPrint(Stdout(), format, FgHiMagenta, a...) }
func HiMagentaE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgHiMagenta, a...) }

func HiCyan(format string, a ...interface{})  { colorPrint(Stdout(), format, FgHiCyan, a...) }
func HiCyanE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgHiCyan, a...) }

func HiWhite(format string, a ...interface{})  { colorPrint(Stdout(), format, FgHiWhite, a...) }
func HiWhiteE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgHiWhite, a...) }

func HiBlackString(format string, a ...interface{}) string {
	return colorString(format, FgHiBlack, a...)
}

func HiRedString(format string, a ...interface{}) string { return colorString(format, FgHiRed, a...) }

func HiGreenString(format string, a ...interface{}) string {
	return colorString(format, FgHiGreen, a...)
}

func HiYellowString(format string, a ...interface{}) string {
	return colorString(format, FgHiYellow, a...)
}

func HiBlueString(format string, a ...interface{}) string { return colorString(format, FgHiBlue, a...) }

func HiMagentaString(format string, a ...interface{}) string {
	return colorString(format, FgHiMagenta, a...)
}

func HiCyanString(format string, a ...interface{}) string { return colorString(format, FgHiCyan, a...) }

func HiWhiteString(format string, a ...interface{}) string {
	return colorString(format, FgHiWhite, a...)
}
