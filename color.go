package color

import (
	"fmt"
	"io"
	"strings"
)

const (
	escape     = "\x1b["
	endCode    = "m"
	lineFeed   = "\n"
	delimiter  = ";"
	colorReset = "\x1b[0m"
)

type writerValuer interface {
	io.Writer
	value(Attribute) *Color
}

func chainSGRCodes(a []Attribute) string {
	if len(a) == 0 {
		return colorReset
	}
	if len(a) == 1 {
		return escape + a[0].Code() + endCode
	}
	var bld strings.Builder
	bld.Grow((len(a) * 2) + len(escape) + len(endCode))
	bld.WriteString(escape)
	delimsAdded := 0
	for i := 0; i < len(a); i++ {
		if delimsAdded > 0 {
			_, _ = bld.WriteString(delimiter)
		}
		bld.WriteString(a[i].Code())
		delimsAdded++
	}
	bld.WriteString(endCode)
	return bld.String()
}

// Color contains methods to create colored strings of text.
type Color struct {
	out        writerValuer
	colorStart string
}

// NewWithWriter create a Color, supplying a writer, as well as desired attributes.
func NewWithWriter(wtr writerValuer, attrs ...Attribute) *Color {
	return &Color{
		colorStart: chainSGRCodes(attrs),
		out:        wtr,
	}
}

// New creates a Color that outputs to Stdout. It takes a list of Attributes to define
// the appearance of output text produced by Color methods.
func New(attrs ...Attribute) *Color {
	return NewWithWriter(Stdout(), attrs...)
}

// NewStderr helper that constructs a Color with stderr as it's underlying io.Writer
func NewStderr(attrs ...Attribute) *Color {
	return NewWithWriter(Stderr(), attrs...)
}

// Fprint writes decorated text to standard out. The number of bytes written is returned.
func (v Color) Fprint(out writerValuer, a ...interface{}) (int, error) {
	_, _ = fmt.Fprint(out, v.colorStart)
	n, err := fmt.Fprint(out, a...)
	_, _ = fmt.Fprint(out, colorReset)
	return n, err
}

// Fprintf formats according to a format specifier and writes decorated text to standard out. The number of bytes
// written is returned.
func (v Color) Fprintf(out writerValuer, format string, a ...interface{}) (int, error) {
	_, _ = fmt.Fprint(out, v.colorStart)
	n, err := fmt.Fprintf(out, format, a...)
	_, _ = fmt.Fprint(out, colorReset)
	return n, err
}

// Fprintln writes decorated text to standard out with a line feed. The number of bytes
// written is returned.
func (v Color) Fprintln(out writerValuer, a ...interface{}) (int, error) {
	_, _ = fmt.Fprint(out, v.colorStart)
	n, err := fmt.Fprint(out, a...)
	_, _ = fmt.Fprintln(out, colorReset)
	return n, err
}

// Print writes decorated text to the io.Writer passed to the Color constructor. The number of bytes written
// is returned.
func (v Color) Print(a ...interface{}) (int, error) {
	return v.Fprint(v.out, a...)
}

// Printf formats according to a format specifier and writes decorated text to the io.Writer passed to the Color
// constructor. The number of bytes written is returned.
func (v Color) Printf(format string, a ...interface{}) (int, error) {
	return v.Fprintf(v.out, format, a...)
}

// Println writes decorated text to the io.Writer passed to the Color constructor.
// The number of bytes written is returned.
func (v Color) Println(a ...interface{}) (int, error) {
	return v.Fprintln(v.out, a...)
}

// Sprint returns text decorated with the display Attributes passed to Color constructor function.
func (v Color) Sprint(a ...interface{}) string {
	return v.wrap(fmt.Sprint(a...))
}

// Sprint formats according to the format specifier and returns text decorated with the display Attributes
// passed to Color constructor function.
func (v Color) Sprintf(format string, a ...interface{}) string {
	return v.wrap(fmt.Sprintf(format, a...))
}

// Sprint returns text decorated with the display Attributes and terminated by a line feed.
func (v Color) Sprintln(a ...interface{}) string {
	s := v.wrap(fmt.Sprint(a...))
	if !strings.HasSuffix(s, lineFeed) {
		s += lineFeed
	}
	return s
}

// FprintFunc returns a function that wraps Fprint.
func (v Color) FprintFunc() func(out writerValuer, a ...interface{}) {
	return func(out writerValuer, a ...interface{}) {
		_, _ = v.Fprint(out, a...)
	}
}

