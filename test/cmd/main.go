package main

import (
	"github.com/murphybytes/color"
)

func main() {
	so := color.Stdout()
	v, _ := color.New(so, color.BgBlack, color.FgGreen)
	_, _ = v.Print("here is some txrt\n")
	//so.StripColors(true)
	//vv := v.Add(color.Underline)
	//vv.Print("something new\n")
	v.Print("something old\n")

	color.HiMagenta("hello megenta %s", "foo")
	color.HiCyan("hello cyan %d", 10)
}
