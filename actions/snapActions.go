//
//   Copyright (C) 2019 moonblue4242@gmail.com
//

package actions

import (
	"strings"

	"github.com/moonblue4242/Gridda/winapi"
)

// SnapActions defines a set of actions for snaping windows to a grid, this structure
// is NOT usable without initialization
type SnapActions struct {
	spans map[int]map[winapi.Hwnd]*Span
}

// NewSnapActions creates a new initialized snap actions structure
func NewSnapActions() *SnapActions {
	snapActions := new(SnapActions)
	snapActions.spans = make(map[int]map[winapi.Hwnd]*Span)
	return snapActions
}

// Span defines the amount of grid tiles the window should span over
type Span struct {
	Columns int
	Rows    int
}

// ToLeft returns the action handler for snapping to the next grid line on the left
func (snapActions *SnapActions) ToLeft() ActionHandler {
	return func(activeConfig ActiveConfig) {
		snapActions.snapLeft(GetTarget(), activeConfig)
	}
}

// ToRight returns the action handler for snapping to the next grid line on the right
func (snapActions *SnapActions) ToRight() ActionHandler {
	return func(activeConfig ActiveConfig) {
		snapActions.snapRight(GetTarget(), activeConfig)
	}
}

// ToTop returns the action handler for snapping to the next grid line on the top
func (snapActions *SnapActions) ToTop() ActionHandler {
	return func(activeConfig ActiveConfig) {
		snapActions.snapTop(GetTarget(), activeConfig)
	}
}

// ToBottom returns the action handler for snapping to the next grid line on the top
func (snapActions *SnapActions) ToBottom() ActionHandler {
	return func(activeConfig ActiveConfig) {
		snapActions.snapBottom(GetTarget(), activeConfig)
	}
}

// SpanHorizontal increases/decrease the horizontal span of grid tiles of the target window
func (snapActions *SnapActions) SpanHorizontal(increase bool) ActionHandler {
	return func(activeConfig ActiveConfig) {
		snapActions.spanHorizontal(GetTarget(), activeConfig, increase)
	}
}

// SpanVertical increases/decrease the vertical span of grid tiles of the target window
func (snapActions *SnapActions) SpanVertical(increase bool) ActionHandler {
	return func(activeConfig ActiveConfig) {
		snapActions.spanVertical(GetTarget(), activeConfig, increase)
	}
}

// snapMovement returns the index of the grid tile the window should be moved
type snapMovement func(gridLeftPos int32, gridLeftIndex int, gridTopPos int32, gridTopIndex int, correctedLeft int32) (newGridLeftIndex int, newGridTopIndex int)

// SnapLeft will attach the window to the next grid line on the left
func (snapActions *SnapActions) snapLeft(target TargetWindow, activeConfig ActiveConfig) {
	snap(target, activeConfig, snapActions,
		func(gridLeftPos int32, gridLeftIndex int, _ int32, gridTopIndex int, correctedLeft int32) (_ int, _ int) {
			if gridLeftPos == correctedLeft && gridLeftIndex > 0 {
				gridLeftIndex--
			}
			return gridLeftIndex, gridTopIndex
		})
}

// SnapLeft will attach the window to the next grid line on the left
func (snapActions *SnapActions) snapTop(target TargetWindow, activeConfig ActiveConfig) {
	snap(target, activeConfig, snapActions,
		func(_ int32, gridLeftIndex int, gridTopPos int32, gridTopIndex int, _ int32) (_ int, _ int) {
			if gridTopPos == target.Size().Top && gridTopIndex > 0 {
				gridTopIndex--
			}
			return gridLeftIndex, gridTopIndex
		})
}

// SnapRight will attach the window to the next grid line on the right
func (snapActions *SnapActions) snapRight(target TargetWindow, activeConfig ActiveConfig) {
	snap(target, activeConfig, snapActions,
		func(_ int32, gridLeftIndex int, _ int32, gridTopIndex int, _ int32) (_ int, _ int) {
			if gridLeftIndex+1 < len(activeConfig.Grid().Columns) {
				gridLeftIndex++
			}
			return gridLeftIndex, gridTopIndex
		})
}

// SnapRight will attach the window to the next grid line on the right
func (snapActions *SnapActions) snapBottom(target TargetWindow, activeConfig ActiveConfig) {
	snap(target, activeConfig, snapActions,
		func(_ int32, gridLeftIndex int, _ int32, gridTopIndex int, _ int32) (_ int, _ int) {
			if gridTopIndex+1 < len(activeConfig.Grid().Rows) {
				gridTopIndex++
			}
			return gridLeftIndex, gridTopIndex
		})
}

// spanHorizontal will increase/decrease the span bound by the amount of columns
// additionally a snap is performed to refresh
func (snapActions *SnapActions) spanHorizontal(target TargetWindow, activeConfig ActiveConfig, increase bool) {
	item := snapActions.getOrSetSpan(activeConfig, target)
	if increase && item.Columns < len(activeConfig.Grid().Columns) {
		item.Columns++
	} else if !increase && item.Columns > 1 {
		item.Columns--
	}
	// refresh by snaping in place
	snap(target, activeConfig, snapActions, nil)
}

