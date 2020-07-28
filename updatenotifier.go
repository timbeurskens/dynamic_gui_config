package dynamic_gui_config

// UpdateNotifier is an interface type a user can implement to get notified when a value changes
type UpdateNotifier interface {
	// OnValueChanged will be called when the user changes the value using the graphical interface
	OnValueChanged()
}
