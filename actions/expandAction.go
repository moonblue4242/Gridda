//
//   Copyright (C) 2019 moonblue4242@gmail.com
//

package actions

import (
	"fmt"

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
	if err == nil {
		if preset := activeConfig.Grid().Presets.FindFirst(name); preset != nil && preset.Expandable {

			distancePerWeightPx := activeConfig.Grid().Columns.fractionPerWeight(int(target.DesktopSize().Width()))
			distancePerHeightPx := activeConfig.Grid().Rows.fractionPerWeight(int(target.DesktopSize().Height()))

			left, _, _, _ := calcCorrectedPosition(target, activeConfig.Grid(), 2, 0, preset.Span, distancePerWeightPx, distancePerHeightPx)
			fmt.Sprintln("Firefox")
			if left == int(target.Size().Left) {
				rect := target.Size()
				fmt.Printf("EXPANDING:: %s\n", rect)
				target.Move(int(rect.Left-400), int(rect.Top), int(rect.Width()+400), int(rect.Bottom))
				expandActions.expandedWindow = target
			}
		}
	}

}
