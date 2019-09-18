package color

import (
	"fmt"
	"io"
)

type ErrTypeColor string
func(et ErrTypeColor) Error() string { return string(et) }

const (
	escape = "\x1b"

	ErrMissingRequiredAttribute = ErrTypeColor("must provide one or more attributes")
)

type Attributes []Attribute

func(a Attributes) String() string {
	var s string
	for i := 0; i < len(a); i++ {
		if len(s) > 0 {
			s += ";"
		}
		s += a[i].String()
	}
	return s
}

type Value struct {
	params Attributes
	out io.Writer
}


func New(out io.Writer, attrs ...Attribute)(*Value, error) {
	if len(attrs) == 0 {
		return nil, ErrMissingRequiredAttribute
	}
	return &Value{
		params: attrs,
		out: out,
	}, nil
}

func(v *Value) Add(attrs ...Attribute) *Value {
	newAttrs := make(Attributes, 0, len(v.params) + len(attrs))
	newAttrs = append(newAttrs, v.params...)
	newAttrs = append(newAttrs, attrs...)
	newVal, _ := New(v.out, newAttrs...)
	return newVal
}

func (v *Value) Print( a ...interface{}) (int, error) {
	fmt.Fprint(v.out, v.format())
	defer fmt.Fprint(v.out, v.unformat())

	return fmt.Fprint(v.out, a...)
}

func (v *Value) Printf( format string, a ...interface{})(int, error) {
	fmt.Fprint(v.out, v.format())
	defer fmt.Fprint(v.out, v.unformat())

	return fmt.Fprintf(v.out, format, a...)
}

func(v *Value) format() string {
	return fmt.Sprintf("%s[%sm", escape, v.params)
}

func(v *Value) unformat() string {
	return fmt.Sprintf("%s[%sm", escape, Reset)
}

func(v *Value) wrap(s string) string {
	return fmt.Sprintf("%s%s%s", v.format(), s, v.unformat())
}






