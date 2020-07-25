package dynamic_gui_config

import (
	"log"
	"os"

	"github.com/andlabs/ui"
)

var (
	// window configuration, set during Start
	windowName   string
	windowWidth  int
	windowHeight int

	window *ui.Window = nil
	tab    *ui.Tab    = nil
)

// perform a check on global variables
func check() {
	if window == nil || tab == nil {
		panic("window pointer is nil, did you call Start()?")
	}
}

func setup() {
	window = ui.NewWindow(windowName, windowWidth, windowHeight, true)
	tab = ui.NewTab()
	window.SetMargined(true)
	window.SetChild(tab)

	// prevent destroy of window
	window.OnClosing(func(w *ui.Window) bool {
		return true
	})
}

func uithread(setupdone chan<- bool) {
	if err := ui.Main(func() {
		setup()
		// signal setup done
		setupdone <- true
	}); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}

// StartDefaults is an alias for Start("config", 640, 400)
func StartDefaults() {
	Start("config", 640, 400)
}

// Start starts the configuration thread and sets up the window handles
// it returns when setup is done and the handles are initialized
func Start(windowname string, width, height int) {
	windowName = windowname
	windowWidth = width
	windowHeight = height

	done := make(chan bool)

	// start the goroutine with ui.Main
	go uithread(done)

	// wait for setup task completion
	<-done
}

// runs in ui thread
func updateUi(fields []structGuiField) (ui.Control, error) {
	container := ui.NewVerticalBox()
	container.SetPadded(true)

	for _, field := range fields {
		hbox := ui.NewHorizontalBox()
		hbox.SetPadded(true)
		hbox.Append(ui.NewLabel(field.Properties.Name), false)
		hbox.Append(field.Factory.Create(), true)

		container.Append(hbox, false)
	}

	return container, nil
}

// Register adds the given struct pointer to the graphical interface as a new tab
func Register(name string, config interface{}) error {
	fields, err := structBreakdown(config)
	if err != nil {
		return err
	}

	// add a new tab to the window and add controls based on struct field
	NewTab(name, func() ui.Control {
		if ctrl, err := updateUi(fields); err != nil {
			log.Println(err)
			return nil
		} else {
			return ctrl
		}
	})

	return nil
}

// NewTab creates a new tab in the configuration window
func NewTab(name string, handle func() ui.Control) {
	ui.QueueMain(func() {
		ctrl := handle()
		if ctrl != nil {
			tab.Append(name, ctrl)
		}
	})
}

// Show displays the configuration window
func Show() {
	ui.QueueMain(window.Show)
}

// Hide hides the configuration window
func Hide() {
	ui.QueueMain(window.Hide)
}