// spanHorizontal will increase/decrease the span bound by the amount of columns
// additionally a snap is performed to refresh
func (snapActions *SnapActions) spanVertical(target TargetWindow, activeConfig ActiveConfig, increase bool) {
	item := snapActions.getOrSetSpan(activeConfig, target)
	if increase && item.Rows < len(activeConfig.Grid().Rows) {
		item.Rows++
	} else if !increase && item.Rows > 1 {
		item.Rows--
	}
	// refresh by snaping in place
	snap(target, activeConfig, snapActions, nil)
}

func (snapActions *SnapActions) getOrSetSpan(activeConfig ActiveConfig, target TargetWindow) *Span {
	spanMap, ok := snapActions.spans[activeConfig.GridIndex()]
	if !ok {
		spanMap = make(map[winapi.Hwnd]*Span)
		snapActions.spans[activeConfig.GridIndex()] = spanMap
	}
	item, ok := spanMap[target.Hwnd()]
	if !ok {
		item = new(Span)
		item.Columns, item.Rows = snapActions.getPresetSpanOrDefault(target, activeConfig)
		spanMap[target.Hwnd()] = item
	}
	return item
}

// getPresetSpanOrDefault will check if for the given base module the window is controlled by a preset exists
func (snapActions *SnapActions) getPresetSpanOrDefault(targetWindow TargetWindow, activeConfig ActiveConfig) (columns int, rows int) {
	columns = 1
	rows = 1
	if moduleName, err := targetWindow.ModuleName(); err == nil {
		for _, preset := range activeConfig.Grid().Presets {
			if strings.ToUpper(preset.Executable) == moduleName {
				columns = preset.Span.Columns
				rows = preset.Span.Rows
				break
			}
		}
	}
	return
}

func snap(target TargetWindow, activeConfig ActiveConfig, snapActions *SnapActions, snapMovement snapMovement) {
	grid := activeConfig.Grid()
	span := snapActions.getOrSetSpan(activeConfig, target)
	// calculate corrections for special borders (e.g. drop shadows)
	deltaH, _ := target.Delta()
	correctedLeft := target.Size().Left + int32(deltaH)
	// calculate the tile index of the grid for x and y axis along with top, left edge of the tile
	gridLeftPos, gridLeftIndex, widthPerWeightPx := getGridTile(correctedLeft, target.DesktopSize().Width(), grid.Columns)
	gridTopPos, gridTopIndex, heightPerWeightPx := getGridTile(target.Size().Top, target.DesktopSize().Height(), grid.Rows)
	// perform movement
	if snapMovement != nil {
		gridLeftIndex, gridTopIndex = snapMovement(gridLeftPos, gridLeftIndex, gridTopPos, gridTopIndex, correctedLeft)
	}
	// position based on the grid tile index
	move(target, span, grid, gridLeftIndex, gridTopIndex, widthPerWeightPx, heightPerWeightPx)
}

// move the target window to the grid tile specified taking into account any necessary border corrections
func move(target TargetWindow, span *Span, grid *Grid, gridLeftIndex int, gridTopIndex int, widthPerWeightPx int, heightPerWeightPx int) {
	left, top, width, height := calcCorrectedPosition(target, grid, gridLeftIndex, gridTopIndex, span, widthPerWeightPx, heightPerWeightPx)
	target.Move(left, top, width, height)
}

func calcCorrectedPosition(target TargetWindow, grid *Grid, column int, row int, span *Span, widthPerWeightPx int, heightPerWeightPx int) (left int, top int, width int, height int) {
	deltaH, deltaV := target.Delta()
	left = grid.Columns.weightedSumFrom(column-1, int(widthPerWeightPx))
	left = left - deltaH
	top = grid.Rows.weightedSumFrom(row-1, int(heightPerWeightPx))
	width = getSpannedDistance(column, grid.Columns, widthPerWeightPx, span.Columns)
	width = width + deltaH*2
	height = getSpannedDistance(row, grid.Rows, heightPerWeightPx, span.Rows)
	height = height + deltaV
	return
}

func getSpannedDistance(weightIndex int, weights []int, distancePerWeightPx int, span int) (distance int) {
	distance = 0
	for i := 0; i < span && weightIndex+i < len(weights); i++ {
		distance += weights[weightIndex+i] * distancePerWeightPx
	}
	return
}

func getGridTile(pos int32, maxDistance int32, weights Weights) (gridStart int32, gridIndex int, distancePerWeightPx int) {
	gridStart = 0
	gridIndex = 0
	distancePerWeightPx = weights.fractionPerWeight(int(maxDistance))
	for idx, weight := range weights {
		additionalDistance := int32(weight) * int32(distancePerWeightPx)
		gridIndex = idx
		if gridStart+additionalDistance > pos {
			break
		}
		gridStart += additionalDistance
	}
	return
}
