package main

import (
	"fmt"

	"github.com/murphybytes/color"
)

func main() {
	v := color.New(color.BgBlack, color.FgGreen)
	_, _ = v.Print("here is some txrt\n")
	v.Print("something old\n")

	color.HiMagenta("hello megenta %s", "foo")
	color.HiCyan("hello cyan %d", 10)

	emphasized := color.New(color.FgBlue, color.FgRed, color.Bold).SprintFunc()
	fmt.Fprintln(color.Stdout(), "Wow! This is", emphasized("exciting!"))

	color.New().Println("no color at all")

}
