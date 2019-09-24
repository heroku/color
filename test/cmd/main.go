package main

import (
	"fmt"

	"github.com/heroku/color"
)

func main() {
	v := color.New(color.BgBlack, color.FgHiMagenta)
	_, _ = color.Stdout().Print(v, "here is some text\n")
	color.HiMagenta("hello megenta %s", "foo")
	color.HiCyan("hello cyan %d", 10)

	emphasized := color.New(color.FgBlue, color.FgRed, color.Bold).SprintFunc()
	_, _ = fmt.Fprintln(color.Stdout(), "Wow! This is", emphasized("exciting!"))

	_, _ = color.Stdout().Println(color.New(), "no color at all")
	c := color.New(color.FgHiMagenta)
	color.Stdout().Set(c)
	fmt.Println("magenta?")
	_, _ = fmt.Fprintln(color.Stdout(), "defin magenta")
	color.Stdout().Unset()
	_, _ = fmt.Fprintln(color.Stdout(), "should be normal")

}
