package main 

import (
	"fmt"
	"github.com/murphybytes/color"
)

func main(){
	so := color.Stdout()
	_, err := so.Println(color.ForegroundRed, "foo hoo hoo")
	if err != nil {
		fmt.Println("error", err)
	}
}