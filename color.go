package color

import (
	"log"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/mattn/go-colorable"
)

type consoleType int 

const (
	typeUnknown consoleType = iota
	typeStdout
	typeStderr
)


type Console struct {
	sync.Mutex
	out io.Writer 
	typ consoleType 
	state ConsoleScreenBufferInfo 
}
	
var stdout *Console
var stdoutOnce sync.Once 

func Stdout() *Console {
	stdoutOnce.Do(func(){
		stdout = &Console{
			out: colorable.NewColorableStdout(),
			typ: typeStdout,
		}
	})
	return stdout 
}

var stderr *Console 
var stderrOnce sync.Once 

func Stderr() *Console {
	stderrOnce.Do(func(){
		stderr = &Console{
			out: colorable.NewColorableStderr(), 
			typ: typeStderr,
		}
	})
}

func(c *Console) Println(attr TextAttribute, a ...interface{})(int, error){

	return fmt.Fprintln(c.out, a...)
}

func reset(c *Console) {
	if err := c.reset(); err != nil {
		log.Fatal(err)
	}
}
