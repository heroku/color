// +build !windows 

package color 


type ConsoleScreenBufferInfo struct {}

func(c *Console) set(attr TextAttribute) error {
	return nil 
}

func(c *Console) reset() error {
	return nil 
}