package color

import (
	"log"
	"fmt"
	"io"
	"os"
	"sync"
)

// TextAttribute describes foreground, and background colors and optionally intensity
type TextAttribute uint16

const (
	ForegroundBlue TextAttribute = 1  
	ForegroundGreen = ForegroundBlue << 1
	ForegroundRed = ForegroundGreen << 1 
	ForegroundIntensity = ForegroundRed << 1 
	BackgroundBlue = ForegroundIntensity << 1 
	BackgroundGreen = BackgroundBlue << 1
	BackgroundRed = BackgroundGreen << 1 
	BackgroundIntensity = BackgroundRed << 1 
	BlackForegroundWhiteBackground = BackgroundBlue | BackgroundGreen | BackgroundRed
	WhiteForegroundBlackBackground = ForegroundBlue | ForegroundGreen | ForegroundRed
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
			out: os.Stdout,
			typ: typeStdout,
		}
	})
	return stdout 
}

func(c *Console) Println(attr TextAttribute, a ...interface{})(int, error){
	if err := c.set(attr); err != nil {
		return 0, err 
	}
	defer reset(c)
	return fmt.Fprintln(c.out, a...)
}

func reset(c *Console) {
	if err := c.reset(); err != nil {
		log.Fatal(err)
	}
}
