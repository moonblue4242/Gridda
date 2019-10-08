package main

import (
	"log"

	"sonnenfroh.de/test/actions"
	"sonnenfroh.de/test/cmds"
	"sonnenfroh.de/test/ui"
	"sonnenfroh.de/test/winapi"
)

// Logic defines the externally callable methods
type Logic interface {
	ui.Actions
	Loop()
}

type logic struct {
	ui               *ui.UI
	config           *cmds.Config
	commander        *cmds.Commander
	currentGridIndex int
}

// NewLogic creates a new business logic controller mapping UI to business functions and vica versa
func NewLogic() Logic {
	logic := new(logic)
	var err error
	// load config and apply
	logic.config, err = cmds.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Config loading failed: %s\n", err)
	}
	// setup key binding service
	logic.commander = new(cmds.Commander)
	logic.commander.Apply(logic.config, logic)
	// setup ui
	logic.ui, err = ui.New(logic, logic.config, func(msg *winapi.Message) {
		logic.commander.HandleHotkey(msg)
	})
	if err != nil {
		log.Fatalf("UI Setup failed: %s\n", err)
	}
	return logic
}

// Run the event loop
func (logic *logic) Loop() {
	logic.ui.Loop()
}

// AttachUI attaches the ui for retrieving events from
func (logic *logic) AttachUI(ui *ui.UI) {
	logic.ui = ui
}

func (logic *logic) OnExitApplication() {
	if logic != nil && logic.ui != nil {
		logic.ui.Quit()
	}
}

func (logic *logic) OnShowApplication() {
	if logic != nil && logic.ui != nil {
		logic.ui.ShowMain()
	}
}

func (logic *logic) Grid() *actions.Grid {
	return &logic.config.Grids[logic.currentGridIndex]
}

func (logic *logic) GridIndex() int {
	return logic.currentGridIndex
}

func (logic *logic) OnPreviousGrid() {
	idx := logic.currentGridIndex - 1
	if idx < 0 {
		idx = len(logic.config.Grids) - 1
	}
	logic.OnGridSelect(idx)
}

func (logic *logic) OnNextGrid() {
	logic.OnGridSelect((logic.currentGridIndex + 1) % len(logic.config.Grids))
}

func (logic *logic) OnGridSelect(index int) {
	logic.ui.SetGrid(logic.config, index)
	logic.currentGridIndex = index
}
