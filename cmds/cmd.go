package cmds

import (
	"sonnenfroh.de/test/actions"
	"sonnenfroh.de/test/winapi"
)

const (
	commandOffset = 4714
)

// Command is the basic interface for all actions executable by a hotkey
type Command func()

// Commander creates bindings for hotkeys to functions
type Commander struct {
	bindings []Command
}

// NewCommander creates a new initialized element
func (commander *Commander) NewCommander() {
	commander.bindings = make([]Command, 0)
}

// Bind will attach a command to a hotkey
func (commander *Commander) Bind(vkey int32, alt bool, ctrl bool, shift bool, command Command) bool {
	commander.bindings = append(commander.bindings, command)
	return winapi.RegisterHotVkey(len(commander.bindings)+commandOffset-1, vkey, alt, ctrl, shift)
}

// Apply will apply the current configuration and its hotkey bindings deleting the previous ones
func (commander *Commander) Apply(config *Config, activeConfig actions.ActiveConfig) {
	// deregister any previous hotkeys
	for idx := range commander.bindings {
		winapi.UnregisterHotVkey(commandOffset + idx)
	}
	// create new bindings
	for _, binding := range config.Bindings {
		action := availableActions[binding.Action] // create new variable inside closure context
		commander.Bind(keys[binding.Key], binding.Alt, binding.Ctrl, binding.Shift, func() { action(activeConfig) })
	}
}

// HandleHotkey must be called during handling of the WM_HOTKEY event
func (commander *Commander) HandleHotkey(msg *winapi.Message) {
	if msg.Message == winapi.WM_HOTKEY && msg.WParam >= commandOffset {
		commander.bindings[msg.WParam-commandOffset]()
	}
}
