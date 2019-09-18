package color

import (
	"io"
	"os"
	"sync"

	"github.com/mattn/go-colorable"
)

type Console struct {
	sync.Mutex
	*valueCache
	out       *os.File
	colorable io.Writer
}

var stdout *Console
var stdoutOnce sync.Once

func Stdout() *Console {
	stdoutOnce.Do(func() {
		stdout = &Console{
			colorable: colorable.NewColorable(os.Stdout),
			out:       os.Stdout,
		}
		stdout.init()
	})
	return stdout
}

var stderr *Console
var stderrOnce sync.Once

func Stderr() *Console {
	stderrOnce.Do(func() {
		stderr = &Console{
			colorable: colorable.NewColorable(os.Stderr),
			out:       os.Stderr,
		}
		stderr.init()
	})
	return stderr
}

func (c *Console) StripColors(strip bool) {
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

type valueMap map[Attribute]*Value

type valueCache struct {
	sync.RWMutex
	cache  valueMap
	parent io.Writer
}

func newValueCache(w io.Writer) *valueCache {
	return &valueCache{
		cache:   make(valueMap),
		parent:  w,
	}
}

func (vc *valueCache) value(attr Attribute) *Value {
	if v := vc.getIfExists(attr); v != nil {
		return v
	}
	vc.Lock()
	defer vc.Unlock()
	v, _ := New(vc.parent, attr)
	vc.cache[attr] = v
	return v
}

func (vc *valueCache) getIfExists(attr Attribute) *Value {
	vc.RLock()
	defer vc.RUnlock()
	if v, ok := vc.cache[attr]; ok {
		return v
	}
	return nil
}
