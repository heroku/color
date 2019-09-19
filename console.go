// Package color produces colored output in terms of ANSI Escape Codes. Posix and Windows platforms are supported.
package color

import (
	"io"
	"os"
	"sync"

	"github.com/mattn/go-colorable"
)

var stdout *Console // Don't use directly use Stdout() instead.
var stdoutOnce sync.Once

// Stdout returns an io.Writer that writes colored text to standard out.
func Stdout() *Console {
	stdoutOnce.Do(func() {
		stdout = NewConsole(os.Stdout)
	})
	return stdout
}

var stderr *Console // Don't use directly use Stderr() instead.
var stderrOnce sync.Once

// Stderr returns an io.Writer that writes colored text to standard error.
func Stderr() *Console {
	stderrOnce.Do(func() {
		stderr = NewConsole(os.Stderr)
	})
	return stderr
}

// Console manages state for output, typically stdout or stderr.
type Console struct {
	sync.Mutex
	*valueCache
	out       *os.File
	colorable io.Writer
}

// NewConsole creates a wrapper around out which will output platform independent colored text.
func NewConsole(out *os.File) *Console {
	c := &Console{
		colorable: colorable.NewColorable(out),
		out:       out,
	}
	c.init()
	return c
}

// DisableColors pass a flag that will remove color information from console output if true,
// otherwise color information is included by default.
func (c *Console) DisableColors(strip bool) {
	c.Lock()
	defer c.Unlock()
	if strip {
		c.colorable = colorable.NewNonColorable(c.out)
		return
	}
	c.colorable = colorable.NewColorable(c.out)
}

// Write so we can treat a console as a Writer
func (c *Console) Write(b []byte) (int, error) {
	c.Lock()
	defer c.Unlock()
	return c.colorable.Write(b)
}

func (c *Console) init() {
	c.valueCache = &valueCache{
		cache:  make(valueMap),
		parent: c,
	}
}

type valueMap map[Attribute]*Color

type valueCache struct {
	sync.RWMutex
	cache  valueMap
	parent writerValuer
}

func newValueCache(w writerValuer) *valueCache {
	return &valueCache{
		cache:  make(valueMap),
		parent: w,
	}
}

func (vc *valueCache) value(attr Attribute) *Color {
	if v := vc.getIfExists(attr); v != nil {
		return v
	}
	vc.Lock()
	defer vc.Unlock()
	v := NewWithWriter(vc.parent, attr)
	vc.cache[attr] = v
	return v
}

func (vc *valueCache) getIfExists(attr Attribute) *Color {
	vc.RLock()
	defer vc.RUnlock()
	if v, ok := vc.cache[attr]; ok {
		return v
	}
	return nil
}
