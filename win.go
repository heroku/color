// +build windows 

package color 

import (
	"syscall"
	"unsafe"
)

var (
	kernel = syscall.NewLazyDLL("kernel32.dll")
	procGetStdHandle = kernel.NewProc("GetStdHandle")
	procSetConsoleTextAttribute = kernel.NewProc("SetConsoleTextAttribute")
	procGetConsoleScreenBufferInfo = kernel.NewProc("GetConsoleScreenBufferInfo")

)

type Coord struct {
	X int16
	Y int16 
}

type SmallRect struct {
	Left int16 
	Top int16 
	Right int16 
	Bottom int16 
}

type ConsoleScreenBufferInfo struct {
	Size Coord 
	CursorPosition Coord 
	Attributes uint32 
	Window SmallRect
	MaximumWindowSize Coord 
}

func(c *Console) set(attr TextAttribute) error {
	var h syscall.Handle 
	switch c.typ {
	case typeStdout :
		h = syscall.Stdout
	case typeStderr:
		h = syscall.Stderr
	}

	_, _, err := procGetConsoleScreenBufferInfo.Call(
		uintptr(h),
		uintptr(unsafe.Pointer(&c.state)),
	) 
	if err != syscall.Errno(0) {
		return err 
	}

	_, _, err = procSetConsoleTextAttribute.Call(
		uintptr(h),
		uintptr(attr),
	)
	if err != syscall.Errno(0) {
		return err 
	}
	return nil 
}

func(c *Console) reset() error {
	var h syscall.Handle 
	switch c.typ {
	case typeStdout :
		h = syscall.Stdout
	case typeStderr:
		h = syscall.Stderr
	}	
	_, _, err := procSetConsoleTextAttribute.Call(
		uintptr(h),
		uintptr(c.state.Attributes),
	)
	return err 
}