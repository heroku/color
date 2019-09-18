package color

import (
	"fmt"
	"strconv"
)

// Attribute defines a single SGR Code
type Attribute int

func (a Attribute) String() string {
	return strconv.Itoa(int(a))
}

func(a Attribute) Name() string {
	m := map[Attribute]string {
		Reset: "Reset",
		Bold: "Bold",
		Faint: "Faint",
		Italic: "Italic",
		Underline: "Underline",
		BlinkSlow: "BlinkSlow",
		BlinkRapid: "BlinkRapid",
		ReverseVideo: "ReverseVideo",
		Concealed: "Concealed",
		CrossedOut: "CrossedOut",
		FgBlack: "FgBlack",
		FgRed: "FgRed",
		FgGreen: "FgGreen",
		FgYellow: "FgYellow",
		FgBlue: "FgBlue",
		FgMagenta: "FgMagenta",
		FgCyan: "FgCyan",
		FgWhite: "FgWhite",
		FgHiBlack: "FgHiBlack",
		FgHiRed: "FgHiRed",
		FgHiGreen: "FgHiGreen",
		FgHiYellow: "FgHiYellow",
		FgHiBlue: "FgHiBlue",
		FgHiMagenta: "FgHiMagenta",
		FgHiCyan: "FgHiCyan",
		FgHiWhite: "FgHiWhite",
		BgBlack: "BgBlack",
		BgRed: "BgRed",
		BgGreen: "BgGreen",
		BgYellow: "BgYellow",
		BgBlue: "BgBlue",
		BgMagenta: "BgMagenta",
		BgCyan: "BgCyan",
		BgWhite: "BgWhite",
		BgHiBlack: "BgHiBlack",
		BgHiRed: "BgHiRed",
		BgHiGreen: "BgHiGreen",
		BgHiYellow: "BgHiYellow",
		BgHiBlue: "BgHiBlue",
		BgHiMagenta: "BgHiMagenta",
		BgHiCyan: "BgHiCyan",
		BgHiWhite: "BgHiWhite",
	}
	if s, ok := m[a]; ok {
		return s
	}
	return fmt.Sprintf("unknown color %q", a)
}

// Base attributes
const (
	Reset Attribute = iota
	Bold
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

// Foreground text colors
const (
	FgBlack Attribute = iota + 30
	FgRed
	FgGreen
	FgYellow
	FgBlue
	FgMagenta
	FgCyan
	FgWhite
)

// Foreground Hi-Intensity text colors
const (
	FgHiBlack Attribute = iota + 90
	FgHiRed
	FgHiGreen
	FgHiYellow
	FgHiBlue
	FgHiMagenta
	FgHiCyan
	FgHiWhite
)

// Background text colors
const (
	BgBlack Attribute = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMagenta
	BgCyan
	BgWhite
)

// Background Hi-Intensity text colors
const (
	BgHiBlack Attribute = iota + 100
	BgHiRed
	BgHiGreen
	BgHiYellow
	BgHiBlue
	BgHiMagenta
	BgHiCyan
	BgHiWhite
)
