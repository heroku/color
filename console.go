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
	*colorCache
	colored    io.Writer
	noncolored io.Writer
	current    io.Writer
}

// NewConsole creates a wrapper around out which will output platform independent colored text.
func NewConsole(out *os.File) *Console {
	c := &Console{
		colored:    colorable.NewColorable(out),
		noncolored: colorable.NewNonColorable(out),
	}
	c.current = c.colored
	c.init()
	return c
}

// DisableColors pass a flag that will remove color information from console output if true,
// otherwise color information is included by default.
func (c *Console) DisableColors(strip bool) {
	c.Lock()
	defer c.Unlock()
	if strip {
		c.current = c.noncolored
		return
	}
	c.current = c.colored
}

// Set will cause the color passed in as an argument to be written until Unset is called.
func (c *Console) Set(color *Color) {
	c.Lock()
	defer c.Unlock()
	_, _ = c.current.Write([]byte(color.colorStart))
}

// Unset will restore console output to default. It will undo colored console output from a call to Set.
func (c *Console) Unset() {
	c.Lock()
	defer c.Unlock()
	_, _ = c.current.Write([]byte(colorReset))
}

// Write so we can treat a console as a Writer
func (c *Console) Write(b []byte) (int, error) {
	c.Lock()
	n, err := c.current.Write(b)
	c.Unlock()
	return n, err
}

func (c *Console) init() {
	c.colorCache = &colorCache{
		cache:  make(colorMap),
		parent: c,
	}
}

type colorMap map[Attribute]*Color

type colorCache struct {
	sync.RWMutex
	cache  colorMap
	parent writerValuer
}

func newValueCache(w writerValuer) *colorCache {
	return &colorCache{
		cache:  make(colorMap),
		parent: w,
	}
}

func (cc *colorCache) value(attrs ...Attribute) *Color {
	key := to_key(attrs)
	if v := cc.getIfExists(key); v != nil {
		return v
	}
	cc.Lock()
	defer cc.Unlock()
	v := &Color{
		colorStart: chainSGRCodes(attrs),
		out:        cc.parent,
	}
	cc.cache[key] = v
	return v
}

func (vc *colorCache) getIfExists(key Attribute) *Color {
	vc.RLock()
	defer vc.RUnlock()

	if v, ok := vc.cache[key]; ok {
		return v
	}
	return nil
}
