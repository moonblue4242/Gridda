//
//   Copyright (C) 2019 moonblue4242@gmail.com
//

package actions

import (
	"fmt"
	"strings"

	"github.com/moonblue4242/Gridda/winapi"
)

// ExpandActions incorporates all the action for expanding a window
type ExpandActions struct {
	expandedWindow TargetWindow
}

// NewExpandActions creates a new initialized expands actions structure which can be called to execute the action
func NewExpandActions() *ExpandActions {
	expandActions := new(ExpandActions)
	return expandActions
}

// ExpandHandler returns the action handler for expanding the currently focused window
func (expandActions *ExpandActions) ExpandHandler() ActionHandler {
	return func(activeConfig ActiveConfig) {
		expandActions.expand(GetTarget(), activeConfig)
	}
}

// Expand will expand the window given via handle if a preset exists
func (expandActions *ExpandActions) Expand(hwnd winapi.Hwnd, activeConfig ActiveConfig) {
	expandActions.expand(GetTargetFromHandle(hwnd), activeConfig)
}

func (expandActions *ExpandActions) expand(target TargetWindow, activeConfig ActiveConfig) {
	name, err := target.ModuleName()
	fmt.Printf("MODULE:: %s\n", name)
	if expandActions.expandedWindow != nil {
		rect := expandActions.expandedWindow.Size()
		// no need to calculate something as target window has the latest data stored
		expandActions.expandedWindow.Move(int(rect.Left), int(rect.Top), int(rect.Width()), int(rect.Bottom))
		expandActions.expandedWindow = nil
	}
	if err == nil && strings.ToUpper("firefox.exe") == name {
		left := getWeightedPosition(2-1, activeConfig.Grid().Columns, calcDistancePerWeightPx(target.DesktopSize().Width(), &activeConfig.Grid().Columns))
		fmt.Sprintln("Firefox")
		deltaH, _ := target.Delta()
		if left == int(target.Size().Left)+deltaH {
			rect := target.Size()
			fmt.Printf("EXPANDING:: %s\n", rect)
			target.Move(int(rect.Left-400), int(rect.Top), int(rect.Width()+400), int(rect.Bottom))
			expandActions.expandedWindow = target
		}
	}

}
