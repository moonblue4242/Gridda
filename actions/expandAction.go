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
	expandedPreset *Preset
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
	// resize to previous size only if not moved
	if expandActions.expandedWindow != nil {
		reaquiredTarget := GetTargetFromHandle(expandActions.expandedWindow.Hwnd())
		if isAt(reaquiredTarget, activeConfig, &expandActions.expandedPreset.Expandable.Where, &expandActions.expandedPreset.Expandable.How) {
			rect := expandActions.expandedWindow.Size()
			// no need to calculate something as target window has the latest data stored
			expandActions.expandedWindow.Move(int(rect.Left), int(rect.Top), int(rect.Width()), int(rect.Height()))
		}
	}
	expandActions.expandedWindow = nil
	if err == nil {
		if preset := activeConfig.Grid().Presets.FindFirst(name); preset != nil && preset.Expandable != nil {
			fmt.Sprintln(preset.Executable)
			if isAt(target, activeConfig, &preset.Expandable.Where, nil) {
				how := preset.Expandable.How
				target.Expand(how.Left, how.Right, how.Top, how.Bottom)
				expandActions.expandedWindow = target
				expandActions.expandedPreset = preset
			}
		}
	}

}

func isAt(target TargetWindow, activeConfig ActiveConfig, where *Where, how *How) bool {
	distancePerWeightPx := activeConfig.Grid().Columns.fractionPerWeight(int(target.DesktopSize().Width()))
	distancePerHeightPx := activeConfig.Grid().Rows.fractionPerWeight(int(target.DesktopSize().Height()))

	span := &Span{0, 0}
	if where.Span != nil {
		span = where.Span
	}

	// initialize to zero if missing
	if how == nil {
		how = &How{0, 0, 0, 0}
	}

	left, top, width, height := calcCorrectedPosition(target, activeConfig.Grid(), where.Column, where.Row, span, distancePerWeightPx, distancePerHeightPx)
	// add expand to allow checking against expanded window
	left = left - how.Left
	top = top - how.Top
	width = width + how.Left + how.Right
	height = height + how.Top + how.Bottom
	fmt.Printf("@: %d,%d,%d,%d\n", left, top, width, height)
	// check against target, use span for width and height if available
	result := left == int(target.Size().Left)
	result = result && top == int(target.Size().Top)
	if where.Span != nil {
		result = result && width == int(target.Size().Width())
		result = result && height == int(target.Size().Height())
	}

	return result

}
