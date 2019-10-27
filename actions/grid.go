//
//   Copyright (C) 2019 moonblue4242@gmail.com
//

package actions

import "errors"

// Weights define numbers used for weighting columns or rows
type Weights []int

// Grid defines a basic grid on the screen at which windows can be snapped to
type Grid struct {
	Name    string
	Columns Weights
	Rows    Weights
	Presets Presets
}

func (weights *Weights) sum() (weight int) {
	weight = 0
	for _, column := range []int(*weights) {
		weight += column
	}
	return
}

func (weights *Weights) fractionPerWeight(total int) int {
	return total / weights.sum()
}

func (weights *Weights) weightedSumFrom(startIndex int, amountPerWeight int) (pos int) {
	pos = 0
	if startIndex >= 0 && startIndex < len(*weights) {
		for i := 0; i <= startIndex; i++ {
			pos += (*weights)[i] * amountPerWeight
		}
	}
	return
}

// Validate determines if the grid object is setup correctly
func (grid *Grid) Validate() (err error) {
	err = nil
	if len(grid.Name) < 1 {
		err = errors.New("Grid must have a name with at least one character")
	}
	return
}
