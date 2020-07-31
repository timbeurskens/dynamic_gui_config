package dynamic_gui_config

import (
	"io"
	"log"

	"github.com/andlabs/ui"
)

// CloseFunc is a function type implementing the io.Closer interface
type CloseFunc func() error

// Close runs given the close function
func (c CloseFunc) Close() error {
	return c()
}

// ControlCloseWrapper acts as a ui.Control and supports close functions.
// This can be used to wrap a control and run a close function when the control gets destroyed.
type ControlCloseWrapper struct {
	Child        ui.Control
	CloseHandles []io.Closer
}

// NewControlCloseWrapper creates a ControlCloseWrapper with the given control type as child element
func NewControlCloseWrapper(child ui.Control) *ControlCloseWrapper {
	return &ControlCloseWrapper{
		Child:        child,
		CloseHandles: make([]io.Closer, 0),
	}
}

// AddClosers adds the given io.Closer objects to the ControlCloseWrapper container
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
