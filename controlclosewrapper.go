package dynamic_gui_config

import (
	"io"
	"log"

	"github.com/andlabs/ui"
)

type CloseFunc func() error

func (c CloseFunc) Close() error {
	return c()
}

type ControlCloseWrapper struct {
	Child        ui.Control
	CloseHandles []io.Closer
}

func NewControlCloseWrapper(child ui.Control) *ControlCloseWrapper {
	return &ControlCloseWrapper{
		Child:        child,
		CloseHandles: make([]io.Closer, 0),
	}
}

func (c *ControlCloseWrapper) AddClosers(closer ...io.Closer) {
	c.CloseHandles = append(c.CloseHandles, closer...)
}

func (c ControlCloseWrapper) LibuiControl() uintptr {
	return c.Child.LibuiControl()
}

func (c ControlCloseWrapper) Destroy() {
	for _, closer := range c.CloseHandles {
		if err := closer.Close(); err != nil {
			log.Println("err closing: ", err)
		}
	}
	c.Child.Destroy()
}

func (c ControlCloseWrapper) Handle() uintptr {
	return c.Child.Handle()
}

func (c ControlCloseWrapper) Visible() bool {
	return c.Child.Visible()
}

func (c ControlCloseWrapper) Show() {
	c.Child.Show()
}

func (c ControlCloseWrapper) Hide() {
	c.Child.Hide()
}

func (c ControlCloseWrapper) Enabled() bool {
	return c.Child.Enabled()
}

func (c ControlCloseWrapper) Enable() {
	c.Child.Enable()
}

func (c ControlCloseWrapper) Disable() {
	c.Child.Disable()
}
