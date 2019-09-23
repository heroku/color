package main

import (
	"fmt"

	"github.com/murphybytes/color"
)

func main() {
	v := color.New(color.BgBlack, color.FgHiMagenta)
	_, _ = color.Stdout().Print(v, "here is some text\n")
	color.HiMagenta("hello megenta %s", "foo")
	color.HiCyan("hello cyan %d", 10)

	emphasized := color.New(color.FgBlue, color.FgRed, color.Bold).SprintFunc()
	fmt.Fprintln(color.Stdout(), "Wow! This is", emphasized("exciting!"))

	color.Stdout().Println(color.New(), "no color at all")
	c := color.New(color.FgHiMagenta)
	color.Stdout().Set(c)
	fmt.Println("magenta?")
	fmt.Fprintln(color.Stdout(), "defin magenta")
	color.Stdout().Unset()
	fmt.Fprintln(color.Stdout(), "should be normal")

}
