package color

import (
	"io"
	"os"
	"sync"

	"github.com/mattn/go-colorable"
)

type Console struct {
	sync.Mutex
	out *os.File
	colorable io.Writer
}

var stdout *Console
var stdoutOnce sync.Once

func Stdout() *Console {
	stdoutOnce.Do(func() {
		stdout = &Console{
			colorable: colorable.NewColorable(os.Stdout),
			out: os.Stdout,
		}
	})
	return stdout
}

var stderr *Console
var stderrOnce sync.Once

func Stderr() *Console {
	stderrOnce.Do(func() {
		stderr = &Console{
			colorable: colorable.NewColorable(os.Stderr),
			out: os.Stderr,
		}
	})
	return stderr
}

func(c *Console) StripColors(strip bool) {
	c.Lock()
	defer c.Unlock()
	if strip {
		c.colorable = colorable.NewNonColorable(c.out)
		return
	}
	c.colorable = colorable.NewColorable(c.out)
}

// Writer so we can treat a console as a Writer
func(c *Console) Write(b []byte)(int, error) {
	return c.colorable.Write(b)
}