// FprintfFunc returns a function that wraps Fprintf.
func (v Color) FprintfFunc() func(out writerValuer, format string, a ...interface{}) {
	return func(out writerValuer, format string, a ...interface{}) {
		_, _ = v.Fprintf(out, format, a...)
	}
}

// FprintlnFunc returns a function that wraps Fprintln.
func (v Color) FprintlnFunc() func(out writerValuer, a ...interface{}) {
	return func(out writerValuer, a ...interface{}) {
		_, _ = v.Fprintln(out, a...)
	}
}

// PrintFunc returns a wrapper function for Print.
func (v Color) PrintFunc() func(a ...interface{}) {
	return func(a ...interface{}) {
		_, _ = v.Print(a...)
	}
}

// PrintfFunc returns a wrapper function for Printf.
func (v Color) PrintfFunc() func(format string, a ...interface{}) {
	return func(format string, a ...interface{}) {
		_, _ = v.Printf(format, a...)
	}
}

// PrintlnFunc returns a wrapper function for Println.
func (v Color) PrintlnFunc() func(a ...interface{}) {
	return func(a ...interface{}) {
		_, _ = v.Println(a...)
	}
}

// SprintFunc returns function that wraps Sprint.
func (v Color) SprintFunc() func(a ...interface{}) string {
	return func(a ...interface{}) string {
		return v.Sprint(a...)
	}
}

// SprintfFunc returns function that wraps Sprintf.
func (v Color) SprintfFunc() func(format string, a ...interface{}) string {
	return func(format string, a ...interface{}) string {
		return v.Sprintf(format, a...)
	}
}

// SprintlnFunc returns function that wraps Sprintln.
func (v Color) SprintlnFunc() func(a ...interface{}) string {
	return func(a ...interface{}) string {
		return v.Sprintln(a...)
	}
}

func (v Color) wrap(s string) string {
	var b strings.Builder
	b.Grow(len(v.colorStart) + len(s) + len(colorReset))
	b.WriteString(v.colorStart)
	b.WriteString(s)
	b.WriteString(colorReset)
	return b.String()
}

func colorPrint(out writerValuer, format string, attr Attribute, a ...interface{}) {
	v := out.value(attr)
	if !strings.HasSuffix(format, lineFeed) {
		_, _ = fmt.Fprint(out, v.Sprintf(format, a...)+lineFeed)
		return
	}
	_, _ = v.Fprintf(out, format, a...)
}

var colorCache *valueCache = newValueCache(Stdout())

func colorString(format string, attr Attribute, a ...interface{}) string {
	return colorCache.value(attr).Sprintf(format, a...)
}

// Black helper to produce black text to stdout.
func Black(format string, a ...interface{}) { colorPrint(Stdout(), format, FgBlack, a...) }

// BlackE helper to produce black text to stderr.
func BlackE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgBlack, a...) }

// Red helper to produce red text to stdout.
func Red(format string, a ...interface{}) { colorPrint(Stdout(), format, FgRed, a...) }

// RedE helper to produce red text to stderr.
func RedE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgRed, a...) }

// Green helper to produce green text to stdout.
func Green(format string, a ...interface{}) { colorPrint(Stdout(), format, FgGreen, a...) }

// GreenE helper to produce green text to stderr.
func GreenE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgGreen, a...) }

// Yellow helper to produce yellow text to stdout.
func Yellow(format string, a ...interface{}) { colorPrint(Stdout(), format, FgYellow, a...) }

// YellowE helper to produce yellow text to stderr.
func YellowE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgYellow, a...) }

// Blue helper to produce blue text to stdout.
func Blue(format string, a ...interface{}) { colorPrint(Stdout(), format, FgBlue, a...) }

// BlueE helper to produce blue text to stderr.
func BlueE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgBlue, a...) }

// Magenta helper to produce magenta text to stdout.
func Magenta(format string, a ...interface{}) { colorPrint(Stdout(), format, FgMagenta, a...) }

// MagentaE produces magenta text to stderr.
func MagentaE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgMagenta, a...) }

// Cyan helper to produce cyan text to stdout.
func Cyan(format string, a ...interface{}) { colorPrint(Stdout(), format, FgCyan, a...) }

// CyanE helper to produce cyan text to stderr.
func CyanE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgCyan, a...) }

// White helper to produce white text to stdout.
func White(format string, a ...interface{}) { colorPrint(Stdout(), format, FgWhite, a...) }

// WhiteE helper to produce white text to stderr.
func WhiteE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgWhite, a...) }

// BlackString returns a string decorated with black attributes.
func BlackString(format string, a ...interface{}) string { return colorString(format, FgBlack, a...) }

// RedString returns a string decorated with red attributes.
func RedString(format string, a ...interface{}) string { return colorString(format, FgRed, a...) }

// GreenString returns a string decorated with green attributes.
func GreenString(format string, a ...interface{}) string { return colorString(format, FgGreen, a...) }

// YellowString returns a string decorated with yellow attributes.
func YellowString(format string, a ...interface{}) string { return colorString(format, FgYellow, a...) }

// BlueString returns a string decorated with blue attributes.
func BlueString(format string, a ...interface{}) string { return colorString(format, FgBlue, a...) }

// MagentaString returns a string decorated with magenta attributes.
func MagentaString(format string, a ...interface{}) string {
	return colorString(format, FgMagenta, a...)
}

// CyanString returns a string decorated with cyan attributes.
func CyanString(format string, a ...interface{}) string { return colorString(format, FgCyan, a...) }

// WhiteString returns a string decorated with white attributes.
func WhiteString(format string, a ...interface{}) string { return colorString(format, FgWhite, a...) }

// HiBlack helper to produce black text to stdout.
func HiBlack(format string, a ...interface{}) { colorPrint(Stdout(), format, FgHiBlack, a...) }

// HiBlackE helper to produce black text to stderr.
func HiBlackE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgHiBlack, a...) }

// HiRed helper to write high contrast red text to stdout.
func HiRed(format string, a ...interface{}) { colorPrint(Stdout(), format, FgHiRed, a...) }

// HiRedE helper to write high contrast red text to stderr.
func HiRedE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgHiRed, a...) }

// HiGreen helper writes high contrast green text to stdout.
func HiGreen(format string, a ...interface{}) { colorPrint(Stdout(), format, FgHiGreen, a...) }

// HiGreenE helper writes high contrast green text to stderr.
func HiGreenE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgHiGreen, a...) }

// HiYellow helper writes high contrast yellow text to stdout.
func HiYellow(format string, a ...interface{}) { colorPrint(Stdout(), format, FgHiYellow, a...) }

// HiYellowE helper writes high contrast yellow text to stderr.
func HiYellowE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgHiYellow, a...) }

// HiBlue helper writes high contrast blue text to stdout.
func HiBlue(format string, a ...interface{}) { colorPrint(Stdout(), format, FgHiBlue, a...) }

// HiBlueE helper writes high contrast blue text to stderr.
func HiBlueE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgHiBlue, a...) }

// HiMagenta writes high contrast magenta text to stdout.
func HiMagenta(format string, a ...interface{}) { colorPrint(Stdout(), format, FgHiMagenta, a...) }

// HiMagentaE writes high contrast magenta text to stderr.
func HiMagentaE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgHiMagenta, a...) }

// HiCyan writes high contrast cyan colored text to stdout.
func HiCyan(format string, a ...interface{}) { colorPrint(Stdout(), format, FgHiCyan, a...) }

// HiCyanE writes high contrast contrast cyan colored text to stderr.
func HiCyanE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgHiCyan, a...) }

// HiWhite writes high contrast white colored text to stdout.
func HiWhite(format string, a ...interface{}) { colorPrint(Stdout(), format, FgHiWhite, a...) }

// HiWhiteE writes high contrast white colored text to stderr.
func HiWhiteE(format string, a ...interface{}) { colorPrint(Stderr(), format, FgHiWhite, a...) }

// HiBlackString returns a high contrast black string.
func HiBlackString(format string, a ...interface{}) string {
	return colorString(format, FgHiBlack, a...)
}

// HiRedString returns a high contrast contrast black string.
func HiRedString(format string, a ...interface{}) string { return colorString(format, FgHiRed, a...) }

// HiGreenString returns a high contrast green string.
func HiGreenString(format string, a ...interface{}) string {
	return colorString(format, FgHiGreen, a...)
}

// HiYellowString returns a high contrast yellow string.
func HiYellowString(format string, a ...interface{}) string {
	return colorString(format, FgHiYellow, a...)
}

// HiBlueString returns a high contrast blue string.
func HiBlueString(format string, a ...interface{}) string { return colorString(format, FgHiBlue, a...) }

// HiMagentaString returns a high contrast magenta string.
func HiMagentaString(format string, a ...interface{}) string {
	return colorString(format, FgHiMagenta, a...)
}

// HiCyanString returns a high contrast cyan string.
func HiCyanString(format string, a ...interface{}) string { return colorString(format, FgHiCyan, a...) }

// HiWhiteString returns a high contrast white string.
func HiWhiteString(format string, a ...interface{}) string {
	return colorString(format, FgHiWhite, a...)
}
